"""Component Utils"""
import logging

from grpc import RpcError

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import MediaItemMetadataRequest  # pylint: disable=no-name-in-module


def grpc_save_mediaitem_metadata(api_stub: APIStub, request: MediaItemMetadataRequest):
    """gRPC call for saving mediaitem metadata"""
    try:
        _ = api_stub.SaveMediaItemMetadata(request)
    except RpcError as e:
        logging.error(
            f'error sending result for mediaitem {request.id}: {str(e)}')

def getval_from_dict(data, keys: list[str], return_type: str = 'str') -> str | int | float | None:
    """Get possible value from dictionary"""
    for key in keys:
        if key in data:
            if return_type == 'int':
                return int(data[key])
            if return_type == 'float':
                return float(data[key])
            return str(data[key])
    return None
