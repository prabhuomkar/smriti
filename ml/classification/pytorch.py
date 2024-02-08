"""ML Classification TorchVision Model"""
import os
import sys

from PIL import Image
import torch
from torchvision import transforms
from torchvision.models import get_model, get_model_weights


VERSION=os.getenv('VERSION', 'dev').replace('.', '')
FILE_NAME = f'classification_v{VERSION}.pt'
MODEL_NAME = 'efficientnet_v2_s' # can be any torchvision classification model

class SmritiClassificationPyTorchModule(torch.nn.Module):
    """EfficientNet TorchScript Module"""
    def __init__(self, model, weights) -> None:
        super(SmritiClassificationPyTorchModule, self).__init__()
        self.model = model
        self.model.eval()
        self.transforms = weights.transforms()
        self.categories = weights.meta['categories']

    def forward(self, img_tensor, topk: int=5):
        """Forward Pass"""
        input_tensor = self.transforms(img_tensor)
        output = self.model(input_tensor)
        probabilities = torch.nn.functional.softmax(output, dim=1)
        top_prob, top_class = torch.topk(probabilities, topk)
        return dict({self.categories[top_class[0][idx].item()]: top_prob[0][idx].item() \
            for idx in range(0, int(topk))})

def script_and_save():
    """Initialize pytorch model with weights, script it and save the torchscript module"""
    print('scripting and saving torchscript module')
    scripted_module = torch.jit.script(SmritiClassificationPyTorchModule(get_model(MODEL_NAME, weights='DEFAULT'), get_model_weights(MODEL_NAME).DEFAULT))
    scripted_module.save(FILE_NAME)

def load_and_run(sample):
    """Loads the saved torchscript module and runs sample image"""
    print('loading and running torchscript module')
    model = torch.jit.load(FILE_NAME)
    img = Image.open(sample)
    transform = transforms.ToTensor()
    input_tensor = transform(img).unsqueeze(0)
    print(model(input_tensor))

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
