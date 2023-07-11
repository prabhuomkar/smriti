"""ML OCR Paddle Model"""
import sys
import os
import urllib.request
import tarfile

from paddleocr import PaddleOCR


MODEL_BASE_URL='https://paddleocr.bj.bcebos.com/'

def download_and_save():
    """Download models for detection, recognition and classification"""
    print('downloading and saving paddle ocr models')
    urllib.request.urlretrieve(f'{MODEL_BASE_URL}PP-OCRv3/english/en_PP-OCRv3_det_infer.tar', 'det_infer.tar')
    with tarfile.open('det_infer.tar', 'r') as tar:
        tar.extractall('.')
    os.rename('en_PP-OCRv3_det_infer', 'det_infer')
    os.remove('det_infer.tar')
    urllib.request.urlretrieve(f'{MODEL_BASE_URL}PP-OCRv3/english/en_PP-OCRv3_rec_infer.tar', 'rec_infer.tar')
    with tarfile.open('rec_infer.tar', 'r') as tar:
        tar.extractall('.')
    os.rename('en_PP-OCRv3_rec_infer', 'rec_infer')
    os.remove('rec_infer.tar')
    urllib.request.urlretrieve(f'{MODEL_BASE_URL}dygraph_v2.0/ch/ch_ppocr_mobile_v2.0_cls_infer.tar', 'cls_infer.tar')
    with tarfile.open('cls_infer.tar', 'r') as tar:
        tar.extractall('.')
    os.rename('ch_ppocr_mobile_v2.0_cls_infer', 'cls_infer')
    os.remove('cls_infer.tar')

def load_and_run(sample):
    """Loads the saved paddle ocr models and runs sample image"""
    print('loading and running paddle ocr models')
    ocr = PaddleOCR(use_angle_cls=True, lang='en', det_model_dir='./det_infer', rec_model_dir='./rec_infer', cls_model_dir='./cls_infer')
    result = ocr.ocr(sample)
    words = []
    for res in result:
        for line in res:
            words.append(line[1])
    print(words)

if __name__ == '__main__':
    args = sys.argv
    if len(args) > 1:
        if args[1] == 'save':
            download_and_save()
            exit(0)
        if args[1] == 'run':
            load_and_run(args[2])
            exit(0)
    print('provide a valid arg: save OR run')
