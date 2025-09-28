"""ML Faces using ONNX model"""
import sys
import os
import urllib.request
import zipfile

import cv2
import numpy as np
import onnxruntime as ort


VERSION=os.getenv("VERSION", "dev").replace(".", "")
DET_MODELS_ZIPNAME = "faces_det.zip"
REC_MODELS_ZIPNAME = "faces_rec.zip"

def export_and_save():
    """Download required models for detection and recognition, export and save for inference"""
    print("downloading face detection model")
    urllib.request.urlretrieve("https://www.dropbox.com/scl/fo/9y86d4qb2nmkdtlf61aw0/AMc1Qfk6TtkbZMkhj8bBOVI"+
                               "?rlkey=2ztw1h9cvj9bsfc2niu42znr5&st=8u3eyse2&dl=1", DET_MODELS_ZIPNAME)
    with zipfile.ZipFile(DET_MODELS_ZIPNAME) as z:
        z.extractall("faces_det")
    os.remove(DET_MODELS_ZIPNAME)
    print("exporting and saving face detection model")
    print("all models are already exported in ONNX format!")
    print("downloading face recognition model")
    urllib.request.urlretrieve("https://www.dropbox.com/scl/fo/pubfb5wkiv5c5iucmlyou/ABKeExgXz5r-KMk5VUw2fPg"+
                               "?rlkey=djw1796psv1ncqun0kgswmkr9&st=ruowqv0e&dl=1", REC_MODELS_ZIPNAME)
    with zipfile.ZipFile(REC_MODELS_ZIPNAME) as z:
        z.extractall("faces_rec")
    os.remove(REC_MODELS_ZIPNAME)
    print("exporting and saving face recognition model")
    print("all models are already exported in ONNX format!")

def nms(dets, thresh=0.4):
    """non maximum supression with threshold"""
    x1 = dets[:, 0]
    y1 = dets[:, 1]
    x2 = dets[:, 2]
    y2 = dets[:, 3]
    areas = (x2 - x1 + 1) * (y2 - y1 + 1)
    order = dets[:, 4].argsort()[::-1]
    keep = []
    while order.size > 0:
        i = order[0]
        keep.append(i)
        xx1, yy1 = np.maximum(x1[i], x1[order[1:]]), np.maximum(y1[i], y1[order[1:]])
        xx2, yy2 = np.minimum(x2[i], x2[order[1:]]), np.minimum(y2[i], y2[order[1:]])
        w, h = np.maximum(0.0, xx2 - xx1 + 1), np.maximum(0.0, yy2 - yy1 + 1)
        inter = w * h
        ovr = inter / (areas[i] + areas[order[1:]] - inter)
        order = order[np.where(ovr <= thresh)[0] + 1]
    return keep

def load_and_run(sample="example.jpg"):
    """Loads the saved onnx models and runs sample image"""
    # input
    img = cv2.imread(sample)
    input_size = (640, 640)
    im_ratio = float(img.shape[0]) / img.shape[1]
    model_ratio = float(input_size[1]) / input_size[0]
    if im_ratio > model_ratio:
        new_height = input_size[1]
        new_width = int(new_height / im_ratio)
    else:
        new_width = input_size[0]
        new_height = int(new_width * im_ratio)
    det_scale = float(new_height) / img.shape[0]
    resized_img = cv2.resize(img, (new_width, new_height))
    det_img = np.zeros( (input_size[1], input_size[0], 3), dtype=np.uint8 )
    det_img[:new_height, :new_width, :] = resized_img
    blob = cv2.dnn.blobFromImage(det_img, 1.0/128.0, input_size, (127.5, 127.5, 127.5), swapRB=True)
    # detection
    for model in os.listdir("faces_det"):
        session = ort.InferenceSession(f"faces_det/{model}", providers=["CPUExecutionProvider"])
        print(model, "inputs:", [(inp.name, inp.shape) for inp in session.get_inputs()])
        print(model, "outputs:", [(out.name, out.shape) for out in session.get_outputs()])
        input_name = session.get_inputs()[0].name
        outputs = session.run(None, {input_name: blob})
        strides = [8, 16, 32]
        for idx, stride in enumerate(strides):
            scores = outputs[idx][0].reshape(-1)
            bboxes = outputs[idx+3][0].reshape(-1, 4) * stride
            h, w = blob.shape[2] // stride, blob.shape[3] // stride
            grid = np.stack(np.mgrid[:h, :w][::-1], axis=-1).astype(np.float32)
            anchors = (grid * stride).reshape(-1, 2)
            anchors = np.repeat(anchors, 2, axis=0).reshape(-1, 2)
            x1, y1 = anchors[:, 0] - bboxes[:, 0], anchors[:, 1] - bboxes[:, 1]
            x2, y2 = anchors[:, 0] + bboxes[:, 2], anchors[:, 1] + bboxes[:, 3]
            bboxes = np.stack([x1, y1, x2, y2], axis=-1) / det_scale
            keep_inds = scores >= 0.5
            scores, bboxes = scores[keep_inds], bboxes[keep_inds]
            dets = np.hstack([bboxes, scores[:, None]])
            keep = nms(dets, thresh=0.4)
            bboxes, scores = dets[keep, :4], dets[keep, 4]
            img = cv2.imread(sample)
            for idx, (box, score) in enumerate(zip(bboxes, scores)):
                if score > 0.8:
                    x1, y1, x2, y2 = box.astype(int)
                    face = img[y1-1:y2-1, x1-1:x2-1]
                    cv2.imwrite(f"result/face_{model.replace(".onnx", "")}_{stride}_{idx}.jpg", face)
    # recognition
    for model in os.listdir("faces_rec"):
        session = ort.InferenceSession(f"faces_rec/{model}", providers=["CPUExecutionProvider"])
        print(model, "inputs:", [(inp.name, inp.shape) for inp in session.get_inputs()])
        print(model, "outputs:", [(out.name, out.shape) for out in session.get_outputs()])
        input_name = session.get_inputs()[0].name
        for img_path in os.listdir("result"):
            img = cv2.imread(f"result/{img_path}")
            img = cv2.resize(img, (112, 112))
            img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
            img = np.transpose(img, (2, 0, 1)).astype(np.float32)
            img = (img / 255.0 - 0.5) / 0.5
            img_np = np.expand_dims(img, axis=0)
            outputs = session.run(None, {input_name: img_np})
            print(outputs[0].shape)

if __name__ == "__main__":
    args = sys.argv
    if len(args) > 1:
        if args[1] == "save":
            export_and_save()
            sys.exit(0)
        if args[1] == "run":
            if len(args)  == 3:
                load_and_run(args[2])
            else:
                load_and_run()
            sys.exit(0)
    print("provide a valid arg: save OR run")
