"""Places Component"""
import logging

from grpc import RpcError

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import MediaItemPlaceRequest  # pylint: disable=no-name-in-module
from src.protos.worker_pb2 import MediaItemComponent, PLACES  # pylint: disable=no-name-in-module
from src.components.component import Component
from src.providers.places.utils import init_places


class Places(Component):
    """Places Component"""
    def __init__(self, api_stub: APIStub, source: str) -> None:
        super().__init__(MediaItemComponent.Name(PLACES), api_stub)
        self.source = init_places(source)

    async def process(self, mediaitem_user_id: str, mediaitem_id: str, _: str, metadata: dict) -> None:
        """Process place details from latitude and longitude"""
        if metadata is None or ('latitude' not in metadata or 'longitude' not in metadata) or \
            (metadata['latitude'] is None and metadata['longitude'] is None):
            return metadata
        coordinates = [metadata['latitude'], metadata['longitude']]
        try:
            result = self.source.reverse_geocode(mediaitem_user_id, mediaitem_id, coordinates)
            logging.debug(f'extracted place for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')
            if result is not None:
                place_keywords = ''
                if result['postcode']:
                    place_keywords += (result['postcode'].lower()+' ')
                if result['city']:
                    place_keywords += (result['city'].lower()+' ')
                if result['town']:
                    place_keywords += (result['town'].lower()+' ')
                if result['state']:
                    place_keywords += (result['state'].lower()+' ')
                if result['country']:
                    place_keywords += (result['country'].lower()+' ')
                place_keywords = place_keywords.strip()
                if 'keywords' not in metadata or metadata['keywords'] == '':
                    metadata['keywords'] = place_keywords
                else:
                    for keyword in place_keywords.split(' '):
                        if keyword not in metadata['keywords']:
                            metadata['keywords'] += (' ' + keyword)
                metadata['keywords'] = metadata['keywords'].strip()
                self._grpc_save_mediaitem_place(result)
        except Exception as exp:
            logging.error('error getting place response for user %s mediaitem %s: %s',
                          mediaitem_user_id, mediaitem_id, exp)
        logging.info(f'processed place for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        return metadata

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
