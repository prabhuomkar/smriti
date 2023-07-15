"""OCR Utils"""
from src.providers.ocr.paddlepaddle import PaddleModule


def init_ocr(name: str, params: list[str]) -> None:
    """Initialize ocr model by name"""
    if name == 'paddlepaddle':
        return PaddleModule(params)
    return None
