"""Place Component"""
import requests


API_URL = 'https://nominatim.openstreetmap.org/reverse.php?zoom=18&format=jsonv2&lat={lat}&lon={lon}'
API_TIMEOUT = 60

def process(lat: float, lon: float) -> dict:
    """Process will do reverse geocoding and return place details for given lat,lon"""
    url = API_URL.format(lat=lat, lon=lon)
    res = requests.get(url=url, timeout=API_TIMEOUT)
    res.raise_for_status()
    body = res.json()
    address = body['address'] if 'address' in body else {}
    return dict({
        'postcode': address['postcode'] if 'postcode' in address else None,
        'country': address['country'] if 'country' in address else None,
        'state': address['state'] if 'state' in address else None,
        'city': address['city'] if 'city' in address else None,
        'town': address['town'] if 'town' in address else None,
    })
