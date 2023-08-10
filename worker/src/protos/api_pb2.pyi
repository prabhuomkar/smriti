from google.protobuf import empty_pb2 as _empty_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ConfigResponse(_message.Message):
    __slots__ = ["config"]
    CONFIG_FIELD_NUMBER: _ClassVar[int]
    config: bytes
    def __init__(self, config: _Optional[bytes] = ...) -> None: ...

class GetUsersResponse(_message.Message):
    __slots__ = ["users"]
    USERS_FIELD_NUMBER: _ClassVar[int]
    users: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, users: _Optional[_Iterable[str]] = ...) -> None: ...

class MediaItemEmbedding(_message.Message):
    __slots__ = ["embedding"]
    EMBEDDING_FIELD_NUMBER: _ClassVar[int]
    embedding: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, embedding: _Optional[_Iterable[float]] = ...) -> None: ...

class MediaItemFaceEmbedding(_message.Message):
    __slots__ = ["embedding", "id", "mediaItemId", "peopleId"]
    EMBEDDING_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    MEDIAITEMID_FIELD_NUMBER: _ClassVar[int]
    PEOPLEID_FIELD_NUMBER: _ClassVar[int]
    embedding: MediaItemEmbedding
    id: str
    mediaItemId: str
    peopleId: str
    def __init__(self, id: _Optional[str] = ..., mediaItemId: _Optional[str] = ..., peopleId: _Optional[str] = ..., embedding: _Optional[_Union[MediaItemEmbedding, _Mapping]] = ...) -> None: ...

class MediaItemFaceEmbeddingsRequest(_message.Message):
    __slots__ = ["userId"]
    USERID_FIELD_NUMBER: _ClassVar[int]
    userId: str
    def __init__(self, userId: _Optional[str] = ...) -> None: ...

class MediaItemFaceEmbeddingsResponse(_message.Message):
    __slots__ = ["mediaItemFaceEmbeddings"]
    MEDIAITEMFACEEMBEDDINGS_FIELD_NUMBER: _ClassVar[int]
    mediaItemFaceEmbeddings: _containers.RepeatedCompositeFieldContainer[MediaItemFaceEmbedding]
    def __init__(self, mediaItemFaceEmbeddings: _Optional[_Iterable[_Union[MediaItemFaceEmbedding, _Mapping]]] = ...) -> None: ...

class MediaItemFacePeople(_message.Message):
    __slots__ = ["facePeople"]
    class FacePeopleEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    FACEPEOPLE_FIELD_NUMBER: _ClassVar[int]
    facePeople: _containers.ScalarMap[str, str]
    def __init__(self, facePeople: _Optional[_Mapping[str, str]] = ...) -> None: ...

class MediaItemFacesRequest(_message.Message):
    __slots__ = ["embeddings", "id", "thumbnails", "userId"]
    EMBEDDINGS_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    USERID_FIELD_NUMBER: _ClassVar[int]
    embeddings: _containers.RepeatedCompositeFieldContainer[MediaItemEmbedding]
    id: str
    thumbnails: _containers.RepeatedScalarFieldContainer[str]
    userId: str
    def __init__(self, userId: _Optional[str] = ..., id: _Optional[str] = ..., embeddings: _Optional[_Iterable[_Union[MediaItemEmbedding, _Mapping]]] = ..., thumbnails: _Optional[_Iterable[str]] = ...) -> None: ...

class MediaItemFinalResultRequest(_message.Message):
    __slots__ = ["embeddings", "id", "keywords", "userId"]
    EMBEDDINGS_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    KEYWORDS_FIELD_NUMBER: _ClassVar[int]
    USERID_FIELD_NUMBER: _ClassVar[int]
    embeddings: _containers.RepeatedCompositeFieldContainer[MediaItemEmbedding]
    id: str
    keywords: str
    userId: str
    def __init__(self, userId: _Optional[str] = ..., id: _Optional[str] = ..., keywords: _Optional[str] = ..., embeddings: _Optional[_Iterable[_Union[MediaItemEmbedding, _Mapping]]] = ...) -> None: ...

