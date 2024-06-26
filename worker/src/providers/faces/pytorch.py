"""Faces: Tensorflow"""
import os
import logging
import base64

from PIL import Image
from facenet_pytorch import MTCNN, InceptionResnetV1


class PyTorchModule:
    """PyTorchModule Faces"""

    def __init__(self, params: dict) -> None:
        os.environ['TORCH_HOME'] = '/'
        try:
            os.symlink('/models/faces/', '/checkpoints')
        except FileExistsError:
            pass
        except Exception as exp:
            logging.error(f'error creating symlink: {str(exp)}')
        self.prob_threshold = float(params['face_threshold'])
        self.det_model = MTCNN(keep_all=True)
        self.rec_model = InceptionResnetV1(pretrained=params['model'], classify=False)
        if self.rec_model:
            self.rec_model.eval()

    def detect(self, mediaitem_user_id: str, mediaitem_id: str, input_file: str) -> dict:
        """Detect faces for mediaitem"""
        embeddings, thumbnails = [], []
        faces, probs = self.det_model.forward(Image.open(input_file), return_prob=True, save_path=f'{mediaitem_id}.jpg')
        logging.debug(f'detected {len(faces) if faces is not None else None} \
                      faces for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        if faces is None or len(faces) == 0:
            return None

        result = self.rec_model.forward(faces)
        logging.debug(f'features calculated for {len(faces)} faces \
                       for user {mediaitem_user_id} mediaitem {mediaitem_id}')
        for idx, (res, prob) in enumerate(zip(result, probs)):
            if prob > self.prob_threshold:
                embeddings.append(res.tolist())
                thumbnails.append(base64.b64encode(Image.open(
                    f'{mediaitem_id}.jpg' if idx == 0 else f'{mediaitem_id}_{idx + 1}.jpg').tobytes()).decode('utf-8'))
        for filename in os.listdir('.'):
            if filename.startswith(mediaitem_id) and filename.endswith('.jpg'):
                file_path = os.path.join('.', filename)
                os.remove(file_path)

        logging.debug(f'got face embeddings for user {mediaitem_user_id} mediaitem {mediaitem_id}: {embeddings}')

        if len(embeddings) == 0:
            return None

        return dict({
            'userId': mediaitem_user_id,
            'id': mediaitem_id,
            'embeddings': embeddings,
            'thumbnails': thumbnails,
        })
