"""Worker"""
import asyncio
import logging
import os

import grpc

from store import init_storage
from components import init_components, process_metadata
from worker_pb2 import MediaItemProcessResponse  # pylint: disable=no-name-in-module
from worker_pb2_grpc import WorkerServicer, add_WorkerServicer_to_server


class WorkerService(WorkerServicer):
    """Worker gRPC Service"""

    def __init__(self) -> None:
        self.file_storage = init_storage(os.getenv('PENSIEVE_STORAGE', 'disk'))
        self.pipeline = init_components()

    def MediaItemProcess(self, request_iterator, context):
        """MediaItem Process"""
        mediaitem_id = None
        for mediaitem in request_iterator:
            try:
                self.file_storage.upload(id=mediaitem.id,
                                         offset=mediaitem.offset, content=mediaitem.content)
                mediaitem_id = mediaitem.id
            except Exception as e:
                logging.error(f'error processing mediaitem for storage: {str(e)}', {
                              'id': mediaitem.id, 'offset': mediaitem.offset})
                return MediaItemProcessResponse(ok=False)
        if mediaitem_id is not None:
            result = process_metadata(
                storage=self.file_storage, id=mediaitem.id)
            logging.info(result)
            return MediaItemProcessResponse(ok=True)
        return MediaItemProcessResponse(ok=False)


async def serve() -> None:
    """Main serve function"""
    server = grpc.aio.server()
    add_WorkerServicer_to_server(WorkerService(), server)
    port = int(os.getenv('PENSIEVE_WORKER_PORT', '15002'))
    server.add_insecure_port(f'[::]:{port}')
    logging.info(f'starting grpc server on: {port}')
    await server.start()
    await server.wait_for_termination()

if __name__ == '__main__':
    logging.basicConfig(level=logging.getLevelName(
        os.getenv('PENSIEVE_LOG_LEVEL', 'INFO')))
    asyncio.run(serve())
