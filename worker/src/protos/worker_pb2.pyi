from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class MediaItemProcessRequest(_message.Message):
    __slots__ = ["command", "content", "id", "offset"]
    COMMAND_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    OFFSET_FIELD_NUMBER: _ClassVar[int]
    command: str
    content: bytes
    id: str
    offset: int
    def __init__(self, id: _Optional[str] = ..., offset: _Optional[int] = ..., command: _Optional[str] = ..., content: _Optional[bytes] = ...) -> None: ...

class MediaItemProcessResponse(_message.Message):
    __slots__ = ["ok"]
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: bool = ...) -> None: ...
