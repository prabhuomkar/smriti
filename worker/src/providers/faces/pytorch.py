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
        self.prob_threshold = float(params[0])
        self.det_model = MTCNN(keep_all=True)
        self.rec_model = InceptionResnetV1(pretrained=params[1], classify=True).eval()

    def detect(self, mediaitem_user_id: str, mediaitem_id: str, input_file: str) -> dict:
        """Detect faces for mediaitem"""
        embeddings = []
        faces, probs = self.det_model(Image.open(input_file), return_prob=True)
        logging.debug(f'detected {len(faces) if faces else None} \
                      faces for user {mediaitem_user_id} mediaitem {mediaitem_id}')

        if not faces or len(faces) == 0:
            return None

        result = self.rec_model(faces)
        for res, prob in zip(result, probs):
            if prob > self.prob_threshold:
                embeddings.append(res)

        logging.debug(f'got face embeddings for user {mediaitem_user_id} mediaitem {mediaitem_id}: {embeddings}')

        if len(embeddings) == 0:
            return None

        return dict({
            'userId': mediaitem_user_id,
            'id': mediaitem_id,
            'embeddings': embeddings,
        })
