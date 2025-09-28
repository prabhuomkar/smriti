"""ML OCR PaddlePaddle Model"""
import sys
import os
import urllib.request
import tarfile

from paddleocr import PaddleOCR


MODEL_BASE_URL="https://paddle-model-ecology.bj.bcebos.com/paddlex/official_inference_model/paddle3.0.0/"
CATEGORY_MODELS = {
    "doc_orient": ["PP-LCNet_x1_0_doc_ori_infer.tar"],
    "text_unwrap": ["UVDoc_infer.tar"],
    "text_line_orient": ["PP-LCNet_x1_0_textline_ori_infer.tar", "PP-LCNet_x0_25_textline_ori_infer.tar"],
    "text_det": ["PP-OCRv5_server_det_infer.tar", "PP-OCRv5_mobile_det_infer.tar"],
    "text_rec": ["PP-OCRv5_server_rec_infer.tar", "PP-OCRv5_mobile_rec_infer.tar"],
}

def download_and_save():
    """Download models for detection, recognition and classification"""
    print("downloading and saving paddle ocr models")
    for category, models in CATEGORY_MODELS.items():
        for model in models:
            urllib.request.urlretrieve(f"{MODEL_BASE_URL}{model}", model)
            model_path = f"ocr_{category}"
            os.makedirs(model_path, exist_ok=True)
            with tarfile.open(model, "r:") as tar:
                tar.extractall(path=model_path)
                full_model_path = f"{model_path}/{model.replace('.tar', '')}"
                # remove redundant files
                for filename in os.listdir(full_model_path):
                    if "demo" in filename:
                        os.remove(full_model_path+"/"+filename)

    for models in CATEGORY_MODELS.values():
        for model in models:
            os.remove(model)

def load_and_run(sample="example.jpg"):
    """Loads the saved paddle ocr models and runs sample image"""
    print("loading and running paddle ocr models")
    ocr = PaddleOCR(
        doc_orientation_classify_model_name="PP-LCNet_x1_0_doc_ori",
        doc_orientation_classify_model_dir="./ocr_doc_orient/PP-LCNet_x1_0_doc_ori_infer",
        textline_orientation_model_name="PP-LCNet_x0_25_textline_ori",
        textline_orientation_model_dir="./ocr_text_line_orient/PP-LCNet_x0_25_textline_ori_infer",
        doc_unwarping_model_name="UVDoc",
        doc_unwarping_model_dir="./ocr_text_unwrap/UVDoc_infer",
        text_detection_model_name="PP-OCRv5_mobile_det",
        text_detection_model_dir="./ocr_text_det/PP-OCRv5_mobile_det_infer",
        text_recognition_model_name="PP-OCRv5_mobile_rec",
        text_recognition_model_dir="./ocr_text_rec/PP-OCRv5_mobile_rec_infer"
    )
    outputs = ocr.predict(sample)
    for output in outputs:
        output.save_to_json("ocr_output")

if __name__ == "__main__":
    args = sys.argv
    if len(args) > 1:
        if args[1] == "save":
            download_and_save()
            sys.exit(0)
        if args[1] == "run":
            if len(args)  == 3:
                load_and_run(args[2])
            else:
                load_and_run()
            sys.exit(0)
    print("provide a valid arg: save OR run")
