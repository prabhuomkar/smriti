"""Component Utils"""
import logging

from grpc import RpcError

from src.protos.api_pb2_grpc import APIStub
from src.protos.api_pb2 import MediaItemMetadataRequest, MediaItemPlaceRequest  # pylint: disable=no-name-in-module


def grpc_save_mediaitem_metadata(api_stub: APIStub, result: dict):
    """gRPC call for saving mediaitem metadata"""
    try:
        request = MediaItemMetadataRequest(
            userId=result['userId'],
            id=result['id'],
            status=result['status'],
            mimeType=result['mimeType'] if 'mimeType' in result else None,
            sourceUrl=result['sourceUrl'] if 'sourceUrl' in result else None,
            type=result['type'] if 'type' in result else None,
            width=result['width'] if 'width' in result else None,
            height=result['height'] if 'height' in result else None,
            creationTime=result['creationTime'] if 'creationTime' in result else None,
            cameraMake=result['cameraMake'] if 'cameraMake' in result else None,
            cameraModel=result['cameraModel'] if 'cameraModel' in result else None,
            focalLength=result['focalLength'] if 'focalLength' in result else None,
            apertureFNumber=result['apertureFNumber'] if 'apertureFNumber' in result else None,
            isoEquivalent=result['isoEquivalent'] if 'isoEquivalent' in result else None,
            exposureTime=result['exposureTime'] if 'exposureTime' in result else None,
            fps=result['fps'] if 'fps' in result else None,
            latitude=result['latitude'] if 'latitude' in result else None,
            longitude=result['longitude'] if 'longitude' in result else None,
        )
        _ = api_stub.SaveMediaItemMetadata(request)
    except RpcError as e:
        logging.error(
            f'error sending metadata for mediaitem {request.id}: {str(e)}')

def grpc_save_mediaitem_place(api_stub: APIStub, result: dict):
    """gRPC call for saving mediaitem place"""
    try:
        request = MediaItemPlaceRequest(
            userId=result['userId'],
            id=result['id'],
            postcode=result['postcode'] if 'postcode' in result else None,
            country=result['country'] if 'country' in result else None,
            state=result['state'] if 'state' in result else None,
            city=result['city'] if 'city' in result else None,
            town=result['town'] if 'town' in result else None,
        )
        _ = api_stub.SaveMediaItemPlace(request)
    except RpcError as e:
        logging.error(
            f'error sending place for mediaitem {request.id}: {str(e)}')

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
