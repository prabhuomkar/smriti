"""Storage Utils Tests"""
from unittest.mock import patch

from src.store.disk import Disk
from src.store.utils import init_storage


@patch('os.path.exists')
def test_init_storage_success(mock_exists):
    f = mock_exists.return_value
    f.method.return_value = True
    assert init_storage('amazons3') == None
    assert isinstance(init_storage('disk'), Disk)
