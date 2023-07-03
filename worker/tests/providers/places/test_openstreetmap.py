"""Providers Places OpenStreetMap Tests"""
from unittest.mock import patch
import pytest
from requests.exceptions import HTTPError

from src.providers.places import init_places


API_URL = 'https://nominatim.openstreetmap.org/reverse.php?zoom=18&format=jsonv2&lat={lat}&lon={lon}'

def test_openstreetmap_success(requests_mock):
    requests_mock.get(url=API_URL.format(lat=19.2195856, lon=73.1056888), json={
        'address': {'city': 'Dombivali', 'state': 'Maharashtra', 'postcode': '421201', 'country': 'India'}
    })
    result = init_places('openstreetmap').reverse_geocode('mediaitem_user_id', 'mediaitem_id', [19.2195856, 73.1056888])
    assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id',
                      'postcode': '421201', 'country': 'India',
                      'state': 'Maharashtra', 'city': 'Dombivali', 'town': None}


def test_openstreetmap_failed_empty_response(requests_mock):
    requests_mock.get(url=API_URL.format(
        lat=19.2195856, lon=73.1056888), json={})
    result = init_places('openstreetmap').reverse_geocode('mediaitem_user_id', 'mediaitem_id', [19.2195856, 73.1056888])
    assert result == None


def test_openstreetmap_failed_exception(requests_mock):
    requests_mock.get(url=API_URL.format(
        lat=19.2195856, lon=73.1056888), json={})
    with patch('requests.get', side_effect=HTTPError()):
        with pytest.raises(HTTPError):
            result = init_places('openstreetmap').reverse_geocode('mediaitem_user_id', 'mediaitem_id', [19.2195856, 73.1056888])
            assert result == None
