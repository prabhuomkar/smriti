"""OCR: Paddle"""
import logging

import cv2
from moviepy.editor import VideoFileClip
from rapidocr_onnxruntime import RapidOCR


class PaddleModule:
    """PaddleModule OCR"""

    def __init__(self, params: dict) -> None:
        self.model = RapidOCR(
            det_model_path=f'/models/ocr/{params["det_model_dir"]}/model.onnx',
            cls_model_path=f'/models/ocr/{params["cls_model_dir"]}/model.onnx',
            rec_model_path=f'/models/ocr/{params["rec_model_dir"]}/model.onnx',
        )

    def extract(self, mediaitem_user_id: str, mediaitem_id: str, mediaitem_type: str, input_file: str) -> dict:
        """Extract text from mediaitem"""
        words = []
        if mediaitem_type == 'photo':
            result = self.model(cv2.imread(input_file))
            if result is not None and len(result) == 2:
                predictions = result[0]
                if predictions is not None:
                    for prediction in predictions:
                        if len(prediction) == 3:
                            words += prediction[1].split()
        else:
            video_clip = VideoFileClip(input_file)
            for frame in video_clip.iter_frames(fps=video_clip.fps):
                _result = self.model(frame)
                if _result is not None and len(_result) == 2:
                    predictions = _result[0]
                    if predictions is not None:
                        for prediction in predictions:
                            if len(prediction) == 3:
                                words += prediction[1].split()
            video_clip.reader.close()
        logging.debug(f'extracted text for user {mediaitem_user_id} mediaitem {mediaitem_id}: {words}')
        if len(words) == 0:
            return None

        return dict({
            'userId': mediaitem_user_id,
            'id': mediaitem_id,
            'words': list(set(words)),
        })
