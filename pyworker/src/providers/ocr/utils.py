"""OCR Utils"""
from src.providers.ocr.paddlepaddle import PaddleModule


def init_ocr(name: str, params: dict) -> None | PaddleModule:
    """Initialize ocr model by name"""
    if name == 'paddlepaddle':
        return PaddleModule(params)
    return None
