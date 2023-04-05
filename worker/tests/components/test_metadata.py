"""Metadata Component Tests"""
from unittest.mock import patch
import pytest

import exiftool
import grpc

from src.store.disk import Disk
from src.components.metadata import process_metadata
from src.protos.api_pb2_grpc import APIStub


@pytest.mark.asyncio
async def test_metadata_success():
    with patch.object(Disk, 'get', return_value='storage-url'):
        with patch.object(exiftool.ExifToolHelper, 'get_metadata', return_value=[{'File:MIMEType': 'image/png'}]):
            _ = await process_metadata(Disk('/tmp/'), APIStub(grpc.insecure_channel('.')), 'mediaitem-user-id', 'mediaitem-id')
    # work(omkar): update tests

@pytest.mark.asyncio
async def test_metadata_failed_extraction():
    with patch.object(Disk, 'get', return_value='storage-url'):
        with patch.object(exiftool.ExifToolHelper, 'get_metadata', side_effect=exiftool.exceptions.ExifToolExecuteError(None, None, None, None)):
            _ = await process_metadata(Disk('/tmp/'), APIStub(grpc.insecure_channel('.')), 'mediaitem-user-id', 'mediaitem-id')
    # work(omkar): update tests
