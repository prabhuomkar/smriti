"""Providers OCR Paddle Tests"""
from unittest import mock
import pytest

from src.providers.ocr.paddlepaddle import PaddleModule


class OCRResult(dict):
    def __getattr__(self, attr):
        return self.get(attr, None)

@mock.patch('fastdeploy.vision.ocr.ppocr.PPOCRv3.__init__', return_value=None)
@mock.patch('fastdeploy.vision.ocr.ppocr.DBDetector.__init__', return_value=None)
@mock.patch('fastdeploy.vision.ocr.ppocr.Classifier.__init__', return_value=None)
@mock.patch('fastdeploy.vision.ocr.ppocr.Recognizer.__init__', return_value=None)
def test_paddle_success(_, __, ___, ____):
    paddle_module = PaddleModule({'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'})
    paddle_module.model = mock.MagicMock()
    paddle_module.model.predict.return_value = OCRResult({'text': ['pizza', '4.99', 'offer']})
    result = paddle_module.extract('mediaitem_user_id', 'mediaitem_id', 'photo', 'previews/file_name')
    if 'value' in result.keys():
        sorted(result['value'])
        assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'name': 'ocr', 'value': ['pizza', '4.99', 'offer']}

@mock.patch('fastdeploy.vision.ocr.ppocr.PPOCRv3.__init__', return_value=None)
@mock.patch('fastdeploy.vision.ocr.ppocr.DBDetector.__init__', return_value=None)
@mock.patch('fastdeploy.vision.ocr.ppocr.Classifier.__init__', return_value=None)
@mock.patch('fastdeploy.vision.ocr.ppocr.Recognizer.__init__', return_value=None)
def test_paddle_failed_empty_response(_, __, ___, ____):
    paddle_module = PaddleModule({'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'})
    paddle_module.model = mock.MagicMock()
    paddle_module.model.predict.return_value = None
    result = paddle_module.extract('mediaitem_user_id', 'mediaitem_id', 'photo', 'previews/file_name')
    assert result == None

@mock.patch('fastdeploy.vision.ocr.ppocr.PPOCRv3.__init__', return_value=None)
@mock.patch('fastdeploy.vision.ocr.ppocr.DBDetector.__init__', return_value=None)
@mock.patch('fastdeploy.vision.ocr.ppocr.Classifier.__init__', return_value=None)
@mock.patch('fastdeploy.vision.ocr.ppocr.Recognizer.__init__', return_value=None)
def test_paddle_failed_exception(_, __, ___, ____):
    paddle_module = PaddleModule({'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'})
    paddle_module.model = mock.MagicMock()
    paddle_module.model.predict.side_effect = Exception('some error')
    with pytest.raises(Exception):
        result = paddle_module.extract('mediaitem_user_id', 'mediaitem_id', 'photo', 'previews/file_name')
        assert result == None
