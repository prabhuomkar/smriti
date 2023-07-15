"""Search: PyTorch"""
import logging
import torch
from transformers import AutoTokenizer


class PyTorchModule:
    """PyTorchModule Search"""

    def __init__(self, params: list[str]) -> None:
        self.module = torch.jit.load(f'/models/{params[0]}')
        self.model_name = params[1]

    def generate_embedding(self, text: str):
        """Generate text embedding from text"""
        tokenizer = AutoTokenizer.from_pretrained(self.model_name)
        input_tensor = tokenizer(text, padding=True, return_tensors='pt')

        res = self.module.forward(**input_tensor)
        logging.debug(f'generated embedding: {res}')

        return res
