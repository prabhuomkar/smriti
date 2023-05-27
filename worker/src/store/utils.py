"""Storage Utils"""
import os

from src.store.disk import Disk
from src.store.minio import MinIO


def init_storage(name: str) -> None | Disk | MinIO:
    """Initialize storage by name"""
    if name == 'amazons3':
        return None
    if name == 'minio':
        return MinIO(endpoint=os.getenv('SMRITI_STORAGE_MINIO_ENDPOINT'),
                     access_key=os.getenv('SMRITI_STORAGE_MINIO_ACCESS_KEY'),
                     secret_key=os.getenv('SMRITI_STORAGE_MINIO_SECRET_KEY'))
    return Disk(root=os.getenv('SMRITI_STORAGE_DISK_ROOT', '../storage'))
