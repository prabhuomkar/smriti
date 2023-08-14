"""Providers OCR Paddle Tests"""
from unittest import mock
import pytest

from src.providers.ocr.paddlepaddle import PaddleModule


@mock.patch('paddleocr.PaddleOCR.__init__', return_value=None)
def test_paddle_success(_):
    paddle_module = PaddleModule({'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'})
    paddle_module.model = mock.MagicMock()
    paddle_module.model.ocr.return_value = [[[[[106.0, 284.0], [266.0, 284.0], [266.0, 321.0], [106.0, 321.0]], ('offer', 0.99)],
                                             [[[106.0, 284.0], [266.0, 284.0], [266.0, 321.0], [106.0, 321.0]], ('pizza', 0.96)],
                                             [[[106.0, 284.0], [266.0, 284.0], [266.0, 321.0], [106.0, 321.0]], ('4.99', 0.94)]]]
    result = paddle_module.extract('mediaitem_user_id', 'mediaitem_id', 'photo', 'previews/file_name')
    if 'value' in result.keys():
        sorted(result['value'])
        assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'name': 'ocr', 'value': ['pizza', '4.99', 'offer']}

@mock.patch('paddleocr.PaddleOCR.__init__', return_value=None)
def test_paddle_failed_empty_response(_):
    paddle_module = PaddleModule({'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'})
    paddle_module.model = mock.MagicMock()
    paddle_module.model.ocr.return_value = {}
    result = paddle_module.extract('mediaitem_user_id', 'mediaitem_id', 'photo', 'previews/file_name')
    assert result == None

@mock.patch('paddleocr.PaddleOCR.__init__', return_value=None)
def test_paddle_failed_exception(_):
    paddle_module = PaddleModule({'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'})
    paddle_module.model = mock.MagicMock()
    paddle_module.model.ocr.side_effect = Exception('some error')
    with pytest.raises(Exception):
        result = paddle_module.extract('mediaitem_user_id', 'mediaitem_id', 'photo', 'previews/file_name')
        assert result == None
