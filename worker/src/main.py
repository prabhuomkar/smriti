"""Worker"""
import asyncio
import logging
import os
import json

import grpc
from google.protobuf.empty_pb2 import Empty   # pylint: disable=no-name-in-module
from prometheus_client import start_http_server

from src.components import Component, Metadata, Places, Classification, OCR
from src.providers.search import init_search, PyTorchModule
from src.protos.api_pb2_grpc import APIStub
from src.protos.worker_pb2 import MediaItemProcessResponse, GenerateEmbeddingResponse  # pylint: disable=no-name-in-module
from src.protos.worker_pb2_grpc import WorkerServicer, add_WorkerServicer_to_server


class WorkerService(WorkerServicer):
    """Worker gRPC Service"""

    def __init__(self, components: list[Component], search_model: PyTorchModule) -> None:
        self.components = components
        self.search_model = search_model

    # pylint: disable=invalid-overridden-method
    async def MediaItemProcess(self, request, context) -> MediaItemProcessResponse:
        """MediaItem Process"""
        mediaitem_user_id = request.userId
        mediaitem_id = request.id
        mediaitem_file_path = request.filePath
        if mediaitem_id is not None and mediaitem_user_id is not None and mediaitem_file_path is not None:
            loop = asyncio.get_event_loop()
            loop.create_task(process_mediaitem(self.components, mediaitem_user_id, mediaitem_id, mediaitem_file_path))
            return MediaItemProcessResponse(ok=True)
        return MediaItemProcessResponse(ok=False)

    # pylint: disable=invalid-overridden-method
    async def GenerateEmbedding(self, request, context) -> GenerateEmbeddingResponse:
        """Generate Embedding"""
        text = request.text
        if self.search_model:
            result = self.search_model.generate_embedding(text)
            return GenerateEmbeddingResponse(embedding=result)
        return GenerateEmbeddingResponse(embedding=None)

# pylint: disable=redefined-builtin,invalid-name
async def process_mediaitem(components: list[Component], user_id: str, id: str, file_path: str) -> None:
    """Process mediaitem"""
    logging.info(f'started processing mediaitem for user {user_id} mediaitem {id}')
    metadata = await components[0].process(user_id, id, file_path, None)
    for i in range(1, len(components)):
        loop = asyncio.get_event_loop()
        loop.create_task(components[i].process(user_id, id, file_path, metadata))
    logging.info(f'finished processing mediaitem for user {user_id} mediaitem {id}')

async def serve() -> None:
    """Main serve function"""
    # start metrics
    start_http_server(int(os.getenv('SMRITI_METRICS_PORT', '5002')))

    # initialize grpc client
    api_host = os.getenv('SMRITI_API_HOST', '127.0.0.1')
    api_port = int(os.getenv('SMRITI_API_PORT', '15001'))
    api_channel = grpc.insecure_channel(f'{api_host}:{api_port}')
    future = grpc.channel_ready_future(api_channel)
    try:
        future.result(timeout=int(os.getenv('SMRITI_API_TIMEOUT', '30')))
        logging.info("grpc channel for api is ready")
    except grpc.FutureTimeoutError:
        logging.error("error as timed out waiting for grpc channel for api")
    api_stub = APIStub(api_channel)

    # get worker config
    worker_cfg = api_stub.GetWorkerConfig(Empty())
    cfg = json.loads(worker_cfg.config)
    logging.info(f'got worker configuration: {cfg}')

    # initialize components
    components = [Metadata(api_stub=api_stub)]
    search_model = None
    for item in cfg:
        if item['name'] == 'places':
            components.append(Places(api_stub=api_stub, source=item['source']))
        elif item['name'] == 'classification':
            components.append(Classification(api_stub=api_stub, source=item['source'], files=item['files']))
        elif item['name'] == 'ocr':
            components.append(OCR(api_stub=api_stub, source=item['source'], files=item['files']))
        elif item['name'] == 'search':
            search_model = init_search(name=item['source'], files=item['files'])

    # initialize grpc server
    server = grpc.aio.server()
    add_WorkerServicer_to_server(WorkerService(components, search_model), server)
    port = int(os.getenv('SMRITI_WORKER_PORT', '15002'))
    server.add_insecure_port(f'[::]:{port}')
    logging.info(f'starting grpc server on: {port}')
    await server.start()
    await server.wait_for_termination()

if __name__ == '__main__':
    logging.basicConfig(level=logging.getLevelName(os.getenv('SMRITI_LOG_LEVEL', 'INFO')))
    asyncio.run(serve())
