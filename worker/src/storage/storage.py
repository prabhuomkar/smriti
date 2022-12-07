"""Storage"""
from abc import ABC, abstractmethod


class Storage(ABC):
    """Storage"""

    def __init__(self, name) -> None:
        self.name = name

    @abstractmethod
    def connect(self) -> None:
        """Initialize connection"""
        raise NotImplementedError

    @abstractmethod
    def reconnect(self) -> None:
        """Re-establish connection"""
        raise NotImplementedError

    @abstractmethod
    def upload(self, id: str, name: str, offset: int, content: bytes) -> str:
        """Upload file chunks"""
        raise NotImplementedError

    @abstractmethod
    def delete(self, id: str) -> str:
        """Delete file"""
        raise NotImplementedError
