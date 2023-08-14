"""Components Faces Tests"""
from unittest import mock
import pytest

import grpc

from src.components.faces import Faces
from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import GetUsersResponse, MediaItemFaceEmbeddingsResponse, MediaItemFaceEmbedding, MediaItemEmbedding


@mock.patch('src.providers.faces.PyTorchModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_faces_process_success(_):
    faces = Faces(APIStub(channel=grpc.insecure_channel('')), 'pytorch', {'minutes':'1','face_threshold':'0.9','model':'vggface2'})
    faces.source = mock.MagicMock()
    faces.source.detect.return_value = dict({'userId':'userId','id':'id','embeddings':[[0.4,0.2]]})
    result = await faces.process('mediaitem_user_id', 'mediaitem_id', None,
                    {'previewPath': 'location/to-preview-file'})
    assert result == {'previewPath':'location/to-preview-file'}

@mock.patch('src.providers.faces.PyTorchModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_faces_process_success_no_result(_):
    result = await Faces(None, 'pytorch', {'minutes':'1','face_threshold':'0.9','model':'vggface2'}).process('mediaitem_user_id', 'mediaitem_id', None, None)
    assert result == None

@mock.patch('src.providers.faces.PyTorchModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_faces_process_failed_process_exception(_):
    faces = Faces(None, 'pytorch', {'minutes':'1','face_threshold':'0.9','model':'vggface2'})
    faces.source = mock.MagicMock()
    faces.source.detect.side_effect = Exception('some exception')
    result = await faces.process('mediaitem_user_id', 'mediaitem_id', None,
                    {'previewPath': 'location/to-preview-file'})
    assert result == {'previewPath':'location/to-preview-file'}


@mock.patch('src.providers.faces.PyTorchModule.__init__', return_value=None)
@mock.patch('src.components.Faces._grpc_save_mediaitem_people', return_value=None)
@pytest.mark.asyncio
async def test_faces_cluster_success(_, __):
    api_stub_mock = mock.MagicMock()
    api_stub_mock.GetUsers.return_value = GetUsersResponse(users=['user-id'])
    api_stub_mock.GetMediaItemFaceEmbeddings.return_value = MediaItemFaceEmbeddingsResponse(
        mediaItemFaceEmbeddings=[
            MediaItemFaceEmbedding(id='face-id-1', mediaItemId='mediaitem-id-1', peopleId='', embedding=MediaItemEmbedding(embedding=[0.42,0.64])),
            MediaItemFaceEmbedding(id='face-id-2', mediaItemId='mediaitem-id-2', peopleId='people-id-1', embedding=MediaItemEmbedding(embedding=[0.41,0.64])),
            MediaItemFaceEmbedding(id='face-id-3', mediaItemId='mediaitem-id-3', peopleId='', embedding=MediaItemEmbedding(embedding=[3.02,4.64])),
            MediaItemFaceEmbedding(id='face-id-4', mediaItemId='mediaitem-id-4', peopleId='', embedding=MediaItemEmbedding(embedding=[3.42,4.24])),
            MediaItemFaceEmbedding(id='face-id-5', mediaItemId='mediaitem-id-5', peopleId='', embedding=MediaItemEmbedding(embedding=[33.42,41.24])),
        ]
    )
    faces = Faces(api_stub_mock, 'pytorch', {'minutes':'1','face_threshold':'0.9','model':'vggface2'})
    faces.cluster()

@mock.patch('src.providers.faces.PyTorchModule.__init__', return_value=None)
@mock.patch('src.components.Faces._grpc_save_mediaitem_people', return_value=None)
@pytest.mark.asyncio
async def test_faces_cluster_success_no_embeddings(_, __):
    api_stub_mock = mock.MagicMock()
    api_stub_mock.GetUsers.return_value = GetUsersResponse(users=['user-id'])
    api_stub_mock.GetMediaItemFaceEmbeddings.return_value = MediaItemFaceEmbeddingsResponse(mediaItemFaceEmbeddings=[])
    faces = Faces(api_stub_mock, 'pytorch', {'minutes':'1','face_threshold':'0.9','model':'vggface2'})
    faces.cluster()

@mock.patch('src.providers.faces.PyTorchModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_faces_cluster_failed_cluster_exception(_):
    api_stub_mock = mock.MagicMock()
    api_stub_mock.GetUsers.return_value = GetUsersResponse(users=['user-id'])
    api_stub_mock.GetMediaItemFaceEmbeddings.side_effect = grpc.RpcError(Exception('some error'))
    faces = Faces(api_stub_mock, 'pytorch', {'minutes':'1','face_threshold':'0.9','model':'vggface2'})
    faces.cluster()

@mock.patch('src.providers.faces.PyTorchModule.__init__', return_value=None)
@pytest.mark.asyncio
async def test_faces_cluster_failed_grpc_exception(_):
    api_stub_mock = mock.MagicMock()
    api_stub_mock.GetUsers.return_value = GetUsersResponse(users=['user-id'])
    api_stub_mock.GetMediaItemFaceEmbeddings.return_value = MediaItemFaceEmbeddingsResponse(
        mediaItemFaceEmbeddings=[
            MediaItemFaceEmbedding(id='face-id-1', mediaItemId='mediaitem-id-1', peopleId='', embedding=MediaItemEmbedding(embedding=[0.42,0.64])),
            MediaItemFaceEmbedding(id='face-id-2', mediaItemId='mediaitem-id-2', peopleId='people-id-1', embedding=MediaItemEmbedding(embedding=[0.41,0.64])),
        ]
    )
    api_stub_mock.SaveMediaItemPeople.side_effect = grpc.RpcError(Exception('some error'))
    faces = Faces(api_stub_mock, 'pytorch', {'minutes':'1','face_threshold':'0.9','model':'vggface2'})
    faces.cluster()
