"""Place Component"""
import logging
import requests

from src.components.utils import getval_from_dict
from src.components.utils import grpc_save_mediaitem_place
from src.protos.api_pb2_grpc import APIStub

API_URL = 'https://nominatim.openstreetmap.org/reverse.php?zoom=18&format=jsonv2&lat={lat}&lon={lon}'
API_TIMEOUT = 60


async def process_place(api_stub: APIStub, mediaitem_user_id: str, mediaitem_id: str, coordinates: list[float]) -> None:
    """Process place details from latitude and longitude"""
    try:
        url = API_URL.format(lat=coordinates[0], lon=coordinates[1])
        res = requests.get(url=url, timeout=API_TIMEOUT)
        res.raise_for_status()

        body = res.json()
        logging.debug(f'place for user {mediaitem_user_id} mediaitem {mediaitem_id}: {body}')

        address = body['address'] if 'address' in body else {}
        result = {
            'userId': mediaitem_user_id,
            'id': mediaitem_id,
            'postcode': getval_from_dict(address, ['postcode']),
            'country': getval_from_dict(address, ['country']),
            'state': getval_from_dict(address, ['state']),
            'city': getval_from_dict(address, ['city']),
            'town':  getval_from_dict(address, ['town']),
        }

        logging.debug(f'extracted place for user {mediaitem_user_id} mediaitem {mediaitem_id}: {result}')

        grpc_save_mediaitem_place(api_stub, result)
    except Exception as e:
        logging.error(f'error getting place response for user {mediaitem_user_id} mediaitem {mediaitem_id}: {str(e)}')

    logging.info(f'processed place for user {mediaitem_user_id} mediaitem {mediaitem_id}')
    return None
