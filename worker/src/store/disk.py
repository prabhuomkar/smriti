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

    def upload(self, id: str, offset: int, content: bytes, type: str = 'originals') -> str:
        """Upload file chunks"""
        with open(f'{self.root}/{type}/{id}', 'ab') as file_bytes:
            file_bytes.write(content)
        return id

    def get(self, id: str, type: str = 'originals') -> str:
        """Get file"""
        return f'{self.root}/{type}/{id}'

    def delete(self, id: str) -> None:
        """Delete file"""
        os.remove(f'{self.root}/originals/{id}')
        os.remove(f'{self.root}/previews/{id}')
        os.remove(f'{self.root}/thumbnails/{id}')
