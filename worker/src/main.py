"""Worker"""
import asyncio
import logging
import os
from typing import AsyncIterable

import grpc

from src.store import init_storage
from src.components import process_metadata
from src.protos.api_pb2_grpc import APIStub
from src.protos.worker_pb2 import MediaItemProcessRequest, MediaItemProcessResponse  # pylint: disable=no-name-in-module
from src.protos.worker_pb2_grpc import WorkerServicer, add_WorkerServicer_to_server


class WorkerService(WorkerServicer):
    """Worker gRPC Service"""

    def __init__(self, file_storage, api_stub: APIStub) -> None:
        self.file_storage = file_storage
        self.api_stub = api_stub

    # pylint: disable=invalid-overridden-method
    async def MediaItemProcess(self, request_iterator: AsyncIterable[
            MediaItemProcessRequest], unused_context) -> MediaItemProcessResponse:
        """MediaItem Process"""
        mediaitem_user_id = None
        mediaitem_id = None
        mediaitem_command = None
        async for mediaitem in request_iterator:
            try:
                self.file_storage.upload(mediaitem_id=mediaitem.id,
                                         content=mediaitem.content)
                mediaitem_user_id = mediaitem.userId
                mediaitem_id = mediaitem.id
                mediaitem_command = mediaitem.command
            except Exception as e:
                logging.error(f'error processing mediaitem for storage: {str(e)}', {
                              'id': mediaitem.id, 'offset': mediaitem.offset})
                mediaitem_id = None
                mediaitem_user_id = None
                return MediaItemProcessResponse(ok=False)
        if mediaitem_id is not None and mediaitem_user_id is not None and 'finish' in mediaitem_command:
            loop = asyncio.get_event_loop()
            loop.create_task(process_metadata(
                self.file_storage, self.api_stub, mediaitem_user_id, mediaitem_id))
        return MediaItemProcessResponse(ok=True)


async def serve() -> None:
    """Main serve function"""
    # initialize storage
    file_storage = init_storage(os.getenv('CAROUSEL_STORAGE', 'disk'))

    # initialize grpc client
    api_host = os.getenv('CAROUSEL_API_HOST', '127.0.0.1')
    api_port = int(os.getenv('CAROUSEL_API_PORT', '15001'))
    api_channel = grpc.insecure_channel(f'{api_host}:{api_port}')
    future = grpc.channel_ready_future(api_channel)
    try:
        future.result(timeout=int(os.getenv('CAROUSEL_API_TIMEOUT', '10')))
        logging.info("grpc channel for api is ready")
    except grpc.FutureTimeoutError:
        logging.error("error as timed out waiting for grpc channel for api")
    api_stub = APIStub(api_channel)

    # initialize grpc server
    server = grpc.aio.server()
    add_WorkerServicer_to_server(WorkerService(file_storage, api_stub), server)
    port = int(os.getenv('CAROUSEL_WORKER_PORT', '15002'))
    server.add_insecure_port(f'[::]:{port}')
    logging.info(f'starting grpc server on: {port}')
    await server.start()
    await server.wait_for_termination()

if __name__ == '__main__':
    logging.basicConfig(level=logging.getLevelName(
        os.getenv('CAROUSEL_LOG_LEVEL', 'INFO')))
    asyncio.run(serve())
