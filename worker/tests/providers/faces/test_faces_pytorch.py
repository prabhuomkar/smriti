"""Providers Faces PyTorch Tests"""
from unittest import mock
import pytest

import torch
import numpy as np

from src.providers.faces.pytorch import PyTorchModule


@mock.patch('facenet_pytorch.MTCNN.__init__', return_value=None)
@mock.patch('PIL.Image.open', return_value=np.ndarray((1, 1)))
def test_pytorch_success(_, __):
    pytorch_module = PyTorchModule(['1', '0.8', None])
    pytorch_module.det_model = mock.MagicMock()
    pytorch_module.det_model.forward.return_value = (torch.tensor([[0.302,0.92,0.38]]), torch.tensor([0.99]))
    pytorch_module.rec_model = mock.MagicMock()
    pytorch_module.rec_model.forward.return_value = torch.tensor([[0.302,0.92,0.38]])
    result = pytorch_module.detect('mediaitem_user_id', 'mediaitem_id', 'previews/file_name')
    assert len(result['thumbnails']) == 1
    del result['thumbnails']
    assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'embeddings': torch.tensor([[0.302,0.92,0.38]]).tolist()}

@mock.patch('facenet_pytorch.MTCNN.__init__', return_value=None)
@mock.patch('PIL.Image.open', return_value=np.ndarray((1, 1)))
def test_pytorch_failed_empty_faces(_, __):
    pytorch_module = PyTorchModule(['1', '0.8', None])
    pytorch_module.det_model = mock.MagicMock()
    pytorch_module.det_model.forward.return_value = (None, None)
    pytorch_module.rec_model = mock.MagicMock()
    pytorch_module.rec_model.forward.return_value = torch.tensor([[0.302,0.92,0.38]])
    result = pytorch_module.detect('mediaitem_user_id', 'mediaitem_id', 'previews/file_name')
    assert result == None

@mock.patch('facenet_pytorch.MTCNN.__init__', return_value=None)
@mock.patch('PIL.Image.open', return_value=np.ndarray((1, 1)))
def test_pytorch_failed_empty_embeddings(_, __):
    pytorch_module = PyTorchModule(['1', '078', None])
    pytorch_module.det_model = mock.MagicMock()
    pytorch_module.det_model.forward.return_value = (torch.tensor([[0.302,0.92,0.38]]), torch.tensor([0.75]))
    pytorch_module.rec_model = mock.MagicMock()
    pytorch_module.rec_model.forward.return_value = torch.tensor([[0.302,0.92,0.38]])
    result = pytorch_module.detect('mediaitem_user_id', 'mediaitem_id', 'previews/file_name')
    assert result == None

@mock.patch('facenet_pytorch.MTCNN.__init__', return_value=None)
@mock.patch('PIL.Image.open', return_value=np.ndarray((1, 1)))
def test_pytorch_failed_exception(_, __):
    pytorch_module = PyTorchModule(['1', '0.8', None])
    pytorch_module.det_model = mock.MagicMock()
    pytorch_module.det_model.forward.side_effect = Exception('some error')
    with pytest.raises(Exception):
        result = pytorch_module.detect('mediaitem_user_id', 'mediaitem_id', 'previews/file_name')
        assert result == None
