"""Providers Places Utils Tests"""
from src.providers.places.openstreetmap import OpenStreetMap
from src.providers.places.utils import init_places


def test_init_places_success():
    assert init_places('googlemaps') == None
    assert isinstance(init_places('openstreetmap'), OpenStreetMap)
