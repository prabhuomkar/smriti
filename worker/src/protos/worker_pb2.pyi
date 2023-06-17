from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class MediaItemProcessRequest(_message.Message):
    __slots__ = ["filePath", "id", "userId"]
    FILEPATH_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    USERID_FIELD_NUMBER: _ClassVar[int]
    filePath: str
    id: str
    userId: str
    def __init__(self, userId: _Optional[str] = ..., id: _Optional[str] = ..., filePath: _Optional[str] = ...) -> None: ...

class MediaItemProcessResponse(_message.Message):
    __slots__ = ["ok"]
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: bool = ...) -> None: ...
