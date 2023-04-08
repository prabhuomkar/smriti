"""Storage Disk Tests"""
import os
from unittest.mock import patch

from src.store.disk import Disk

@patch('os.mkdir')
@patch('builtins.open')
@patch('os.remove')
def test_disk_success(mock_mkdir, mock_open, mock_remove):
    f = mock_mkdir.return_value
    f.method.return_value = True
    f = mock_open.return_value
    f.method.return_value = None
    f = mock_remove.return_value
    f.method.return_value = None
    disk = Disk(root='.')
    assert disk.upload('mediaitem_id', None, 'previews') == os.path.abspath('./previews/mediaitem_id')
    assert disk.get('mediaitem_id', 'thumbnails') == os.path.abspath('./thumbnails/mediaitem_id')
    assert disk.delete('') == None
