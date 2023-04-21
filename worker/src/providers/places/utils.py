"""Places Utils"""
from src.providers.places.openstreetmap import OpenStreetMap


def init_places(name: str) -> None | OpenStreetMap:
    """Initialize places by name"""
    if name == 'googlemaps':
        return None
    return OpenStreetMap()
