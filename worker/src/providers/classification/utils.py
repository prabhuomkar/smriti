"""Classification Utils"""
from src.providers.classification.pytorch import PyTorchModule


def init_classification(name: str, files: list[str]) -> None | PyTorchModule:
    """Initialize classification model by name"""
    if name == 'pytorch':
        return PyTorchModule(files)
    return None
