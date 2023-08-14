"""Providers Classification Utils Tests"""
from unittest import mock

from src.providers.classification.pytorch import PyTorchModule
from src.providers.classification.utils import init_classification


@mock.patch('torch.jit.load', return_value=None)
def test_init_classification_success(_):
    assert init_classification('tensorflow', []) == None
    assert isinstance(init_classification('pytorch', {'file':'model_name.pt'}), PyTorchModule)
