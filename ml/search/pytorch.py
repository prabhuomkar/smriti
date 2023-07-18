"""ML Search HuggingFace Model"""
import sys

import torch
from transformers import AutoTokenizer, CLIPTextModelWithProjection


VERSION='20230731'
FILE_NAME = f'search_v{VERSION}.pt'
TOKENIZER_DIR_NAME = 'search_tokenizer'
MODEL_NAME = 'openai/clip-vit-base-patch32' # can be any huggingface model

class SmritiSearchPyTorchModule(torch.nn.Module):
    """CLIP TorchScript Module"""
    def __init__(self, model) -> None:
        super(SmritiSearchPyTorchModule, self).__init__()
        self.model = model
        self.model.eval()

    def forward(self, input_ids: torch.tensor, attention_mask: torch.tensor):
        """Forward Pass"""
        output = self.model(input_ids, attention_mask)
        return output['text_embeds'][0]

def script_and_save():
    """Initialize pytorch model with weights, script it and save the torchscript module"""
    print('scripting and saving torchscript module')
    tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME)
    tokenizer.save_pretrained(TOKENIZER_DIR_NAME)
    inputs = tokenizer(["a photo of a cat"], padding=True, return_tensors="pt")
    traced_module = torch.jit.trace(SmritiSearchPyTorchModule(CLIPTextModelWithProjection.from_pretrained(MODEL_NAME)), (inputs['input_ids'], inputs['attention_mask']))
    traced_module.save(FILE_NAME)

def load_and_run(sample):
    """Loads the saved torchscript module and runs sample image"""
    print('loading and running torchscript module')
    model = torch.jit.load(FILE_NAME)
    tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME)
    input = tokenizer(sample, padding=True, return_tensors='pt')
    input_ids, attention_mask = input['input_ids'], input['attention_mask']
    print(model(input_ids, attention_mask))

if __name__ == '__main__':
    args = sys.argv
    if len(args) > 1:
        if args[1] == 'save':
            script_and_save()
            exit(0)
        if args[1] == 'run':
            load_and_run(args[2])
            exit(0)
    print('provide a valid arg: save OR run')
