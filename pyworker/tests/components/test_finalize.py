"""Components Finalize Tests"""
from unittest import mock
import pytest
from requests.exceptions import HTTPError

import grpc

from src.components.finalize import Finalize
from src.protos.api_pb2_grpc import APIStub


@mock.patch('src.components.Finalize._grpc_final_save_mediaitem', return_value=None)
@pytest.mark.asyncio
async def test_finalize_process_success(_):
    result = await Finalize(None).process('mediaitem_user_id', 'mediaitem_id', None,
                    {'keywords': '', 'embeddings': []})
    assert result == None

@pytest.mark.asyncio
async def test_finalize_process_grpc_exception():
    grpc_mock = mock.MagicMock()
    grpc_mock.side_effect = grpc.RpcError(Exception('some error'))
    with mock.patch('src.protos.API.SaveMediaItemFinalResult', grpc_mock):
        result = await Finalize(APIStub(channel=grpc.insecure_channel(''))).process('mediaitem_user_id', 'mediaitem_id', 
                        None, {'keywords': '', 'embeddings': []})
        assert result == None