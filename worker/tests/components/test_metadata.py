"""Components Metadata Tests"""
from unittest import mock
import pytest

import grpc

from src.components.metadata import Metadata
from src.protos.api_pb2_grpc import APIStub


@mock.patch('src.components.Metadata._grpc_save_mediaitem_metadata', return_value=None)
@pytest.mark.asyncio
async def test_metadata_process_photo_success(_):
    exiftool_mock = mock.MagicMock()
    exiftool_mock.return_value.__getitem__.return_value = {
        'File:MIMEType': 'image/jpeg', 'Composite:ImageSize': '14124 3100',
        'EXIF:Make': 'Apple', 'EXIF:Model': 'iPhone 12 mini', 'EXIF:FocalLength': '5.49', 'EXIF:ExposureTime': '1/20',
        'EXIF:FNumber': '1.8', 'EXIF:ISO': '758', 'EXIF:GPSLatitudeRef': 'N', 'EXIF:GPSLatitude': '19.292902',
        'EXIF:GPSLongitudeRef': 'E', 'EXIF:GPSLongitude': '70.2822',
    }
    with mock.patch('exiftool.ExifToolHelper.get_metadata', exiftool_mock):
        result = await Metadata(None).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', None)
        del result['sourcePath']
        del result['exifdata']
        assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'status': 'READY', 
                        'type': 'photo', 'category': 'panorama', 'mimeType': 'image/jpeg', 'width': 14124, 
                        'height': 3100, 'cameraMake': 'Apple', 'cameraModel': 'iPhone 12 mini', 
                        'focalLength': '5.49', 'apertureFNumber': '1.8', 'isoEquivalent': '758', 'exposureTime': '1/20', 
                        'fps': None, 'latitude': 19.292902, 'longitude': 70.2822}

@mock.patch('src.components.Metadata._grpc_save_mediaitem_metadata', return_value=None)
@pytest.mark.asyncio
async def test_metadata_process_video_success(_):
    exiftool_mock = mock.MagicMock()
    exiftool_mock.return_value.__getitem__.return_value = {
        'File:MIMEType': 'video/mpeg', 'QuickTime:ImageWidth': '1080', 'QuickTime:ImageHeight': '720',
        'EXIF:Make': 'Apple', 'EXIF:Model': 'iPhone 12 mini', 'EXIF:FocalLength': '5.49', 'QuickTime:LivePhotoAuto': '1',
        'EXIF:FNumber': '1.8', 'EXIF:ISO': '758', 'EXIF:GPSLatitudeRef': 'N', 'EXIF:GPSLatitude': '19.292902',
        'EXIF:GPSLongitudeRef': 'E', 'EXIF:GPSLongitude': '70.2822', 'QuickTime:VideoFrameRate': '60', 'EXIF:ExposureTime': '1/20',
    }
    with mock.patch('exiftool.ExifToolHelper.get_metadata', exiftool_mock):
        result = await Metadata(None).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', None)
        del result['sourcePath']
        del result['exifdata']
        assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'status': 'READY', 
                        'type': 'video', 'category': 'live', 'mimeType': 'video/mpeg', 'width': 1080, 
                        'height': 720, 'cameraMake': 'Apple', 'cameraModel': 'iPhone 12 mini', 
                        'focalLength': '5.49', 'apertureFNumber': '1.8', 'isoEquivalent': '758', 'exposureTime': '1/20', 
                        'fps': '60', 'latitude': 19.292902, 'longitude': 70.2822}

@mock.patch('src.components.Metadata._grpc_save_mediaitem_metadata', return_value=None)
@pytest.mark.asyncio
async def test_metadata_process_video_category_slow_success(_):
    exiftool_mock = mock.MagicMock()
    exiftool_mock.return_value.__getitem__.return_value = {
        'File:MIMEType': 'video/mpeg', 'QuickTime:ImageWidth': '1080', 'QuickTime:ImageHeight': '720',
        'EXIF:Make': 'Apple', 'EXIF:Model': 'iPhone 12 mini', 'EXIF:FocalLength': '5.49', 'QuickTime:LivePhotoAuto': '1',
        'EXIF:FNumber': '1.8', 'EXIF:ISO': '758', 'EXIF:GPSLatitudeRef': 'N', 'EXIF:GPSLatitude': '19.292902',
        'EXIF:GPSLongitudeRef': 'E', 'EXIF:GPSLongitude': '70.2822', 'QuickTime:VideoFrameRate': '200', 'EXIF:ExposureTime': '1/20',
    }
    with mock.patch('exiftool.ExifToolHelper.get_metadata', exiftool_mock):
        result = await Metadata(None).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', None)
        del result['sourcePath']
        del result['exifdata']
        assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'status': 'READY', 
                        'type': 'video', 'category': 'slow', 'mimeType': 'video/mpeg', 'width': 1080, 
                        'height': 720, 'cameraMake': 'Apple', 'cameraModel': 'iPhone 12 mini', 
                        'focalLength': '5.49', 'apertureFNumber': '1.8', 'isoEquivalent': '758', 'exposureTime': '1/20', 
                        'fps': '200', 'latitude': 19.292902, 'longitude': 70.2822}

