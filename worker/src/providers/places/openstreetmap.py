"""Places: OpenStreetMap"""
import logging
import requests

from src.utils import getval_from_dict


class OpenStreetMap:
    """OpenStreetMap Places"""

    def __init__(self) -> None:
        self.url = 'https://nominatim.openstreetmap.org/reverse.php?zoom=18&format=jsonv2&lat={lat}&lon={lon}'
        self.timeout = 60

    def reverse_geocode(self, mediaitem_user_id: str, mediaitem_id: str, coordinates: list[float]) -> dict:
        """Reverse geocode location from co-ordinates"""
        url = self.url.format(lat=coordinates[0], lon=coordinates[1])
        res = requests.get(url=url, headers={'accept-language':'en-GB,en-US'}, timeout=self.timeout)
        res.raise_for_status()

        body = res.json()
        logging.debug(f'place for user {mediaitem_user_id} mediaitem {mediaitem_id}: {body}')

        if 'address' not in body:
            return None
        if 'address' in body and ('postcode' not in body['address'] or \
                                   'country' not in body['address'] or \
                                   'state' not in body['address']):
            return None

        address = body['address']
        return dict({
            'userId': mediaitem_user_id,
            'id': mediaitem_id,
            'postcode': getval_from_dict(address, ['postcode']),
            'country': getval_from_dict(address, ['country']),
            'state': getval_from_dict(address, ['state']),
            'city': getval_from_dict(address, ['city', 'county']),
            'town':  getval_from_dict(address, ['town']),
        })
