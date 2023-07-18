"""Search: PyTorch"""
import logging
import torch
from transformers import AutoTokenizer


class PyTorchModule:
    """PyTorchModule Search"""

    def __init__(self, params: list[str]) -> None:
        self.module = torch.jit.load(f'/models/search/{params[0]}')
        self.tokenizer = AutoTokenizer.from_pretrained(f'/models/search/{params[1]}')

    def generate_embedding(self, text: str):
        """Generate text embedding from text"""
        input_tensor = self.tokenizer(text, padding=True, return_tensors='pt')

        res = self.module.forward(**input_tensor)
        if res is not None:
            res = res.tolist()

        logging.debug(f'generated embedding: {res}')

        return res
