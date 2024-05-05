"""ML OCR Paddle Model"""
import sys
import os
import urllib.request
import shutil

import cv2
from paddleocr_convert import PaddleOCRModelConvert
from rapidocr_onnxruntime import RapidOCR


MODEL_BASE_URL='https://paddleocr.bj.bcebos.com/'

def download_and_save():
    """Download models for detection, recognition and classification"""
    print('downloading and saving paddle ocr models')
    urllib.request.urlretrieve(f'{MODEL_BASE_URL}PP-OCRv3/english/en_PP-OCRv3_det_infer.tar', 'det_infer.tar')
    urllib.request.urlretrieve(f'{MODEL_BASE_URL}PP-OCRv3/english/en_PP-OCRv3_rec_infer.tar', 'rec_infer.tar')
    urllib.request.urlretrieve(f'{MODEL_BASE_URL}dygraph_v2.0/ch/ch_ppocr_mobile_v2.0_cls_infer.tar', 'cls_infer.tar')
    urllib.request.urlretrieve('https://raw.githubusercontent.com/PaddlePaddle/PaddleOCR/release/2.7/ppocr/utils/en_dict.txt', 'rec_en_dict.txt')
    converter = PaddleOCRModelConvert()
    converter('det_infer.tar', 'det_onnx')
    converter = PaddleOCRModelConvert()
    converter('cls_infer.tar', 'cls_onnx')
    converter = PaddleOCRModelConvert()
    converter('rec_infer.tar', 'rec_onnx', txt_path='rec_en_dict.txt')
    os.remove('det_infer.tar')
    os.remove('cls_infer.tar')
    os.remove('rec_infer.tar')
    os.remove('rec_en_dict.txt')
    os.rename('det_onnx/det_infer/det_infer.onnx', 'det_onnx/model.onnx')
    os.rename('cls_onnx/cls_infer/cls_infer.onnx', 'cls_onnx/model.onnx')
    os.rename('rec_onnx/rec_infer/rec_infer.onnx', 'rec_onnx/model.onnx')
    shutil.rmtree('det_onnx/det_infer')
    shutil.rmtree('cls_onnx/cls_infer')
    shutil.rmtree('rec_onnx/rec_infer')
    
def load_and_run(sample='example.jpg'):
    """Loads the saved paddle ocr models and runs sample image"""
    print('loading and running paddle ocr models')
    engine = RapidOCR(rec_model_path='rec_onnx/model.onnx', det_model_path='det_onnx/model.onnx', cls_model_path='cls_onnx/model.onnx')
    result = engine(cv2.imread(sample))
    print(result)
    detections = result[0]
    text = ''
    for detection in detections:
        text += ' ' + detection[1]
    text = text.strip()
    print(text)

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
