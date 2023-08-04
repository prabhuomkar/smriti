from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class MediaItemProcessRequest(_message.Message):
    __slots__ = ["userId", "id", "filePath"]
    USERID_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    FILEPATH_FIELD_NUMBER: _ClassVar[int]
    userId: str
    id: str
    filePath: str
    def __init__(self, userId: _Optional[str] = ..., id: _Optional[str] = ..., filePath: _Optional[str] = ...) -> None: ...

class MediaItemProcessResponse(_message.Message):
    __slots__ = ["ok"]
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: bool = ...) -> None: ...

class GenerateEmbeddingRequest(_message.Message):
    __slots__ = ["text"]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    text: str
    def __init__(self, text: _Optional[str] = ...) -> None: ...

class GenerateEmbeddingResponse(_message.Message):
    __slots__ = ["embedding"]
    EMBEDDING_FIELD_NUMBER: _ClassVar[int]
    embedding: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, embedding: _Optional[_Iterable[float]] = ...) -> None: ...
