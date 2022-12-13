"""Worker"""
import asyncio
import logging
import os

import grpc

from store import init_storage
from components import process_metadata
from api_pb2 import MediaItemMetadataRequest  # pylint: disable=no-name-in-module
from api_pb2_grpc import APIStub
from worker_pb2 import MediaItemProcessResponse  # pylint: disable=no-name-in-module
from worker_pb2_grpc import WorkerServicer, add_WorkerServicer_to_server


class WorkerService(WorkerServicer):
    """Worker gRPC Service"""

    def __init__(self, file_storage, api_stub: APIStub) -> None:
        self.file_storage = file_storage
        self.api_stub = api_stub

    def MediaItemProcess(self, request_iterator, context):
        """MediaItem Process"""
        mediaitem_id = None
        mediaitem_command = None
        for mediaitem in request_iterator:
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
            result = process_metadata(
                storage=self.file_storage, id=mediaitem_id)
            logging.debug(result)
            _ = self.api_stub.SaveMediaItemMetadata(MediaItemMetadataRequest(
                id=result['id'],
                status=result['status'],
                mimeType=result['mimeType'] if 'mimeType' in result else None,
                sourceUrl=result['sourceUrl'] if 'sourceUrl' in result else None,
                previewUrl=result['previewUrl'] if 'previewUrl' in result else None,
                thumbnailUrl=result['thumbnailUrl'] if 'thumbnailUrl' in result else None,
                type=result['type'] if 'type' in result else None,
                width=result['width'] if 'width' in result else None,
                height=result['height'] if 'height' in result else None,
                creationTime=result['creationTime'] if 'creationTime' in result else None,
                cameraMake=result['cameraMake'] if 'cameraMake' in result else None,
                cameraModel=result['cameraModel'] if 'cameraModel' in result else None,
                focalLength=result['focalLength'] if 'focalLength' in result else None,
                apertureFNumber=result['apertureFNumber'] if 'apertureFNumber' in result else None,
                isoEquivalent=result['isoEquivalent'] if 'isoEquivalent' in result else None,
                exposureTime=result['exposureTime'] if 'exposureTime' in result else None,
                fps=result['fps'] if 'fps' in result else None,
                latitude=result['latitude'] if 'latitude' in result else None,
                longitude=result['longitude'] if 'longitude' in result else None,
            ))
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
