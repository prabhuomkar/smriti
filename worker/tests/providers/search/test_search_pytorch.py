"""Providers Search PyTorch Tests"""
from unittest import mock
import pytest

import numpy as np
import torch

from src.providers.search.pytorch import PyTorchModule


@mock.patch('torch.jit.load', return_value=torch.jit.ScriptModule())
@mock.patch('transformers.AutoTokenizer.from_pretrained', return_value=None)
@mock.patch('transformers.AutoImageProcessor.from_pretrained', return_value=None)
def test_pytorch_text_success(_, __, ___):
    pytorch_module = PyTorchModule({'tokenizer_dir':'tokenizer_path','processor_dir':'processor_path','text_file':'text_model.pt','vision_file':'vision_model.pt'})
    pytorch_module.text_module = mock.MagicMock()
    pytorch_module.text_module.forward.return_value = torch.tensor([0.302,0.92,0.38])
    pytorch_module.vision_module = mock.MagicMock()
    pytorch_module.vision_module.forward.return_value = torch.tensor([0.302,0.92,0.38])
    pytorch_module.tokenizer = mock.MagicMock()
    pytorch_module.tokenizer.tokenize.return_value = torch.rand(1, 1)
    pytorch_module.processor = mock.MagicMock()
    pytorch_module.processor.process.return_value = torch.rand(1, 1)
    result = pytorch_module.generate_embedding('text', 'search query text')
    assert result == torch.tensor([0.302,0.92,0.38]).tolist()

@mock.patch('torch.jit.load', return_value=torch.jit.ScriptModule())
@mock.patch('transformers.AutoTokenizer.from_pretrained', return_value=None)
@mock.patch('transformers.AutoImageProcessor.from_pretrained', return_value=None)
@mock.patch('PIL.Image.open', return_value=np.ndarray((1, 1)))
def test_pytorch_photo_success(_, __, ___, ____):
    pytorch_module = PyTorchModule({'tokenizer_dir':'tokenizer_path','processor_dir':'processor_path','text_file':'text_model.pt','vision_file':'vision_model.pt'})
    pytorch_module.text_module = mock.MagicMock()
    pytorch_module.text_module.forward.return_value = torch.tensor([0.302,0.92,0.38])
    pytorch_module.vision_module = mock.MagicMock()
    pytorch_module.vision_module.forward.return_value = torch.tensor([0.302,0.92,0.38])
    pytorch_module.tokenizer = mock.MagicMock()
    pytorch_module.tokenizer.tokenize.return_value = torch.rand(1, 1)
    pytorch_module.processor = mock.MagicMock()
    pytorch_module.processor.process.return_value = torch.rand(1, 1)
    result = pytorch_module.generate_embedding('file', {'type':'photo','previewPath':'filePath'})
    assert result == [torch.tensor([0.302,0.92,0.38]).tolist()]

@mock.patch('torch.jit.load', return_value=torch.jit.ScriptModule())
@mock.patch('transformers.AutoTokenizer.from_pretrained', return_value=None)
@mock.patch('transformers.AutoImageProcessor.from_pretrained', return_value=None)
def test_pytorch_failed_empty_response(_, __, ___):
    pytorch_module = PyTorchModule({'tokenizer_dir':'tokenizer_path','processor_dir':'processor_path','text_file':'text_model.pt','vision_file':'vision_model.pt'})
    pytorch_module.text_module = mock.MagicMock()
    pytorch_module.text_module.forward.return_value = None
    pytorch_module.vision_module = mock.MagicMock()
    pytorch_module.vision_module.forward.return_value = torch.tensor([0.302,0.92,0.38])
    pytorch_module.tokenizer = mock.MagicMock()
    pytorch_module.tokenizer.tokenize.return_value = torch.rand(1, 1)
    pytorch_module.processor = mock.MagicMock()
    pytorch_module.processor.process.return_value = torch.rand(1, 1)
    result = pytorch_module.generate_embedding('text', 'search query text')
    assert result == None

@mock.patch('torch.jit.load', return_value=torch.jit.ScriptModule())
@mock.patch('transformers.AutoTokenizer.from_pretrained', return_value=None)
@mock.patch('transformers.AutoImageProcessor.from_pretrained', return_value=None)
def test_pytorch_failed_exception(_, __, ___):
    pytorch_module = PyTorchModule({'tokenizer_dir':'tokenizer_path','processor_dir':'processor_path','text_file':'text_model.pt','vision_file':'vision_model.pt'})
    pytorch_module.text_module = mock.MagicMock()
    pytorch_module.text_module.forward.side_effect = Exception('some error')
    pytorch_module.vision_module = mock.MagicMock()
    pytorch_module.vision_module.forward.return_value = torch.tensor([0.302,0.92,0.38])
    pytorch_module.tokenizer = mock.MagicMock()
    pytorch_module.tokenizer.tokenize.return_value = torch.rand(1, 1)
    pytorch_module.processor = mock.MagicMock()
    pytorch_module.processor.process.return_value = torch.rand(1, 1)
    with pytest.raises(Exception):
        result = pytorch_module.generate_embedding('text', 'search query text')
        assert result == None
