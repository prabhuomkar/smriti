"""OCR Component"""
import logging

from grpc import RpcError

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import MediaItemMLResultRequest  # pylint: disable=no-name-in-module
from src.components.component import Component
from src.providers.ocr.utils import init_ocr


class OCR(Component):
    """OCR Component"""
    def __init__(self, api_stub: APIStub, source: str, files: list[str]) -> None:
        super().__init__('ocr', api_stub)
        self.model = init_ocr(source, files)

    async def process(self, mediaitem_user_id: str, mediaitem_id: str, _: str, metadata: dict) -> None:
        """Process ocr from mediaitem"""
        if metadata is None or 'previewPath' not in metadata or ('type' in metadata and metadata['type'] == 'video'):
            return None
        try:
            result = self.model.extract(mediaitem_user_id, mediaitem_id, metadata['previewPath'])

            logging.debug(f'extracted ocr for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')

            if result is not None:
                self._grpc_save_mediaitem_ml_result(result)
        except Exception as exp:
            logging.error(f'error getting thing response for user {mediaitem_user_id} '+
                          f'mediaitem {mediaitem_id}: {str(exp)}')

        logging.info(f'processed ocr for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return None


    def _grpc_save_mediaitem_ml_result(self, result: dict):
        """gRPC call for saving mediaitem text"""
        try:
            request = MediaItemMLResultRequest(
                userId=result['userId'],
                id=result['id'],
                name='ocr',
                value=result['words'] if result else None
            )
            _ = self.api_stub.SaveMediaItemMLResult(request)
        except RpcError as rpc_exp:
            logging.error(
                f'error sending ml result for mediaitem {request.id}: {str(rpc_exp)}')
