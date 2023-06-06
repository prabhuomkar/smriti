"""Components Metadata Tests"""
from unittest import mock
import pytest

import grpc

from src.components.metadata import Metadata
from src.protos.api_pb2_grpc import APIStub


@mock.patch('src.components.Metadata._grpc_save_mediaitem_metadata', return_value=None)
@mock.patch('src.components.Metadata._generate_photo_preview_and_thumbnail', return_value=(bytes(), bytes()))
@pytest.mark.asyncio
async def test_metadata_process_photo_success(_, __):
    exiftool_mock = mock.MagicMock()
    exiftool_mock.return_value.__getitem__.return_value = {
        'File:MIMEType': 'image/jpeg', 'Composite:ImageSize': '14124 3100',
        'EXIF:Make': 'Apple', 'EXIF:Model': 'iPhone 12 mini', 'EXIF:FocalLength': '5.49', 'EXIF:ExposureTime': '1/20',
        'EXIF:FNumber': '1.8', 'EXIF:ISO': '758', 'EXIF:GPSLatitude': '19.292902', 'EXIF:GPSLongitude': '70.2822',
    }
    with mock.patch('exiftool.ExifToolHelper.get_metadata', exiftool_mock):
        result = await Metadata(None).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', None)
        del result['sourceUrl']
        del result['previewUrl']
        del result['thumbnailUrl']
        assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'status': 'READY', 
                        'type': 'photo', 'category': 'panorama', 'mimeType': 'image/jpeg', 'width': 14124, 
                        'height': 3100, 'creationTime': None, 'cameraMake': 'Apple', 'cameraModel': 'iPhone 12 mini', 
                        'focalLength': '5.49', 'apertureFNumber': '1.8', 'isoEquivalent': '758', 'exposureTime': '1/20', 
                        'fps': None, 'latitude': 19.292902, 'longitude': 70.2822}

@mock.patch('src.components.Metadata._grpc_save_mediaitem_metadata', return_value=None)
@mock.patch('src.components.Metadata._generate_video_preview_and_thumbnail', return_value=(bytes(), bytes()))
@pytest.mark.asyncio
async def test_metadata_process_video_success(_, __):
    exiftool_mock = mock.MagicMock()
    exiftool_mock.return_value.__getitem__.return_value = {
        'File:MIMEType': 'video/mpeg', 'QuickTime:ImageWidth': '1080', 'QuickTime:ImageHeight': '720',
        'EXIF:Make': 'Apple', 'EXIF:Model': 'iPhone 12 mini', 'EXIF:FocalLength': '5.49', 'QuickTime:LivePhotoAuto': '1',
        'EXIF:FNumber': '1.8', 'EXIF:ISO': '758', 'EXIF:GPSLatitude': '19.292902', 'EXIF:GPSLongitude': '70.2822',
        'QuickTime:VideoFrameRate': '60', 'EXIF:ExposureTime': '1/20',
    }
    with mock.patch('exiftool.ExifToolHelper.get_metadata', exiftool_mock):
        result = await Metadata(None).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', None)
        del result['sourceUrl']
        del result['previewUrl']
        del result['thumbnailUrl']
        assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'status': 'READY', 
                        'type': 'video', 'category': 'live', 'mimeType': 'video/mpeg', 'width': 1080, 
                        'height': 720, 'creationTime': None, 'cameraMake': 'Apple', 'cameraModel': 'iPhone 12 mini', 
                        'focalLength': '5.49', 'apertureFNumber': '1.8', 'isoEquivalent': '758', 'exposureTime': '1/20', 
                        'fps': '60', 'latitude': 19.292902, 'longitude': 70.2822}

@mock.patch('src.components.Metadata._grpc_save_mediaitem_metadata', return_value=None)
@pytest.mark.asyncio
async def test_metadata_process_failed_process_exception(_):
    exiftool_mock = mock.MagicMock()
    exiftool_mock.return_value.__getitem__.return_value = None
    with mock.patch('exiftool.ExifToolHelper.get_metadata', exiftool_mock):
        result = await Metadata(None).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', None)
        assert result == None

