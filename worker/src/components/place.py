"""Place Component"""
import requests

from src.components.utils import getval_from_dict

API_URL = 'https://nominatim.openstreetmap.org/reverse.php?zoom=18&format=jsonv2&lat={lat}&lon={lon}'
API_TIMEOUT = 60


def process_place(lat: float, lon: float) -> dict:
    """Process place details from latitude and longitude"""
    url = API_URL.format(lat=lat, lon=lon)
    res = requests.get(url=url, timeout=API_TIMEOUT)
    res.raise_for_status()
    body = res.json()
    address = body['address'] if 'address' in body else {}
    return dict({
        'postcode': getval_from_dict(address, ['postcode']),
        'country': getval_from_dict(address, ['country']),
        'state': getval_from_dict(address, ['state']),
        'city': getval_from_dict(address, ['city']),
        'town':  getval_from_dict(address, ['town']),
    })
