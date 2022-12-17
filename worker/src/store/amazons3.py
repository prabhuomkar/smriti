"""Storage: AmazonS3"""


class AmazonS3:
    """AmazonS3 Storage"""

    def __init__(self, access_key: str, secret_key: str) -> None:
        self.access_key = access_key
        self.secret_key = secret_key

    def upload(self, mediaitem_id: str, offset: int, content: bytes, mediaitem_type: str = 'originals') -> str:
        """Upload file chunks"""
        raise NotImplementedError

    def get(self, mediaitem_id: str, mediaitem_type: str = 'originals') -> str:
        """Get file"""
        raise NotImplementedError

    def delete(self, mediaitem_id: str) -> None:
        """Delete file"""
        raise NotImplementedError
