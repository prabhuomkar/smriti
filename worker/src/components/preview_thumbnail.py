"""PreviewThumbnail Component"""
import logging

from grpc import RpcError

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import MediaItemMetadataRequest  # pylint: disable=no-name-in-module
from src.components.component import Component


class PreviewThumbnail(Component):
    """PreviewThumbnail Component"""

    def __init__(self, api_stub: APIStub) -> None:
        super().__init__('preview_thumbnail', api_stub)

    async def process(self, mediaitem_user_id: str, mediaitem_id: str, _: str, metadata: dict) -> dict:
        """Process to generate preview and thumbnail for mediaitem"""
        logging.debug(f'generating preview and thumbnail of mediaitem for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        result = {}
        result['userId'] = mediaitem_user_id
        result['id'] = mediaitem_id
        result['status'] = 'READY'
        self._grpc_save_mediaitem_preview_thumbnail(result)
        logging.debug(f'generated preview and thumbnail of mediaitem for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return None

    def _grpc_save_mediaitem_preview_thumbnail(self, result: dict):
        """gRPC call for saving mediaitem metadata"""
        try:
            request = MediaItemMetadataRequest(
                userId=result['userId'],
                id=result['id'],
                status=result['status'],
                placeholder=result['placeholder'] if 'placeholder' in result else None,
                sourcePath=result['sourcePath'] if 'sourcePath' in result else None,
                previewPath=result['previewPath'] if 'previewPath' in result else None,
                thumbnailPath=result['thumbnailPath'] if 'thumbnailPath' in result else None,
            )
            _ = self.api_stub.SaveMediaItemMetadata(request)
        except RpcError as rpc_exp:
            logging.error(
                f'error sending preview and thumbnail for mediaitem {request.id}: {str(rpc_exp)}')
