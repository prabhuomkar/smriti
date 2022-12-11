from google.protobuf import empty_pb2 as _empty_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class MediaItemMetadataRequest(_message.Message):
    __slots__ = ["apertureFNumber", "cameraMake", "cameraModel", "creationTime", "exposureTime", "filename", "focalLength", "fps", "height", "id", "isoEquivalent", "latitude", "longitude", "mimeType", "previewUrl", "sourceUrl", "status", "thumbnailUrl", "type", "width"]
    APERTUREFNUMBER_FIELD_NUMBER: _ClassVar[int]
    CAMERAMAKE_FIELD_NUMBER: _ClassVar[int]
    CAMERAMODEL_FIELD_NUMBER: _ClassVar[int]
    CREATIONTIME_FIELD_NUMBER: _ClassVar[int]
    EXPOSURETIME_FIELD_NUMBER: _ClassVar[int]
    FILENAME_FIELD_NUMBER: _ClassVar[int]
    FOCALLENGTH_FIELD_NUMBER: _ClassVar[int]
    FPS_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    ISOEQUIVALENT_FIELD_NUMBER: _ClassVar[int]
    LATITUDE_FIELD_NUMBER: _ClassVar[int]
    LONGITUDE_FIELD_NUMBER: _ClassVar[int]
    MIMETYPE_FIELD_NUMBER: _ClassVar[int]
    PREVIEWURL_FIELD_NUMBER: _ClassVar[int]
    SOURCEURL_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILURL_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    WIDTH_FIELD_NUMBER: _ClassVar[int]
    apertureFNumber: str
    cameraMake: str
    cameraModel: str
    creationTime: str
    exposureTime: str
    filename: str
    focalLength: str
    fps: str
    height: int
    id: str
    isoEquivalent: str
    latitude: float
    longitude: float
    mimeType: str
    previewUrl: str
    sourceUrl: str
    status: str
    thumbnailUrl: str
    type: str
    width: int
    def __init__(self, id: _Optional[str] = ..., status: _Optional[str] = ..., filename: _Optional[str] = ..., mimeType: _Optional[str] = ..., sourceUrl: _Optional[str] = ..., previewUrl: _Optional[str] = ..., thumbnailUrl: _Optional[str] = ..., type: _Optional[str] = ..., width: _Optional[int] = ..., height: _Optional[int] = ..., creationTime: _Optional[str] = ..., cameraMake: _Optional[str] = ..., cameraModel: _Optional[str] = ..., focalLength: _Optional[str] = ..., apertureFNumber: _Optional[str] = ..., isoEquivalent: _Optional[str] = ..., exposureTime: _Optional[str] = ..., fps: _Optional[str] = ..., latitude: _Optional[float] = ..., longitude: _Optional[float] = ...) -> None: ...

class MediaItemPlaceRequest(_message.Message):
    __slots__ = ["city", "country", "id", "postcode", "state", "town"]
    CITY_FIELD_NUMBER: _ClassVar[int]
    COUNTRY_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    POSTCODE_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    TOWN_FIELD_NUMBER: _ClassVar[int]
    city: str
    country: str
    id: str
    postcode: str
    state: str
    town: str
    def __init__(self, id: _Optional[str] = ..., postcode: _Optional[str] = ..., country: _Optional[str] = ..., state: _Optional[str] = ..., city: _Optional[str] = ..., town: _Optional[str] = ...) -> None: ...
