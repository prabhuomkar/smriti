"""Places Component"""
import logging

from grpc import RpcError

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import MediaItemPlaceRequest  # pylint: disable=no-name-in-module
from src.components.component import Component
from src.providers.places.utils import init_places


class Places(Component):
    """Places Component"""
    def __init__(self, api_stub: APIStub, source: str) -> None:
        super().__init__('places', None, api_stub)
        self.source = init_places(source)

    async def process(self, mediaitem_user_id: str, mediaitem_id: str, metadata: dict) -> None:
        """Process place details from latitude and longitude"""
        coordinates = [metadata['latitude'], metadata['longitude']]
        try:
            result = self.source.reverse_geocode(mediaitem_user_id, mediaitem_id, coordinates)

            logging.debug(f'extracted place for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')

            self._grpc_save_mediaitem_place(result)
        except Exception as exp:
            logging.error(f'error getting place response for user {mediaitem_user_id} '+
                          f'mediaitem {mediaitem_id}: {str(exp)}')

        logging.info(f'processed place for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return None


    def _grpc_save_mediaitem_place(self, result: dict):
        """gRPC call for saving mediaitem place"""
        try:
            request = MediaItemPlaceRequest(
                userId=result['userId'],
                id=result['id'],
                postcode=result['postcode'] if 'postcode' in result else None,
                country=result['country'] if 'country' in result else None,
                state=result['state'] if 'state' in result else None,
                city=result['city'] if 'city' in result else None,
                town=result['town'] if 'town' in result else None,
            )
            _ = self.api_stub.SaveMediaItemPlace(request)
        except RpcError as rpc_exp:
            logging.error(
                f'error sending place for mediaitem {request.id}: {str(rpc_exp)}')
