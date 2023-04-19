from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class MediaItemProcessRequest(_message.Message):
    __slots__ = ["command", "content", "features", "id", "offset", "userId"]
    COMMAND_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    FEATURES_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    OFFSET_FIELD_NUMBER: _ClassVar[int]
    USERID_FIELD_NUMBER: _ClassVar[int]
    command: str
    content: bytes
    features: _containers.RepeatedScalarFieldContainer[str]
    id: str
    offset: int
    userId: str
    def __init__(self, userId: _Optional[str] = ..., id: _Optional[str] = ..., offset: _Optional[int] = ..., command: _Optional[str] = ..., content: _Optional[bytes] = ..., features: _Optional[_Iterable[str]] = ...) -> None: ...

class MediaItemProcessResponse(_message.Message):
    __slots__ = ["ok"]
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: bool = ...) -> None: ...
