"""Storage: Disk"""
import os


class Disk:
    """Disk Storage"""

    def __init__(self, root: str) -> None:
        self.root = root
        if not os.path.exists(f'{self.root}/originals'):
            os.makedirs(f'{self.root}/originals')
        if not os.path.exists(f'{self.root}/thumbnails'):
            os.makedirs(f'{self.root}/thumbnails')

    def upload(self, id: str, offset: int, content: bytes, type: str = 'originals') -> str:
        """Upload file chunks"""
        with open(f'{self.root}/{type}/{id}', 'ab') as file_bytes:
            file_bytes.write(content)
        return id

    def delete(self, id: str) -> None:
        """Delete file"""
        os.remove(f'{self.root}/originals/{id}')
        os.remove(f'{self.root}/thumbnails/{id}')
