"""Providers Search Utils Tests"""
from unittest import mock

from src.providers.search.pytorch import PyTorchModule
from src.providers.search.utils import init_search


@mock.patch('src.providers.search.PyTorchModule.__init__', return_value=None)
def test_init_search_success(_):
    assert init_search('tensorflow', []) == None
    assert isinstance(init_search('pytorch', ['/det_models', '/rec_models', '/cls_models']), PyTorchModule)
