"""Storage Utils"""
import os

from store.disk import Disk
from store.amazons3 import AmazonS3


def init_storage(name: str) -> AmazonS3 | Disk:
    """Initialize storage by name"""
    if name == 'amazons3':
        return AmazonS3()
    return Disk(root=os.getenv('PENSIEVE_WORKER_STORAGE_ROOT', '/storage'))
