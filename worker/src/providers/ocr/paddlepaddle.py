"""OCR: Paddle"""
import logging

from moviepy.editor import VideoFileClip
from paddleocr import PaddleOCR


logging.getLogger().handlers = []

class PaddleModule:
    """PaddleModule OCR"""

    def __init__(self, params: dict) -> None:
        self.model = PaddleOCR(show_log=False, use_angle_cls=True, lang='en',
                               det_model_dir=f'/models/ocr/{params["det_model_dir"]}',
                               rec_model_dir=f'/models/ocr/{params["rec_model_dir"]}',
                               cls_model_dir=f'/models/ocr/{params["cls_model_dir"]}')

    def extract(self, mediaitem_user_id: str, mediaitem_id: str, mediaitem_type: str, input_file: str) -> dict:
        """Extract text from mediaitem"""
        result, words = [], []
        if mediaitem_type == 'photo':
            result = self.model.ocr(input_file)
        else:
            video_clip = VideoFileClip(input_file)
            for frame in video_clip.iter_frames(fps=video_clip.fps):
                _result = self.model.ocr(frame)
                result += _result
            video_clip.reader.close()
        for res in result:
            for line in res:
                if line[1][1] > 0.90:
                    words += line[1][0].split()
        logging.debug(f'extracted text for user {mediaitem_user_id} mediaitem {mediaitem_id}: {words}')

        if len(words) == 0:
            return None

        return dict({
            'userId': mediaitem_user_id,
            'id': mediaitem_id,
            'words': list(set(words)),
        })
