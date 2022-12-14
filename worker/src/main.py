"""Worker"""
import asyncio
import logging
import os
from typing import AsyncIterable, Iterable

import grpc

from store import init_storage
from components import process_metadata
from api_pb2_grpc import APIStub
from worker_pb2 import MediaItemProcessRequest, MediaItemProcessResponse  # pylint: disable=no-name-in-module
from worker_pb2_grpc import WorkerServicer, add_WorkerServicer_to_server


class WorkerService(WorkerServicer):
    """Worker gRPC Service"""

    def __init__(self, file_storage, api_stub: APIStub) -> None:
        self.file_storage = file_storage
        self.api_stub = api_stub

    async def MediaItemProcess(self, request_iterator: AsyncIterable[
            MediaItemProcessRequest], unused_context) -> MediaItemProcessResponse:
        """MediaItem Process"""
        mediaitem_id = None
        mediaitem_command = None
        async for mediaitem in request_iterator:
            try:
                self.file_storage.upload(id=mediaitem.id,
                                         offset=mediaitem.offset,
                                         content=mediaitem.content)
                mediaitem_id = mediaitem.id
                mediaitem_command = mediaitem.command
            except Exception as e:
                logging.error(f'error processing mediaitem for storage: {str(e)}', {
                              'id': mediaitem.id, 'offset': mediaitem.offset})
                mediaitem_id = None
                return MediaItemProcessResponse(ok=False)
        if mediaitem_id is not None and 'finish' in mediaitem_command:
            loop = asyncio.get_event_loop()
            loop.create_task(process_metadata(
                self.file_storage, self.api_stub, mediaitem_id))
        return MediaItemProcessResponse(ok=True)


async def serve() -> None:
    """Main serve function"""
    # initialize storage
    file_storage = init_storage(os.getenv('PENSIEVE_STORAGE', 'disk'))

    # initialize grpc client
    api_host = os.getenv('PENSIEVE_API_HOST', 'localhost')
    api_port = int(os.getenv('PENSIEVE_API_PORT', '15001'))
    api_channel = grpc.insecure_channel(f'{api_host}:{api_port}')
    api_stub = APIStub(api_channel)

    # initialize grpc server
    server = grpc.aio.server()
    add_WorkerServicer_to_server(WorkerService(file_storage, api_stub), server)
    port = int(os.getenv('PENSIEVE_WORKER_PORT', '15002'))
    server.add_insecure_port(f'[::]:{port}')
    logging.info(f'starting grpc server on: {port}')
    await server.start()
    await server.wait_for_termination()

if __name__ == '__main__':
    logging.basicConfig(level=logging.getLevelName(
        os.getenv('PENSIEVE_LOG_LEVEL', 'INFO')))
    asyncio.run(serve())
