"""Providers OCR Utils Tests"""
from unittest import mock

from src.providers.ocr.paddlepaddle import PaddleModule
from src.providers.ocr.utils import init_ocr


@mock.patch('src.providers.ocr.PaddleModule.__init__', return_value=None)
def test_init_ocr_success(_):
    assert init_ocr('tesseract', []) == None
    assert isinstance(init_ocr('paddlepaddle', ['/det_models', '/rec_models', '/cls_models']), PaddleModule)
