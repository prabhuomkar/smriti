"""Storage: AmazonS3"""


class AmazonS3:
    """AmazonS3 Storage"""

    def __init__(self, access_key: str, secret_key: str) -> None:
        self.access_key = access_key
        self.secret_key = secret_key

    def upload(self, id: str, offset: int, content: bytes) -> str:
        """Upload file chunks"""
        raise NotImplementedError

    def get(self, id: str, type: str = 'originals') -> str:
        """Get file"""
        raise NotImplementedError

    def delete(self, id: str) -> None:
        """Delete file"""
        raise NotImplementedError
