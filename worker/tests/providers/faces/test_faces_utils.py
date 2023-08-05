"""Providers Faces Utils Tests"""
from unittest import mock

from src.providers.faces.pytorch import PyTorchModule
from src.providers.faces.utils import init_faces


@mock.patch('src.providers.faces.PyTorchModule.__init__', return_value=None)
def test_init_faces_success(_):
    assert init_faces('tensorflow', []) == None
    assert isinstance(init_faces('pytorch', ['0.9', 'vggface2']), PyTorchModule)
