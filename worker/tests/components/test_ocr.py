"""Components OCR Tests"""
from unittest import mock
import pytest

import grpc

from src.components.ocr import OCR
from src.protos.api_pb2_grpc import APIStub


@mock.patch('src.providers.ocr.PaddleModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_ocr_process_success(_):
    ocr = OCR(APIStub(channel=grpc.insecure_channel('')), 'paddle', {'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'})
    ocr.model = mock.MagicMock()
    ocr.model.extract.return_value = dict({'userId':'userId','id':'id','words':['first','second']})
    result = await ocr.process('mediaitem_user_id', 'mediaitem_id', None,
                    {'previewPath': 'location/to-preview-file'})
    assert result == {'keywords':'first second','previewPath':'location/to-preview-file'}

@mock.patch('src.providers.ocr.PaddleModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_ocr_process_success_with_keywords(_):
    ocr = OCR(APIStub(channel=grpc.insecure_channel('')), 'paddle', {'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'})
    ocr.model = mock.MagicMock()
    ocr.model.extract.return_value = dict({'userId':'userId','id':'id','words':['first','second']})
    result = await ocr.process('mediaitem_user_id', 'mediaitem_id', None,
                    {'previewPath': 'location/to-preview-file', 'keywords': 'exists'})
    assert result == {'keywords':'exists first second','previewPath':'location/to-preview-file'}

@mock.patch('src.providers.ocr.PaddleModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_ocr_process_success_no_result(_):
    result = await OCR(None, 'paddle', {'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'}).process('mediaitem_user_id', 'mediaitem_id', None, None)
    assert result == None

@mock.patch('src.providers.ocr.PaddleModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_ocr_process_failed_process_exception(_):
    ocr = OCR(None, 'paddle', {'det_model_dir':'/det_models','rec_model_dir':'/rec_models','cls_model_dir':'/cls_models'})
    ocr.model = mock.MagicMock()
    ocr.model.extract.side_effect = Exception('some exception')
    result = await ocr.process('mediaitem_user_id', 'mediaitem_id', None,
                    {'previewPath': 'location/to-preview-file'})
    assert result == {'previewPath':'location/to-preview-file'}
