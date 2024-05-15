from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MediaItemComponent(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    METADATA: _ClassVar[MediaItemComponent]
    PREVIEW_THUMBNAIL: _ClassVar[MediaItemComponent]
    PLACES: _ClassVar[MediaItemComponent]
    CLASSIFICATION: _ClassVar[MediaItemComponent]
    FACES: _ClassVar[MediaItemComponent]
    OCR: _ClassVar[MediaItemComponent]
    SEARCH: _ClassVar[MediaItemComponent]
METADATA: MediaItemComponent
PREVIEW_THUMBNAIL: MediaItemComponent
PLACES: MediaItemComponent
CLASSIFICATION: MediaItemComponent
FACES: MediaItemComponent
OCR: MediaItemComponent
SEARCH: MediaItemComponent

class MediaItemProcessRequest(_message.Message):
    __slots__ = ("userId", "id", "filePath", "components", "payload")
    class PayloadEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    USERID_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    FILEPATH_FIELD_NUMBER: _ClassVar[int]
    COMPONENTS_FIELD_NUMBER: _ClassVar[int]
    PAYLOAD_FIELD_NUMBER: _ClassVar[int]
    userId: str
    id: str
    filePath: str
    components: _containers.RepeatedScalarFieldContainer[MediaItemComponent]
    payload: _containers.ScalarMap[str, str]
    def __init__(self, userId: _Optional[str] = ..., id: _Optional[str] = ..., filePath: _Optional[str] = ..., components: _Optional[_Iterable[_Union[MediaItemComponent, str]]] = ..., payload: _Optional[_Mapping[str, str]] = ...) -> None: ...

class MediaItemProcessResponse(_message.Message):
    __slots__ = ("ok",)
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: bool = ...) -> None: ...

class GenerateEmbeddingRequest(_message.Message):
    __slots__ = ("text",)
    TEXT_FIELD_NUMBER: _ClassVar[int]
    text: str
    def __init__(self, text: _Optional[str] = ...) -> None: ...

class GenerateEmbeddingResponse(_message.Message):
    __slots__ = ("embedding",)
    EMBEDDING_FIELD_NUMBER: _ClassVar[int]
    embedding: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, embedding: _Optional[_Iterable[float]] = ...) -> None: ...
