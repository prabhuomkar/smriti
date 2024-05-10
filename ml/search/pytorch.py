"""ML Search HuggingFace Model"""
import os
import sys

from PIL import Image
import torch
from transformers import AutoTokenizer, AutoImageProcessor, CLIPModel, BlipModel


VERSION=os.getenv('VERSION', 'dev').replace('.', '')
TEXT_FILE_NAME = f'search_text_v{VERSION}.pt'
VISION_FILE_NAME = f'search_vision_v{VERSION}.pt'
TOKENIZER_DIR_NAME = 'search_tokenizer'
PROCESSOR_DIR_NAME = 'search_processor'
MODEL_TYPE = 'blip' # options: blip, clip
MODEL_NAME = 'openai/clip-vit-base-patch32' if MODEL_TYPE == 'clip' else 'Salesforce/blip-image-captioning-base' # can be any huggingface model
MODEL = CLIPModel if MODEL_TYPE == 'clip' else BlipModel

class SmritiSearchTextPyTorchModule(torch.nn.Module):
    """Search TorchScript Text Module"""
    def __init__(self, model) -> None:
        super(SmritiSearchTextPyTorchModule, self).__init__()
        self.model = model
        self.model.eval()

    def forward(self, input_ids: torch.tensor, attention_mask: torch.tensor):
        """Forward Pass"""
        output = self.model.get_text_features(input_ids, attention_mask)
        return output[0]
    
class SmritiSearchVisionPyTorchModule(torch.nn.Module):
    """Search TorchScript Vision Module"""
    def __init__(self, model) -> None:
        super(SmritiSearchVisionPyTorchModule, self).__init__()
        self.model = model
        self.model.eval()

    def forward(self, pixel_values: torch.tensor):
        """Forward Pass"""
        output = self.model.get_image_features(pixel_values)
        return output[0]

def script_and_save():
    """Initialize pytorch model with weights, script it and save the torchscript module"""
    print('scripting and saving torchscript module')
    # text
    tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME)
    tokenizer.save_pretrained(TOKENIZER_DIR_NAME)
    text_inputs = tokenizer(["a photo of a cat"], padding=True, return_tensors="pt")
    traced_text_module = torch.jit.optimize_for_inference(torch.jit.trace(SmritiSearchTextPyTorchModule(MODEL.from_pretrained(MODEL_NAME, torch_dtype=torch.float16)),
                                         (text_inputs['input_ids'], text_inputs['attention_mask'])))
    traced_text_module.save(TEXT_FILE_NAME)
    tokenizer.save_pretrained(TOKENIZER_DIR_NAME)
    # vision
    processor = AutoImageProcessor.from_pretrained(MODEL_NAME)
    processor.save_pretrained(PROCESSOR_DIR_NAME)
    vision_inputs = processor(Image.open('example.jpg'), return_tensors="pt")
    traced_vision_module = torch.jit.optimize_for_inference(torch.jit.trace(SmritiSearchVisionPyTorchModule(MODEL.from_pretrained(MODEL_NAME, torch_dtype=torch.float16)),
                                           (vision_inputs['pixel_values'])))
    traced_vision_module.save(VISION_FILE_NAME)

def load_and_run(sample='example.jpg'):
    """Loads the saved torchscript module and runs sample image"""
    print('loading and running torchscript module')
    # text
    text_model = torch.jit.load(TEXT_FILE_NAME)
    tokenizer = AutoTokenizer.from_pretrained(TOKENIZER_DIR_NAME)
    input = tokenizer(sample, padding=True, return_tensors='pt')
    input_ids, attention_mask = input['input_ids'], input['attention_mask']
    print(text_model(input_ids, attention_mask))
    # vision
    vision_model = torch.jit.load(VISION_FILE_NAME)
    processor = AutoImageProcessor.from_pretrained(PROCESSOR_DIR_NAME)
    input = processor(Image.open('example.jpg'), return_tensors='pt')
    pixel_values = input['pixel_values']
    print(vision_model(pixel_values))

if __name__ == '__main__':
    args = sys.argv
    if len(args) > 1:
        if args[1] == 'save':
            script_and_save()
            exit(0)
        if args[1] == 'run':
            if len(args)  == 3:
                load_and_run(args[2])
            else:
                load_and_run()
            exit(0)
    print('provide a valid arg: save OR run')
