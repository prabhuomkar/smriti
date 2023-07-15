"""OCR Utils"""
from src.providers.ocr.paddlepaddle import PaddleModule


def init_ocr(name: str, files: list[str]) -> None:
    """Initialize ocr model by name"""
    if name == 'paddlepaddle':
        return PaddleModule(files)
    return None
