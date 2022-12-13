"""Components Utils"""
import os


def init_components():
    """Initialize components"""
    result = []
    if bool(os.getenv('PENSIEVE_FEATURE_EXPLORE_PLACES', 'True')):
        result += []
    if bool(os.getenv('PENSIEVE_FEATURE_EXPLORE_THINGS', 'True')):
        result += []
    if bool(os.getenv('PENSIEVE_FEATURE_EXPLORE_PEOPLE', 'True')):
        result += []
    return result
