"""Storage Utils Tests"""
from unittest.mock import patch
import pytest

from rawpy import LibRawNoThumbnailError

from src.store.utils import generate_thumbnail


def test_generate_thumbnail_failed_exception():
    with patch('rawpy.imread', side_effect=LibRawNoThumbnailError()):
        with pytest.raises(LibRawNoThumbnailError):
            thumbnail_bytes = generate_thumbnail('sample.jpeg')
            assert len(thumbnail_bytes) == 0
