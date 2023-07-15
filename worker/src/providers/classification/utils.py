"""Classification Utils"""
from src.providers.classification.pytorch import PyTorchModule


def init_classification(name: str, params: list[str]) -> None | PyTorchModule:
    """Initialize classification model by name"""
    if name == 'pytorch':
        return PyTorchModule(params)
    return None
