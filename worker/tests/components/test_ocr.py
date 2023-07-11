"""Components OCR Tests"""
from unittest import mock
import pytest

import grpc

from src.components.ocr import OCR
from src.protos.api_pb2_grpc import APIStub


@mock.patch('src.providers.ocr.PaddleModule.__init__', return_value=None)
@mock.patch('src.components.OCR._grpc_save_mediaitem_ml_result', return_value=None)
@pytest.mark.asyncio
async def test_ocr_process_success(_, __):
    ocr = OCR(APIStub(channel=grpc.insecure_channel('')), 'paddle', ['/det_models', '/rec_models', '/cls_models'])
    ocr.model = mock.MagicMock()
    ocr.model.extract.return_value = dict({'userId':'userId','id':'id','name':'name'})
    result = await ocr.process('mediaitem_user_id', 'mediaitem_id', None,
                    {'previewPath': 'location/to-preview-file'})
    assert result == None

@mock.patch('src.providers.ocr.PaddleModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_ocr_process_success_no_result(_):
    result = await OCR(None, 'paddle', ['/det_models', '/rec_models', '/cls_models']).process('mediaitem_user_id', 'mediaitem_id', None, None)
    assert result == None

@mock.patch('src.providers.ocr.PaddleModule.__init__', return_value=None)
@mock.patch('src.components.OCR._grpc_save_mediaitem_ml_result', return_value=None)
@pytest.mark.asyncio
async def test_ocr_process_failed_process_exception(_, __):
    ocr = OCR(None, 'paddle', ['/det_models', '/rec_models', '/cls_models'])
    ocr.model = mock.MagicMock()
    ocr.model.extract.side_effect = Exception('some exception')
    result = await ocr.process('mediaitem_user_id', 'mediaitem_id', None,
                    {'previewPath': 'location/to-preview-file'})
    assert result == None

@mock.patch('src.providers.ocr.PaddleModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_ocr_process_grpc_exception(_):
    ocr = OCR(APIStub(channel=grpc.insecure_channel('')), 'paddle', ['/det_models', '/rec_models', '/cls_models'])
    ocr.model = mock.MagicMock()
    ocr.model.extract.return_value = dict({'userId':'userId','id':'id','words':['value']})
    grpc_mock = mock.MagicMock()
    grpc_mock.side_effect = grpc.RpcError(Exception('some error'))
    with mock.patch('src.protos.API.SaveMediaItemMLResult', grpc_mock):
        result = await ocr.process('mediaitem_user_id', 'mediaitem_id', 
                        None, {'previewPath': 'location/to-preview-file'})
        assert result == None
