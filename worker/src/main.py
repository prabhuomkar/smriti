"""Worker"""
import asyncio
import logging
import os
import json
import signal

import grpc
import schedule
from google.protobuf.empty_pb2 import Empty  # pylint: disable=no-name-in-module
from prometheus_client import start_http_server

from src.components.component import Component
from src.components.finalize import Finalize
from src.providers.search import init_search, PyTorchModule
from src.protos.api_pb2_grpc import APIStub
from src.protos.worker_pb2 import (  # pylint: disable=no-name-in-module
    MediaItemComponent, MediaItemProcessResponse, GenerateEmbeddingResponse,
    METADATA, PREVIEW_THUMBNAIL, CLASSIFICATION, FACES, OCR, PLACES, SEARCH)
from src.protos.worker_pb2_grpc import WorkerServicer, add_WorkerServicer_to_server


class WorkerService(WorkerServicer):
    """Worker gRPC Service"""

    def __init__(self, components: list[Component], search_model: PyTorchModule) -> None:
        self.ordered_components = components
        self.search_model = search_model

    # pylint: disable=invalid-overridden-method
    async def MediaItemProcess(self, request, context) -> MediaItemProcessResponse:
        """MediaItem Process"""
        mediaitem_user_id = request.userId
        mediaitem_id = request.id
        mediaitem_file_path = request.filePath
        mediaitem_components = request.components
        mediaitem_payload = request.payload
        logging.info(f'mediaitem process request user {mediaitem_user_id} id {mediaitem_id} \
                     path {mediaitem_file_path} components {mediaitem_components} payload {mediaitem_payload}')
        if mediaitem_id is not None and mediaitem_user_id is not None \
            and mediaitem_file_path is not None and mediaitem_components is not None:
            components = []
            mediaitem_components = [
                MediaItemComponent.Name(mediaitem_component) for mediaitem_component in mediaitem_components]
            for _component in self.ordered_components:
                if _component.name in mediaitem_components or _component.name == 'FINALIZE':
                    components.append(_component)
            logging.info([_comp.name for _comp in components])
            loop = asyncio.get_event_loop()
            loop.create_task(process_mediaitem(components,self.search_model if MediaItemComponent.Name(SEARCH) \
                                               in mediaitem_components else None, mediaitem_user_id, mediaitem_id,
                                               mediaitem_file_path, mediaitem_payload))
            return MediaItemProcessResponse(ok=True)
        return MediaItemProcessResponse(ok=False)

    # pylint: disable=invalid-overridden-method
    async def GenerateEmbedding(self, request, context) -> GenerateEmbeddingResponse:
        """Generate Embedding"""
        text = request.text
        if self.search_model:
            result = self.search_model.generate_embedding('text', text)
            return GenerateEmbeddingResponse(embedding=result)
        return GenerateEmbeddingResponse(embedding=None)

def run_pending() -> None:
    """Run scheduled jobs in background"""
    while True:
        schedule.run_pending()
        asyncio.run(asyncio.sleep(1))

# pylint: disable=redefined-builtin,invalid-name,too-many-arguments
async def process_mediaitem(components: list[Component], search_model: PyTorchModule,
                            user_id: str, id: str, file_path: str, payload: any) -> None:
    """Process mediaitem"""
    logging.info(f'started processing mediaitem for user {user_id} mediaitem {id}')
    result = dict(payload)
    for i in range(len(components)-1):
        loop = asyncio.get_event_loop()
        task = loop.create_task(components[i].process(user_id, id, file_path, result))
        result = await task
    if search_model:
        result['embeddings'] = search_model.generate_embedding('file', result)
        if 'keywords' in result:
            result['embeddings'] += [search_model.generate_embedding('text', result['keywords'])]
    await components[-1].process(user_id, id, file_path, result)
    logging.info(f'finished processing mediaitem for user {user_id} mediaitem {id}')

async def serve() -> None: # pylint: disable=too-many-locals
    """Start gRPC and metrics server"""
    # start metrics
    metrics_server, metrics_thread = start_http_server(int(os.getenv('SMRITI_METRICS_PORT', '5002')))

    # initialize api grpc client
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
    components = []
    search_model = None
    for item in cfg:
        if 'params' in item:
            item['params'] = json.loads(item['params'])
        if item['name'] == MediaItemComponent.Name(METADATA):
            from src.components.metadata import Metadata
            components.append(Metadata(api_stub=api_stub))
        elif item['name'] == MediaItemComponent.Name(PREVIEW_THUMBNAIL):
            from src.components.preview_thumbnail import PreviewThumbnail
            components.append(PreviewThumbnail(api_stub=api_stub, params=item['params']))
        elif item['name'] == MediaItemComponent.Name(PLACES):
            from src.components.places import Places
            components.append(Places(api_stub=api_stub, source=item['source']))
        elif item['name'] == MediaItemComponent.Name(CLASSIFICATION):
            from src.components.classification import Classification
            components.append(Classification(api_stub=api_stub, source=item['source'], params=item['params']))
        elif item['name'] == MediaItemComponent.Name(OCR):
            from src.components.ocr import OCR as OCRComponent
            components.append(OCRComponent(api_stub=api_stub, source=item['source'], params=item['params']))
        elif item['name'] == MediaItemComponent.Name(FACES):
            from src.components.faces import Faces
            components.append(Faces(api_stub=api_stub, source=item['source'], params=item['params']))
        elif item['name'] == MediaItemComponent.Name(SEARCH):
            search_model = init_search(name=item['source'], params=item['params'])
    components.append(Finalize(api_stub=api_stub))

    # initialize worker grpc server
    grpc_server = grpc.aio.server(maximum_concurrent_rpcs=int(os.getenv('SMRITI_WORKER_CONCURRENT_RPCS', '5')))
    add_WorkerServicer_to_server(WorkerService(components, search_model), grpc_server)
    port = int(os.getenv('SMRITI_WORKER_PORT', '15002'))
    grpc_server.add_insecure_port(f'[::]:{port}')
    await grpc_server.start()
    logging.info(f'starting grpc server on: {port}')
    return grpc_server, metrics_server, metrics_thread

async def shutdown(grpc_server, metrics_server, metrics_thread):
    """Shutdown gRPC and metrics server"""
    logging.info('stopping grpc server')
    await grpc_server.stop(10)
    logging.info('stopping metrics server')
    metrics_server.shutdown()
    metrics_thread.join()

async def main():
    """Main driver function"""
    grpc_server, metrics_server, metrics_thread = await serve()
    loop = asyncio.get_running_loop()
    loop.run_in_executor(None, run_pending)
    for sig in [signal.SIGINT, signal.SIGTERM]:
        loop.add_signal_handler(sig, lambda sig=sig: asyncio.create_task(
            shutdown(grpc_server, metrics_server, metrics_thread)))
    await grpc_server.wait_for_termination()


if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s %(levelname)s %(message)s',
                        datefmt='%Y/%m/%d %H:%M:%S',
                        level=logging.getLevelName(os.getenv('SMRITI_LOG_LEVEL', 'INFO')))
    asyncio.run(main())
