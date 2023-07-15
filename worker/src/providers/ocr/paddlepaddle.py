"""OCR: Paddle"""
import logging

from paddleocr import PaddleOCR


logging.getLogger().handlers = []

class PaddleModule:
    """PaddleModule OCR"""

    def __init__(self, params: list[str]) -> None:
        self.model = PaddleOCR(show_log=False, use_angle_cls=True, lang='en', det_model_dir=f'/models/{params[0]}',
                               rec_model_dir=f'/models/{params[1]}', cls_model_dir=f'/models/{params[2]}')

    def extract(self, mediaitem_user_id: str, mediaitem_id: str, input_file: str) -> dict:
        """Extract text from mediaitem"""
        result = self.model.ocr(input_file)
        words = []
        for res in result:
            for line in res:
                if line[1][1] > 0.85:
                    words += line[1][0].split()
        logging.debug(f'extracted text for user {mediaitem_user_id} mediaitem {mediaitem_id}: {words}')

        if len(words) == 0:
            return None

        return dict({
            'userId': mediaitem_user_id,
            'id': mediaitem_id,
            'words': list(set(words)),
        })
