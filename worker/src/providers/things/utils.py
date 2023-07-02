"""Things Utils"""
from src.providers.things.pytorch import PyTorchModule


def init_things(name: str, files: list[str]) -> None | PyTorchModule:
    """Initialize things model by name"""
    if name == 'pytorch':
        return PyTorchModule(files)
    return None
