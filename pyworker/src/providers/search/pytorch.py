"""Search: PyTorch"""
import logging

from PIL import Image
import torch
from transformers import AutoTokenizer, AutoImageProcessor


class PyTorchModule:
    """PyTorchModule Search"""

    def __init__(self, params: dict) -> None:
        self.tokenizer = AutoTokenizer.from_pretrained(f'/models/search/{params["tokenizer_dir"]}')
        self.processor = AutoImageProcessor.from_pretrained(f'/models/search/{params["processor_dir"]}')
        self.text_module = torch.jit.load(f'/models/search/{params["text_file"]}')
        self.vision_module = torch.jit.load(f'/models/search/{params["vision_file"]}')

    def generate_embedding(self, input_type: str, data: any):
        """Generate text embedding from text"""
        if input_type == 'text':
            input_tensor = self.tokenizer(data, padding=True, return_tensors='pt')
            res = self.text_module.forward(**input_tensor)
            if res is not None:
                res = res.tolist()
                logging.debug(f'generated text embedding: {res}')
                return res
            return None
        if data and data['type'] == 'photo':
            input_tensor = self.processor(Image.open(data['previewPath']), return_tensors='pt')
            res = self.vision_module.forward(**input_tensor)
            if res is not None:
                res = res.tolist()
                logging.debug(f'generated photo embedding: {res}')
                return [res]
            return []
        return []
