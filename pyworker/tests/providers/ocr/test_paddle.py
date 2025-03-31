"""Providers OCR Paddle Tests"""
from unittest import mock
import pytest

from src.providers.ocr.paddlepaddle import PaddleModule


class OCRResult(dict):
    def __getattr__(self, attr):
        return self.get(attr, None)

@mock.patch('rapidocr_onnxruntime.RapidOCR.__init__', return_value=None)
def test_paddle_success(_):
    paddle_module = PaddleModule({'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'})
    paddle_module.model = mock.MagicMock()
    paddle_module.model.return_value = ([
        [[[23.0, 7.0], [746.0, 11.0], [745.0, 118.0], [22.0, 114.0]], 'pizza', 0.42], 
        [[[724.0, 18.0], [888.0, 26.0], [885.0, 94.0], [721.0, 87.0]], '4.99', 0.69], 
        [[[24.0, 228.0], [468.0, 226.0], [468.0, 326.0], [25.0, 328.0]], 'offer', 0.37]], [0.1, 0.2, 0.3])
    result = paddle_module.extract('mediaitem_user_id', 'mediaitem_id', 'photo', 'previews/file_name')
    if 'value' in result.keys():
        sorted(result['value'])
        assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'name': 'ocr', 'value': ['pizza', '4.99', 'offer']}

@mock.patch('rapidocr_onnxruntime.RapidOCR.__init__', return_value=None)
def test_paddle_failed_empty_response(_):
    paddle_module = PaddleModule({'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'})
    paddle_module.model = mock.MagicMock()
    paddle_module.model.return_value = None
    result = paddle_module.extract('mediaitem_user_id', 'mediaitem_id', 'photo', 'previews/file_name')
    assert result == None

@mock.patch('rapidocr_onnxruntime.RapidOCR.__init__', return_value=None)
def test_paddle_failed_exception(_):
    paddle_module = PaddleModule({'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'})
    paddle_module.model = mock.MagicMock()
    paddle_module.model.side_effect = Exception('some error')
    with pytest.raises(Exception):
        result = paddle_module.extract('mediaitem_user_id', 'mediaitem_id', 'photo', 'previews/file_name')
        assert result == None
