"""Faces Utils"""
from src.providers.faces.pytorch import PyTorchModule


def init_faces(name: str, params: list[str]) -> None | PyTorchModule:
    """Initialize faces by name"""
    if name == 'pytorch':
        return PyTorchModule(params)
    return None
