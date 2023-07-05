"""Providers Things Utils Tests"""
from unittest import mock

from src.providers.things.pytorch import PyTorchModule
from src.providers.things.utils import init_things


@mock.patch('torch.jit.load', return_value=None)
def test_init_things_success(_):
    assert init_things('tensorflow', []) == None
    assert isinstance(init_things('pytorch', ['model_name.pt']), PyTorchModule)
