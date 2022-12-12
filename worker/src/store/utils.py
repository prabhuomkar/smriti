"""Storage Utils"""
import os

from .disk import Disk
from .amazons3 import AmazonS3


def init_storage(name: str) -> AmazonS3 | Disk:
    """Initialize storage by name"""
    if name == 'amazons3':
        return AmazonS3()
    return Disk(root=os.getenv('PENSIEVE_STORAGE_ROOT', '/storage'))
