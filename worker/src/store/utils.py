"""Storage Utils"""
import os
import io

import rawpy
from PIL import Image

from .disk import Disk
from .amazons3 import AmazonS3


def init_storage(name: str) -> AmazonS3 | Disk:
    """Initialize storage by name"""
    if name == 'amazons3':
        return AmazonS3()
    return Disk(root=os.getenv('PENSIEVE_STORAGE_ROOT', '/storage'))


def generate_thumbnail(original_file_path: str) -> bytes:
    """Generate thumbnail"""
    with rawpy.imread(original_file_path) as raw:
        thumb = raw.extract_thumb()
    if thumb.format == rawpy.ThumbFormat.JPEG:
        return thumb.data
    elif thumb.format == rawpy.ThumbFormat.BITMAP:
        img = Image.frombytes(data=thumb.data)
        img_bytes = io.BytesIO()
        img.save(img_bytes, format='JPEG')
        return img_bytes.getvalue()
