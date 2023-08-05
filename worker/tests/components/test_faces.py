"""Components Faces Tests"""
from unittest import mock
import pytest

import grpc

from src.components.faces import Faces
from src.protos.api_pb2_grpc import APIStub


@mock.patch('src.providers.faces.PyTorchModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_faces_process_success(_):
    faces = Faces(APIStub(channel=grpc.insecure_channel('')), 'pytorch', ['0.9', 'vggface2'])
    faces.source = mock.MagicMock()
    faces.source.detect.return_value = dict({'userId':'userId','id':'id','embeddings':[[0.4,0.2]]})
    result = await faces.process('mediaitem_user_id', 'mediaitem_id', None,
                    {'previewPath': 'location/to-preview-file'})
    assert result == {'previewPath':'location/to-preview-file'}

@mock.patch('src.providers.faces.PyTorchModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_faces_process_success_no_result(_):
    result = await Faces(None, 'pytorch', ['0.9', 'vggface2']).process('mediaitem_user_id', 'mediaitem_id', None, None)
    assert result == None

@mock.patch('src.providers.faces.PyTorchModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_faces_process_failed_process_exception(_):
    faces = Faces(None, 'pytorch', ['0.9', 'vggface2'])
    faces.source = mock.MagicMock()
    faces.source.detect.side_effect = Exception('some exception')
    result = await faces.process('mediaitem_user_id', 'mediaitem_id', None,
                    {'previewPath': 'location/to-preview-file'})
    assert result == {'previewPath':'location/to-preview-file'}
