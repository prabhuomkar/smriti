"""Storage Utils"""
import os

from src.store.disk import Disk
from src.store.amazons3 import AmazonS3


def init_storage(name: str) -> AmazonS3 | Disk:
    """Initialize storage by name"""
    if name == 'amazons3':
        return AmazonS3(access_key=None, secret_key=None)
    return Disk(root=os.getenv('PENSIEVE_STORAGE_ROOT', '../storage'))
