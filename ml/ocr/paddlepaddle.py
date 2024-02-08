"""ML OCR Paddle Model"""
import sys
import os
import urllib.request
import tarfile
import subprocess
import shutil

import cv2


MODEL_BASE_URL='https://paddleocr.bj.bcebos.com/'

def download_and_save():
    """Download models for detection, recognition and classification"""
    print('downloading and saving paddle ocr models')
    urllib.request.urlretrieve(f'{MODEL_BASE_URL}PP-OCRv3/english/en_PP-OCRv3_det_infer.tar', 'det_infer.tar')
    with tarfile.open('det_infer.tar', 'r') as tar:
        tar.extractall('.')
    if os.path.exists('det_infer'):
        shutil.rmtree('det_infer')
    os.rename('en_PP-OCRv3_det_infer', 'det_infer')
    os.remove('det_infer.tar')
    urllib.request.urlretrieve(f'{MODEL_BASE_URL}PP-OCRv3/english/en_PP-OCRv3_rec_infer.tar', 'rec_infer.tar')
    with tarfile.open('rec_infer.tar', 'r') as tar:
        tar.extractall('.')
    if os.path.exists('rec_infer'):
        shutil.rmtree('rec_infer')
    os.rename('en_PP-OCRv3_rec_infer', 'rec_infer')
    os.remove('rec_infer.tar')
    urllib.request.urlretrieve(f'{MODEL_BASE_URL}dygraph_v2.0/ch/ch_ppocr_mobile_v2.0_cls_infer.tar', 'cls_infer.tar')
    with tarfile.open('cls_infer.tar', 'r') as tar:
        tar.extractall('.')
    if os.path.exists('cls_infer'):
        shutil.rmtree('cls_infer')
    os.rename('ch_ppocr_mobile_v2.0_cls_infer', 'cls_infer')
    os.remove('cls_infer.tar')

    commands = [
        "paddle2onnx --model_dir det_infer "
        "--model_filename inference.pdmodel "
        "--params_filename inference.pdiparams "
        "--save_file det_onnx/model.onnx "
        "--opset_version 10 "
        "--input_shape_dict=\"{'x':[-1,3,-1,-1]}\" "
        "--enable_onnx_checker True",

        "paddle2onnx --model_dir rec_infer "
        "--model_filename inference.pdmodel "
        "--params_filename inference.pdiparams "
        "--save_file rec_onnx/model.onnx "
        "--opset_version 10 "
        "--input_shape_dict=\"{'x':[-1,3,-1,-1]}\" "
        "--enable_onnx_checker True",

        "paddle2onnx --model_dir cls_infer "
        "--model_filename inference.pdmodel "
        "--params_filename inference.pdiparams "
        "--save_file cls_onnx/model.onnx "
        "--opset_version 10 "
        "--input_shape_dict=\"{'x':[-1,3,-1,-1]}\" "
        "--enable_onnx_checker True"
    ]
    for command in commands:
        subprocess.run(command, shell=True)
    shutil.rmtree('det_infer')
    shutil.rmtree('rec_infer')
    shutil.rmtree('cls_infer')
    urllib.request.urlretrieve('https://raw.githubusercontent.com/PaddlePaddle/PaddleOCR/release/2.7/ppocr/utils/en_dict.txt', 'rec_onnx/en_dict.txt')
    
def load_and_run(sample):
    """Loads the saved paddle ocr models and runs sample image"""
    print('loading and running paddle ocr models')
    import fastdeploy as fd
    default_option = fd.RuntimeOption()
    default_option.use_ort_backend()
    det_model = fd.vision.ocr.DBDetector(
            model_file='det_onnx/model.onnx', 
            runtime_option=default_option, model_format=fd.ModelFormat.ONNX)
    cls_model = fd.vision.ocr.Classifier(
        model_file='cls_onnx/model.onnx', 
        runtime_option=default_option, model_format=fd.ModelFormat.ONNX)
    rec_model = fd.vision.ocr.Recognizer(
        model_file='rec_onnx/model.onnx', 
        label_path='rec_infer/en_dict.txt', runtime_option=default_option, model_format=fd.ModelFormat.ONNX)
    ocr = fd.vision.ocr.PPOCRv3(det_model=det_model, cls_model=cls_model, rec_model=rec_model)
    result = ocr.predict(cv2.imread(sample))
    print(result.text)

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
