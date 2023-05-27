"""Storage: MinIO"""
from minio import Minio


class MinIO:
    """MinIO Storage"""

    def __init__(self, endpoint: str, access_key: str, secret_key: str) -> None:
        self.client = Minio(endpoint, access_key, secret_key)
        for bucket_name in ["originals", "previews", "thumbnails"]:
            found = self.client.bucket_exists(bucket_name)
            if not found:
                self.client.make_bucket(bucket_name)

    def upload(self, mediaitem_id: str, content: bytes, mediaitem_type: str = 'originals') -> str:
        """Upload file chunks"""
        self.client.put_object(bucket_name=mediaitem_type, object_name=mediaitem_id,
                               data=content, length=-1, part_size=len(content))
        return f'{mediaitem_type}/{mediaitem_id}'

    def get(self, mediaitem_id: str, mediaitem_type: str = 'originals') -> str:
        """Get file"""
        file_path = f'{mediaitem_type}-{mediaitem_id}'
        self.client.fget_object(mediaitem_type, mediaitem_id, file_path)
        return file_path

    def delete(self, mediaitem_id: str) -> None:
        """Delete file"""
        self.client.remove_object('originals', mediaitem_id)
        self.client.remove_object('previews', mediaitem_id)
        self.client.remove_object('thumbnails', mediaitem_id)
