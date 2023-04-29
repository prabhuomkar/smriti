"""Components Places Tests"""
from unittest import mock
import pytest
from requests.exceptions import HTTPError

import grpc

from src.components.places import Places
from src.protos.api_pb2_grpc import APIStub


API_URL = 'https://nominatim.openstreetmap.org/reverse.php?zoom=18&format=jsonv2&lat={lat}&lon={lon}'

@pytest.mark.asyncio
async def test_places_process_success():
    result = await Places(None, 'openstreetmap').process('mediaitem_user_id', 'mediaitem_id', 
                    {'latitude': 19.2195856, 'longitude': 73.1056888})
    assert result == None

@pytest.mark.asyncio
async def test_places_process_success_no_metadata():
    result = await Places(None, 'openstreetmap').process('mediaitem_user_id', 'mediaitem_id', None)
    assert result == None

@mock.patch('src.components.Places._grpc_save_mediaitem_place', return_value=None)
@pytest.mark.asyncio
async def test_places_process_failed_process_exception(_):
    with mock.patch('requests.get', side_effect=HTTPError()):
        result = await Places(None, 'openstreetmap').process('mediaitem_user_id', 'mediaitem_id', 
                        {'latitude': 19.2195856, 'longitude': 73.1056888})
        assert result == None

@pytest.mark.asyncio
async def test_places_process_grpc_exception():
    grpc_mock = mock.MagicMock()
    grpc_mock.side_effect = grpc.RpcError(Exception('some error'))
    with mock.patch('src.protos.API.SaveMediaItemPlace', grpc_mock):
        result = await Places(APIStub(channel=grpc.insecure_channel('')), 'openstreetmap').process('mediaitem_user_id', 'mediaitem_id', 
                        {'latitude': 19.2195856, 'longitude': 73.1056888})
        assert result == None