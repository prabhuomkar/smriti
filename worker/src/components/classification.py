"""Classification Component"""
import logging

from grpc import RpcError

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import MediaItemThingRequest  # pylint: disable=no-name-in-module
from src.components.component import Component
from src.providers.classification.utils import init_classification


class Classification(Component):
    """Classification Component"""
    def __init__(self, api_stub: APIStub, source: str, params: list[str]) -> None:
        super().__init__('classification', api_stub)
        self.model = init_classification(source, params)

    async def process(self, mediaitem_user_id: str, mediaitem_id: str, _: str, metadata: dict) -> None:
        """Process classification from mediaitem"""
        if metadata is None or 'previewPath' not in metadata or ('type' in metadata and metadata['type'] == 'video'):
            return None
        try:
            result = self.model.classify(mediaitem_user_id, mediaitem_id, metadata['previewPath'])

            logging.debug(f'extracted classification for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')

            if result is not None:
                self._grpc_save_mediaitem_thing(result)
        except Exception as exp:
            logging.error(f'error getting thing response for user {mediaitem_user_id} '+
                          f'mediaitem {mediaitem_id}: {str(exp)}')

        logging.info(f'processed classification for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return None


    def _grpc_save_mediaitem_thing(self, result: dict):
        """gRPC call for saving mediaitem thing"""
        try:
            request = MediaItemThingRequest(
                userId=result['userId'],
                id=result['id'],
                name=result['name'].title() if 'name' in result else None
            )
            _ = self.api_stub.SaveMediaItemThing(request)
        except RpcError as rpc_exp:
            logging.error(
                f'error sending thing for mediaitem {request.id}: {str(rpc_exp)}')
