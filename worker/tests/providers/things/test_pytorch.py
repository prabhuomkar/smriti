"""Providers Things PyTorch Tests"""
from unittest import mock
import pytest

import numpy as np
import torch

from src.providers.things.pytorch import PyTorchModule


@mock.patch('torch.jit.load', return_value=torch.jit.ScriptModule())
@mock.patch('PIL.Image.open', return_value=np.ndarray((1, 1)))
def test_pytorch_success(_, __):
    pytorch_module = PyTorchModule(['model_name.pt'])
    pytorch_module.module = mock.MagicMock()
    pytorch_module.module.forward.return_value = {'pizza':0.7939030,'burger':0.1933852}
    result = pytorch_module.classify('mediaitem_user_id', 'mediaitem_id', 'previews/file_name')
    assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'name': 'pizza'}

@mock.patch('torch.jit.load', return_value=torch.jit.ScriptModule())
@mock.patch('PIL.Image.open', return_value=np.ndarray((1, 1)))
def test_pytorch_failed_empty_response(_, __):
    pytorch_module = PyTorchModule(['model_name.pt'])
    pytorch_module.module = mock.MagicMock()
    pytorch_module.module.forward.return_value = {}
    result = pytorch_module.classify('mediaitem_user_id', 'mediaitem_id', 'previews/file_name')
    assert result == None

@mock.patch('torch.jit.load', return_value=torch.jit.ScriptModule())
@mock.patch('PIL.Image.open', return_value=np.ndarray((1, 1)))
def test_pytorch_failed_exception(_, __):
    pytorch_module = PyTorchModule(['model_name.pt'])
    pytorch_module.module = mock.MagicMock()
    pytorch_module.module.forward.side_effect = Exception('some error')
    with pytest.raises(Exception):
        result = pytorch_module.classify('mediaitem_user_id', 'mediaitem_id', 'previews/file_name')
        assert result == None
