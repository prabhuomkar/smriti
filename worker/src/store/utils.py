"""Storage Utils"""
import os

from src.store.disk import Disk


def init_storage(name: str) -> None | Disk:
    """Initialize storage by name"""
    if name == 'amazons3':
        return None
    return Disk(root=os.getenv('CAROUSEL_STORAGE_ROOT', '../storage'))
