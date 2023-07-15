"""Search Utils"""
from src.providers.search.pytorch import PyTorchModule


def init_search(name: str, files: list[str]) -> None | PyTorchModule:
    """Initialize search model by name"""
    if name == 'pytorch':
        return PyTorchModule(files)
    return None
