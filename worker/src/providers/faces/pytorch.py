"""Faces: Tensorflow"""
import os
import logging

from PIL import Image
from facenet_pytorch import MTCNN, InceptionResnetV1


class PyTorchModule:
    """PyTorchModule Faces"""

    def __init__(self, params: list[str]) -> None:
        os.environ['TORCH_HOME'] = '/'
        try:
            os.symlink('/models/faces/', '/checkpoints')
        except Exception as exp:
            logging.error(f'error creating symlink: {str(exp)}')
        self.prob_threshold = float(params[1])
        self.det_model = MTCNN(keep_all=True)
        self.rec_model = InceptionResnetV1(pretrained=params[2], classify=False)
        if self.rec_model:
            self.rec_model.eval()

    def detect(self, mediaitem_user_id: str, mediaitem_id: str, input_file: str) -> dict:
        """Detect faces for mediaitem"""
        embeddings = []
        faces, probs = self.det_model.forward(Image.open(input_file), return_prob=True)
        logging.debug(f'detected {len(faces) if faces is not None else None} \
                      faces for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        if faces is None or len(faces) == 0:
            return None

        result = self.rec_model.forward(faces)
        logging.debug(f'features calculated for {len(faces)} faces \
                       for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        for res, prob in zip(result, probs):
            if prob > self.prob_threshold:
                embeddings.append(res.tolist())

        logging.debug(f'got face embeddings for user {mediaitem_user_id} mediaitem {mediaitem_id}: {embeddings}')

        if len(embeddings) == 0:
            return None

        return dict({
            'userId': mediaitem_user_id,
            'id': mediaitem_id,
            'embeddings': embeddings,
        })