@mock.patch('src.components.Metadata._grpc_save_mediaitem_metadata', return_value=None)
@pytest.mark.asyncio
async def test_metadata_process_video_category_timelapse_success(_):
    exiftool_mock = mock.MagicMock()
    exiftool_mock.return_value.__getitem__.return_value = {
        'File:MIMEType': 'video/mpeg', 'QuickTime:ImageWidth': '1080', 'QuickTime:ImageHeight': '720',
        'EXIF:Make': 'Apple', 'EXIF:Model': 'iPhone 12 mini', 'EXIF:FocalLength': '5.49', 'QuickTime:CaptureMode': 'time-lapse',
        'EXIF:FNumber': '1.8', 'EXIF:ISO': '758', 'EXIF:GPSLatitudeRef': 'N', 'EXIF:GPSLatitude': '19.292902',
        'EXIF:GPSLongitudeRef': 'E', 'EXIF:GPSLongitude': '70.2822', 'QuickTime:VideoFrameRate': '60', 'EXIF:ExposureTime': '1/20',
    }
    with mock.patch('exiftool.ExifToolHelper.get_metadata', exiftool_mock):
        result = await Metadata(None).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', None)
        del result['sourcePath']
        del result['exifdata']
        assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'status': 'READY', 
                        'type': 'video', 'category': 'timelapse', 'mimeType': 'video/mpeg', 'width': 1080, 
                        'height': 720, 'cameraMake': 'Apple', 'cameraModel': 'iPhone 12 mini', 
                        'focalLength': '5.49', 'apertureFNumber': '1.8', 'isoEquivalent': '758', 'exposureTime': '1/20', 
                        'fps': '60', 'latitude': 19.292902, 'longitude': 70.2822}

@mock.patch('src.components.Metadata._grpc_save_mediaitem_metadata', return_value=None)
@pytest.mark.asyncio
async def test_metadata_process_video_category_default_success(_):
    exiftool_mock = mock.MagicMock()
    exiftool_mock.return_value.__getitem__.return_value = {
        'File:MIMEType': 'video/mpeg', 'QuickTime:ImageWidth': '1080', 'QuickTime:ImageHeight': '720',
        'EXIF:Make': 'Apple', 'EXIF:Model': 'iPhone 12 mini', 'EXIF:FocalLength': '5.49',
        'EXIF:FNumber': '1.8', 'EXIF:ISO': '758', 'EXIF:GPSLatitudeRef': 'N', 'EXIF:GPSLatitude': '19.292902',
        'EXIF:GPSLongitudeRef': 'E', 'EXIF:GPSLongitude': '70.2822', 'QuickTime:VideoFrameRate': '60', 'EXIF:ExposureTime': '1/20',
    }
    with mock.patch('exiftool.ExifToolHelper.get_metadata', exiftool_mock):
        result = await Metadata(None).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', None)
        del result['sourcePath']
        del result['exifdata']
        assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'status': 'READY', 
                        'type': 'video', 'category': 'default', 'mimeType': 'video/mpeg', 'width': 1080, 
                        'height': 720, 'cameraMake': 'Apple', 'cameraModel': 'iPhone 12 mini', 
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

@pytest.mark.asyncio
async def test_metadata_process_grpc_exception():
    exiftool_mock = mock.MagicMock()
    exiftool_mock.return_value.__getitem__.return_value = {
        'File:MIMEType': 'image/jpeg', 'File:ImageWidth': '4022', 'File:ImageHeight': '3100', 'EXIF:UserComment': 'Screenshot',
        'EXIF:Make': 'Apple', 'EXIF:Model': 'iPhone 12 mini', 'EXIF:FocalLength': '5.49', 'EXIF:ExposureTime': '1/20',
        'EXIF:FNumber': '1.8', 'EXIF:ISO': '758', 'Composite:GPSPosition': '19.292902 70.2822',
    }
    grpc_mock = mock.MagicMock()
    grpc_mock.side_effect = grpc.RpcError(Exception('some error'))
    with mock.patch('exiftool.ExifToolHelper.get_metadata', exiftool_mock):
        with mock.patch('src.protos.API.SaveMediaItemMetadata', grpc_mock):
            result = await Metadata(APIStub(channel=grpc.insecure_channel(''))).process('mediaitem_user_id', 'mediaitem_id', 'mediaitem_file_path', None)
            del result['sourcePath']
            del result['exifdata']
            assert result == {'userId': 'mediaitem_user_id', 'id': 'mediaitem_id', 'status': 'READY', 
                        'type': 'photo', 'category': 'screenshot', 'mimeType': 'image/jpeg', 'width': 4022, 
                        'height': 3100, 'cameraMake': 'Apple', 'cameraModel': 'iPhone 12 mini', 
                        'focalLength': '5.49', 'apertureFNumber': '1.8', 'isoEquivalent': '758', 'exposureTime': '1/20', 
                        'fps': None, 'latitude': 19.292902, 'longitude': 70.2822}
