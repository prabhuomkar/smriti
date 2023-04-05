"""Place Component Tests"""
import pytest

import grpc

from src.components.place import API_URL, process_place
from src.protos.api_pb2_grpc import APIStub


@pytest.mark.asyncio
async def test_place_success(requests_mock):
    requests_mock.get(url=API_URL.format(lat=19.2195856, lon=73.1056888), json={
        'address': {'city': 'Dombivali', 'state': 'Maharashtra', 'postcode': '421201', 'country': 'India'}
    })
    _ = process_place(APIStub(grpc.insecure_channel('.')), 'mediaitem-user-id', 'mediaitem-id', [19.2195856, 73.1056888])
    # work(omkar): update tests

@pytest.mark.asyncio
async def test_place_failed_response(requests_mock):
    requests_mock.get(url=API_URL.format(
        lat=19.2195856, lon=73.1056888), json={}, status_code=500)
    _ = process_place(APIStub(grpc.insecure_channel('.')), 'mediaitem-user-id', 'mediaitem-id', [19.2195856, 73.1056888])
    # work(omkar): update tests
