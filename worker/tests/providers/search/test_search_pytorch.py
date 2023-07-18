"""Providers Search PyTorch Tests"""
from unittest import mock
import pytest

import torch
from transformers import PreTrainedTokenizerFast

from src.providers.search.pytorch import PyTorchModule


@mock.patch('torch.jit.load', return_value=torch.jit.ScriptModule())
@mock.patch('transformers.AutoTokenizer.from_pretrained', return_value=None)
def test_pytorch_success(_, __):
    pytorch_module = PyTorchModule(['model_name.pt','tokenizer_path'])
    pytorch_module.module = mock.MagicMock()
    pytorch_module.module.forward.return_value = torch.tensor([0.302,0.92,0.38])
    pytorch_module.tokenizer = mock.MagicMock()
    pytorch_module.tokenizer.tokenize.return_value = torch.rand(1, 1)
    result = pytorch_module.generate_embedding('search query text')
    assert result == torch.tensor([0.302,0.92,0.38]).tolist()

@mock.patch('torch.jit.load', return_value=torch.jit.ScriptModule())
@mock.patch('transformers.AutoTokenizer.from_pretrained', return_value=None)
def test_pytorch_failed_empty_response(_, __):
    pytorch_module = PyTorchModule(['model_name.pt','tokenizer_path'])
    pytorch_module.module = mock.MagicMock()
    pytorch_module.module.forward.return_value = None
    pytorch_module.tokenizer = mock.MagicMock()
    pytorch_module.tokenizer.tokenize.return_value = torch.rand(1, 1)
    result = pytorch_module.generate_embedding('search query text')
    assert result == None

@mock.patch('torch.jit.load', return_value=torch.jit.ScriptModule())
@mock.patch('transformers.AutoTokenizer.from_pretrained', return_value=None)
def test_pytorch_failed_exception(_, __):
    pytorch_module = PyTorchModule(['model_name.pt','tokenizer_path'])
    pytorch_module.module = mock.MagicMock()
    pytorch_module.module.forward.side_effect = Exception('some error')
    pytorch_module.tokenizer = mock.MagicMock()
    pytorch_module.tokenizer.tokenize.return_value = torch.rand(1, 1)
    with pytest.raises(Exception):
        result = pytorch_module.generate_embedding('search query text')
        assert result == None
