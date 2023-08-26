"""OCR: Paddle"""
import logging

import cv2
from moviepy.editor import VideoFileClip
import fastdeploy as fd


class PaddleModule:
    """PaddleModule OCR"""

    def __init__(self, params: dict) -> None:
        self.default_option = fd.RuntimeOption()
        self.default_option.use_ort_backend()
        self.det_model = fd.vision.ocr.DBDetector(
            model_file=f'/models/ocr/{params["det_model_dir"]}/model.onnx',
            runtime_option=self.default_option,
            model_format=fd.ModelFormat.ONNX,
        )
        self.cls_model = fd.vision.ocr.Classifier(
            model_file=f'/models/ocr/{params["cls_model_dir"]}/model.onnx',
            runtime_option=self.default_option,
            model_format=fd.ModelFormat.ONNX,
        )
        self.rec_model = fd.vision.ocr.Recognizer(
            model_file=f'/models/ocr/{params["rec_model_dir"]}/model.onnx',
            label_path=f'/models/ocr/{params["rec_model_dir"]}/en_dict.txt',
            runtime_option=self.default_option,
            model_format=fd.ModelFormat.ONNX,
        )
        self.model = fd.vision.ocr.PPOCRv3(
            det_model=self.det_model,
            cls_model=self.cls_model,
            rec_model=self.rec_model,
        )

    def extract(self, mediaitem_user_id: str, mediaitem_id: str, mediaitem_type: str, input_file: str) -> dict:
        """Extract text from mediaitem"""
        result, words = [], []
        if mediaitem_type == 'photo':
            result = self.model.predict(cv2.imread(input_file))
            if result is not None:
                words = result.text
        else:
            video_clip = VideoFileClip(input_file)
            for frame in video_clip.iter_frames(fps=video_clip.fps):
                _result = self.model.predict(frame)
                if _result is not None:
                    result += _result.text
            video_clip.reader.close()
            words = result
        logging.debug(f'extracted text for user {mediaitem_user_id} mediaitem {mediaitem_id}: {words}')
        if len(words) == 0:
            return None

        return dict({
            'userId': mediaitem_user_id,
            'id': mediaitem_id,
            'words': list(set(words)),
        })
