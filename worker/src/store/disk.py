"""Storage: Disk"""
import os


class Disk:
    """Disk Storage"""

    def __init__(self, root: str) -> None:
        self.root = root
        if not os.path.exists(f'{self.root}/originals'):
            os.mkdir(f'{self.root}/originals')
        if not os.path.exists(f'{self.root}/previews'):
            os.mkdir(f'{self.root}/previews')
        if not os.path.exists(f'{self.root}/thumbnails'):
            os.mkdir(f'{self.root}/thumbnails')

    def upload(self, mediaitem_id: str, _: int, content: bytes, mediaitem_type: str = 'originals') -> str:
        """Upload file chunks"""
        with open(f'{self.root}/{mediaitem_type}/{mediaitem_id}', 'ab') as file_bytes:
            file_bytes.write(content)
        return f'{self.root}/{mediaitem_type}/{mediaitem_id}'

    def get(self, mediaitem_id: str, mediaitem_type: str = 'originals') -> str:
        """Get file"""
        return f'{self.root}/{mediaitem_type}/{mediaitem_id}'

    def delete(self, mediaitem_id: str) -> None:
        """Delete file"""
        os.remove(f'{self.root}/originals/{mediaitem_id}')
        os.remove(f'{self.root}/previews/{mediaitem_id}')
        os.remove(f'{self.root}/thumbnails/{mediaitem_id}')
