"""Components Classification Tests"""
from unittest import mock
import pytest

import grpc

from src.components.classification import Classification
from src.protos.api_pb2_grpc import APIStub


@mock.patch('torch.jit.load', return_value=None)
@mock.patch('src.components.Classification._grpc_save_mediaitem_thing', return_value=None)
@pytest.mark.asyncio
async def test_classification_process_success(_, __):
    classification = Classification(APIStub(channel=grpc.insecure_channel('')), 'pytorch', ['model_name.pt'])
    classification.model = mock.MagicMock()
    classification.model.classify.return_value = dict({'userId':'userId','id':'id','name':'name'})
    result = await classification.process('mediaitem_user_id', 'mediaitem_id', None,
                    {'previewPath': 'location/to-preview-file'})
    assert result == None

@mock.patch('torch.jit.load', return_value=None)
@pytest.mark.asyncio
async def test_classification_process_success_no_metadata(_):
    result = await Classification(None, 'pytorch', ['model_name.pt']).process('mediaitem_user_id', 'mediaitem_id', None, None)
    assert result == None

@mock.patch('torch.jit.load', return_value=None)
@mock.patch('src.components.Classification._grpc_save_mediaitem_thing', return_value=None)
@pytest.mark.asyncio
async def test_classification_process_failed_process_exception(_, __):
    classification = Classification(None, 'pytorch', ['model_name.pt'])
    classification.model = mock.MagicMock()
    classification.model.classify.side_effect = Exception('some exception')
    result = await classification.process('mediaitem_user_id', 'mediaitem_id', None,
                    {'previewPath': 'location/to-preview-file'})
    assert result == None

@mock.patch('torch.jit.load', return_value=None)
@pytest.mark.asyncio
async def test_classification_process_grpc_exception(_):
    classification = Classification(APIStub(channel=grpc.insecure_channel('')), 'pytorch', ['model_name.pt'])
    classification.model = mock.MagicMock()
    classification.model.classify.return_value = dict({'userId':'userId','id':'id','name':'name'})
    grpc_mock = mock.MagicMock()
    grpc_mock.side_effect = grpc.RpcError(Exception('some error'))
    with mock.patch('src.protos.API.SaveMediaItemThing', grpc_mock):
        result = await classification.process('mediaitem_user_id', 'mediaitem_id', 
                        None, {'previewPath': 'location/to-preview-file'})
        assert result == None