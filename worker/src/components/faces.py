"""Faces Component"""
import logging

from grpc import RpcError

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import MediaItemFacesRequest  # pylint: disable=no-name-in-module
from src.components.component import Component
from src.providers.faces.utils import init_faces


class Faces(Component):
    """Faces Component"""
    def __init__(self, api_stub: APIStub, source: str, params: list[str]) -> None:
        super().__init__('faces', api_stub)
        self.source = init_faces(source, params)

    async def process(self, mediaitem_user_id: str, mediaitem_id: str, _: str, metadata: dict) -> None:
        """Process faces detection for mediaitem"""
        if metadata is None or 'previewPath' not in metadata or ('type' in metadata and metadata['type'] == 'video'):
            return metadata
        try:
            result = self.source.detect(mediaitem_user_id, mediaitem_id, metadata['previewPath'])
            logging.debug(f'extracted faces for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')
            if result is not None:
                self._grpc_save_mediaitem_faces(result)
        except Exception as exp:
            logging.error(f'error getting faces response for user {mediaitem_user_id} '+
                          f'mediaitem {mediaitem_id}: {str(exp)}')
        logging.info(f'processed faces for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return metadata

    def _grpc_save_mediaitem_faces(self, result: dict):
        """gRPC call for saving mediaitem faces"""
        try:
            request = MediaItemFacesRequest(
                userId=result['userId'],
                id=result['id'],
                embeddings=result['embeddings']
            )
            _ = self.api_stub.SaveMediaItemFaces(request)
        except RpcError as rpc_exp:
            logging.error(
                f'error sending faces for mediaitem {request.id}: {str(rpc_exp)}')