class MediaItemMetadataRequest(_message.Message):
    __slots__ = ["apertureFNumber", "cameraMake", "cameraModel", "category", "creationTime", "exposureTime", "focalLength", "fps", "height", "id", "isoEquivalent", "latitude", "longitude", "mimeType", "placeholder", "previewPath", "sourcePath", "status", "thumbnailPath", "type", "userId", "width"]
    APERTUREFNUMBER_FIELD_NUMBER: _ClassVar[int]
    CAMERAMAKE_FIELD_NUMBER: _ClassVar[int]
    CAMERAMODEL_FIELD_NUMBER: _ClassVar[int]
    CATEGORY_FIELD_NUMBER: _ClassVar[int]
    CREATIONTIME_FIELD_NUMBER: _ClassVar[int]
    EXPOSURETIME_FIELD_NUMBER: _ClassVar[int]
    FOCALLENGTH_FIELD_NUMBER: _ClassVar[int]
    FPS_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    ISOEQUIVALENT_FIELD_NUMBER: _ClassVar[int]
    LATITUDE_FIELD_NUMBER: _ClassVar[int]
    LONGITUDE_FIELD_NUMBER: _ClassVar[int]
    MIMETYPE_FIELD_NUMBER: _ClassVar[int]
    PLACEHOLDER_FIELD_NUMBER: _ClassVar[int]
    PREVIEWPATH_FIELD_NUMBER: _ClassVar[int]
    SOURCEPATH_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILPATH_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    USERID_FIELD_NUMBER: _ClassVar[int]
    WIDTH_FIELD_NUMBER: _ClassVar[int]
    apertureFNumber: str
    cameraMake: str
    cameraModel: str
    category: str
    creationTime: str
    exposureTime: str
    focalLength: str
    fps: str
    height: int
    id: str
    isoEquivalent: str
    latitude: float
    longitude: float
    mimeType: str
    placeholder: str
    previewPath: str
    sourcePath: str
    status: str
    thumbnailPath: str
    type: str
    userId: str
    width: int
    def __init__(self, userId: _Optional[str] = ..., id: _Optional[str] = ..., status: _Optional[str] = ..., mimeType: _Optional[str] = ..., sourcePath: _Optional[str] = ..., placeholder: _Optional[str] = ..., previewPath: _Optional[str] = ..., thumbnailPath: _Optional[str] = ..., type: _Optional[str] = ..., category: _Optional[str] = ..., width: _Optional[int] = ..., height: _Optional[int] = ..., creationTime: _Optional[str] = ..., cameraMake: _Optional[str] = ..., cameraModel: _Optional[str] = ..., focalLength: _Optional[str] = ..., apertureFNumber: _Optional[str] = ..., isoEquivalent: _Optional[str] = ..., exposureTime: _Optional[str] = ..., fps: _Optional[str] = ..., latitude: _Optional[float] = ..., longitude: _Optional[float] = ...) -> None: ...

class MediaItemPeopleRequest(_message.Message):
    __slots__ = ["mediaItemFacePeople", "userId"]
    class MediaItemFacePeopleEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MediaItemFacePeople
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MediaItemFacePeople, _Mapping]] = ...) -> None: ...
    MEDIAITEMFACEPEOPLE_FIELD_NUMBER: _ClassVar[int]
    USERID_FIELD_NUMBER: _ClassVar[int]
    mediaItemFacePeople: _containers.MessageMap[str, MediaItemFacePeople]
    userId: str
    def __init__(self, userId: _Optional[str] = ..., mediaItemFacePeople: _Optional[_Mapping[str, MediaItemFacePeople]] = ...) -> None: ...

class MediaItemPlaceRequest(_message.Message):
    __slots__ = ["city", "country", "id", "postcode", "state", "town", "userId"]
    CITY_FIELD_NUMBER: _ClassVar[int]
    COUNTRY_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    POSTCODE_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    TOWN_FIELD_NUMBER: _ClassVar[int]
    USERID_FIELD_NUMBER: _ClassVar[int]
    city: str
    country: str
    id: str
    postcode: str
    state: str
    town: str
    userId: str
    def __init__(self, userId: _Optional[str] = ..., id: _Optional[str] = ..., postcode: _Optional[str] = ..., country: _Optional[str] = ..., state: _Optional[str] = ..., city: _Optional[str] = ..., town: _Optional[str] = ...) -> None: ...

class MediaItemThingRequest(_message.Message):
    __slots__ = ["id", "name", "userId"]
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    USERID_FIELD_NUMBER: _ClassVar[int]
    id: str
    name: str
    userId: str
    def __init__(self, userId: _Optional[str] = ..., id: _Optional[str] = ..., name: _Optional[str] = ...) -> None: ...
