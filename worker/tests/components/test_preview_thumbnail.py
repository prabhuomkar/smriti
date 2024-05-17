"""Components Metadata Tests"""
from unittest import mock
import pytest

import grpc

from src.components.preview_thumbnail import PreviewThumbnail
from src.protos.api_pb2_grpc import APIStub


@mock.patch('src.components.PreviewThumbnail._grpc_save_mediaitem_preview_thumbnail', return_value=None)
@mock.patch('src.components.PreviewThumbnail._generate_photo_preview_and_thumbnail_and_placeholder', return_value=(bytes(), bytes(), bytes()))
@pytest.mark.asyncio
async def test_preview_thumbnail_process_photo_success(_, __):
    metadata = {
        'type': 'photo', 'sourcePath': 'sourcePath', 'mimeType': 'image/png',
    }
    result = await PreviewThumbnail(None, params={'thumbnail_size':'512'}).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', metadata)
    del result['sourcePath']
    del result['previewPath']
    del result['thumbnailPath']
    del result['placeholder']
    assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'status': 'READY', 'mimeType': 'image/png', 'type': 'photo'}

@mock.patch('src.components.PreviewThumbnail._grpc_save_mediaitem_preview_thumbnail', return_value=None)
@mock.patch('src.components.PreviewThumbnail._generate_video_preview_and_thumbnail_and_placeholder', return_value=(bytes(), bytes(), bytes()))
@pytest.mark.asyncio
async def test_preview_thumbnail_process_video_success(_, __):
    metadata = {
        'type': 'video', 'sourcePath': 'sourcePath',
    }
    result = await PreviewThumbnail(None, params={'thumbnail_size':'512'}).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', metadata)
    del result['sourcePath']
    del result['previewPath']
    del result['thumbnailPath']
    del result['placeholder']
    assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'status': 'READY', 'type': 'video'}

@mock.patch('src.components.PreviewThumbnail._generate_photo_preview_and_thumbnail_and_placeholder', return_value=(bytes(), bytes(), bytes()))
@pytest.mark.asyncio
async def test_preview_thumbnail_process_grpc_exception(_):
    metadata = {
        'type': 'photo', 'sourcePath': 'sourcePath', 'mimeType': 'image/png',
    }
    grpc_mock = mock.MagicMock()
    grpc_mock.side_effect = grpc.RpcError(Exception('some error'))
    with mock.patch('src.protos.API.SaveMediaItemPreviewThumbnail', grpc_mock):
        result = await PreviewThumbnail(APIStub(channel=grpc.insecure_channel('')), params={'thumbnail_size':'512'}).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', metadata)
        del result['sourcePath']
        del result['previewPath']
        del result['thumbnailPath']
        del result['placeholder']
        assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'status': 'READY', 'mimeType': 'image/png', 'type': 'photo'}

@mock.patch('src.components.PreviewThumbnail._grpc_save_mediaitem_preview_thumbnail', return_value=None)
@pytest.mark.asyncio
async def test_preview_thumbnail_process_photo_preview_thumbnail_exception(_):
    metadata = {
        'type': 'photo', 'sourcePath': 'sourcePath', 'mimeType': 'image/png',
    }
    preview_thumbnail_mock = mock.MagicMock()
    preview_thumbnail_mock.side_effect = grpc.RpcError(Exception('some error'))
    with mock.patch('src.components.PreviewThumbnail._generate_photo_preview_and_thumbnail_and_placeholder', preview_thumbnail_mock):
        result = await PreviewThumbnail(None, params={'thumbnail_size':'512'}).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', metadata)
        assert result == None

@mock.patch('src.components.PreviewThumbnail._grpc_save_mediaitem_preview_thumbnail', return_value=None)
@pytest.mark.asyncio
async def test_preview_thumbnail_process_video_preview_thumbnail_exception(_):
    metadata = {
        'type': 'video', 'sourcePath': 'sourcePath',
    }
    preview_thumbnail_mock = mock.MagicMock()
    preview_thumbnail_mock.side_effect = grpc.RpcError(Exception('some error'))
    with mock.patch('src.components.PreviewThumbnail._generate_video_preview_and_thumbnail_and_placeholder', preview_thumbnail_mock):
        result = await PreviewThumbnail(None, params={'thumbnail_size':'512'}).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', metadata)
        assert result == None