@mock.patch('src.components.Metadata._generate_photo_preview_and_thumbnail', return_value=(bytes(), bytes()))
@pytest.mark.asyncio
async def test_metadata_process_grpc_exception(_):
    exiftool_mock = mock.MagicMock()
    exiftool_mock.return_value.__getitem__.return_value = {
        'File:MIMEType': 'image/jpeg', 'File:ImageWidth': '4022', 'File:ImageHeight': '3100', 'EXIF:UserComment': 'Screenshot',
        'EXIF:Make': 'Apple', 'EXIF:Model': 'iPhone 12 mini', 'EXIF:FocalLength': '5.49', 'EXIF:ExposureTime': '1/20',
        'EXIF:FNumber': '1.8', 'EXIF:ISO': '758', 'EXIF:GPSLatitude': '19.292902', 'EXIF:GPSLongitude': '70.2822',
    }
    grpc_mock = mock.MagicMock()
    grpc_mock.side_effect = grpc.RpcError(Exception('some error'))
    with mock.patch('exiftool.ExifToolHelper.get_metadata', exiftool_mock):
        with mock.patch('src.protos.API.SaveMediaItemMetadata', grpc_mock):
            result = await Metadata(APIStub(channel=grpc.insecure_channel(''))).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', None)
            del result['sourceUrl']
            del result['previewUrl']
            del result['thumbnailUrl']
            assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'status': 'READY', 
                        'type': 'photo', 'category': 'screenshot', 'mimeType': 'image/jpeg', 'width': 4022, 
                        'height': 3100, 'creationTime': None, 'cameraMake': 'Apple', 'cameraModel': 'iPhone 12 mini', 
                        'focalLength': '5.49', 'apertureFNumber': '1.8', 'isoEquivalent': '758', 'exposureTime': '1/20', 
                        'fps': None, 'latitude': 19.292902, 'longitude': 70.2822}

@mock.patch('src.components.Metadata._grpc_save_mediaitem_metadata', return_value=None)
@pytest.mark.asyncio
async def test_metadata_process_photo_preview_thumbnail_exception(_):
    exiftool_mock = mock.MagicMock()
    exiftool_mock.return_value.__getitem__.return_value = {
        'File:MIMEType': 'image/jpeg', 'File:ImageWidth': '14124', 'File:ImageHeight': '3100',
        'EXIF:Make': 'Apple', 'EXIF:Model': 'iPhone 12 mini', 'EXIF:FocalLength': '5.49', 'EXIF:ExposureTime': '1/20',
        'EXIF:FNumber': '1.8', 'EXIF:ISO': '758', 'EXIF:GPSLatitude': '19.292902', 'EXIF:GPSLongitude': '70.2822',
    }
    preview_thumbnail_mock = mock.MagicMock()
    preview_thumbnail_mock.side_effect = grpc.RpcError(Exception('some error'))
    with mock.patch('exiftool.ExifToolHelper.get_metadata', exiftool_mock):
        with mock.patch('src.components.Metadata._generate_photo_preview_and_thumbnail', preview_thumbnail_mock):
            result = await Metadata(None).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', None)
            assert result == None

@mock.patch('src.components.Metadata._grpc_save_mediaitem_metadata', return_value=None)
@pytest.mark.asyncio
async def test_metadata_process_video_preview_thumbnail_exception(_):
    exiftool_mock = mock.MagicMock()
    exiftool_mock.return_value.__getitem__.return_value = {
        'File:MIMEType': 'video/mpeg', 'QuickTime:ImageWidth': '1080', 'QuickTime:ImageHeight': '720',
        'EXIF:Make': 'Apple', 'EXIF:Model': 'iPhone 12 mini', 'EXIF:FocalLength': '5.49', 'QuickTime:LivePhotoAuto': '1',
        'EXIF:FNumber': '1.8', 'EXIF:ISO': '758', 'EXIF:GPSLatitude': '19.292902', 'EXIF:GPSLongitude': '70.2822',
        'QuickTime:VideoFrameRate': '60', 'EXIF:ExposureTime': '1/20',
    }
    preview_thumbnail_mock = mock.MagicMock()
    preview_thumbnail_mock.side_effect = grpc.RpcError(Exception('some error'))
    with mock.patch('exiftool.ExifToolHelper.get_metadata', exiftool_mock):
        with mock.patch('src.components.Metadata._generate_video_preview_and_thumbnail', preview_thumbnail_mock):
            result = await Metadata(None).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', None)
            assert result == None
