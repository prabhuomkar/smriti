"""Things: PyTorch"""
import logging
import torch
from torchvision import transforms
from PIL import Image


class PyTorchModule:
    """PyTorchModule Things"""

    def __init__(self, files: list[str]) -> None:
        self.module = torch.jit.load(f'/models/{files[0]}')

    def classify(self, mediaitem_user_id: str, mediaitem_id: str, input_file: str) -> dict:
        """Classify things for mediaitem"""
        transform = transforms.ToTensor()
        input_tensor = transform(Image.open(input_file)).unsqueeze(0)

        res = self.module.forward(input_tensor)
        logging.debug(f'thing for user {mediaitem_user_id} mediaitem {mediaitem_id}: {res}')

        return dict({
            'userId': mediaitem_user_id,
            'id': mediaitem_id,
            'name': next(iter(res)),
        })
