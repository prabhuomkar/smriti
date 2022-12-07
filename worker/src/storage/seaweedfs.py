"""Storage: SeaweedFS"""
from .storage import Storage


class SeaweedFS(Storage):
    """SeaweedFS Storage"""

    def __init__(self, host: str, port: int) -> None:
        super().__init__(name='seaweedfs')

    def connect(self) -> None:
        """Initialize connection"""
        raise NotImplementedError

    def reconnect(self) -> None:
        """Re-establish connection"""
        raise NotImplementedError

    def upload(self, id: str, name: str, offset: int, content: bytes) -> str:
        """Upload file chunks"""
        raise NotImplementedError

    def delete(self, id: str) -> str:
        """Delete file"""
        raise NotImplementedError
