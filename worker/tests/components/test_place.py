"""Place Component Tests"""
from unittest.mock import patch
import pytest
from requests.exceptions import HTTPError

from src.components.place import API_URL, process_place


def test_place_success(requests_mock):
    requests_mock.get(url=API_URL.format(lat=19.2195856, lon=73.1056888), json={
        'address': {'city': 'Dombivali', 'state': 'Maharashtra', 'postcode': '421201', 'country': 'India'}
    })
    _ = process_place(None, 'mediaitem-user-id', 'mediaitem-id', [19.2195856, 73.1056888])
    # work(omkar): update tests

def test_place_failed_empty_response(requests_mock):
    requests_mock.get(url=API_URL.format(
        lat=19.2195856, lon=73.1056888), json={})
    _ = process_place(None, 'mediaitem-user-id', 'mediaitem-id', [19.2195856, 73.1056888])
    # work(omkar): update tests

def test_place_failed_exception(requests_mock):
    requests_mock.get(url=API_URL.format(
        lat=19.2195856, lon=73.1056888), json={})
    _ = process_place(None, 'mediaitem-user-id', 'mediaitem-id', [19.2195856, 73.1056888])
    # work(omkar): update tests
