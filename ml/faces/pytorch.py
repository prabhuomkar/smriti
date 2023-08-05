"""ML Faces PyTorch Model"""
import sys
import os
import urllib.request

from PIL import Image
from facenet_pytorch import MTCNN, InceptionResnetV1


def download_and_save():
    """Download models for detection, recognition and classification"""
    print('downloading and saving pytorch models')
    urllib.request.urlretrieve('https://github.com/timesler/facenet-pytorch/releases/download/v2.2.9/20180402-114759-vggface2.pt', 'vggface2.pt')

def load_and_run(sample='example.jpg'):
    """Loads the saved pytorch models and runs sample image"""
    print('loading and running pytorch models')
    os.environ['TORCH_HOME'] = os.getcwd()
    try:
        os.symlink(os.getcwd(), f'{os.getcwd()}/checkpoints')
    except:
        pass
    det_model = MTCNN(keep_all=True)
    rec_model = InceptionResnetV1(pretrained='vggface2', classify=True).eval()
    faces, probs = det_model(Image.open(sample), return_prob=True)
    result = rec_model(faces)
    for res, prob in zip(result, probs):
        print(res.shape, prob)

if __name__ == '__main__':
    args = sys.argv
    if len(args) > 1:
        if args[1] == 'save':
            download_and_save()
            exit(0)
        if args[1] == 'run':
            if len(args)  == 3:
                load_and_run(args[2])
            else:
                load_and_run()
            exit(0)
    print('provide a valid arg: save OR run')
