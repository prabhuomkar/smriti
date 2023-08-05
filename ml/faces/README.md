# Smriti ML OCR
This is a reference for making ML models ready for inference to detect text in images.  
The idea is for smriti to support multiple runtimes for different types of models.

## Providers
Following types of providers are available for running ML Classification:
- [PaddlePaddle](https://github.com/PaddlePaddle/PaddleOCR)

### PaddlePaddle
Refer to [paddlepaddle.py](paddlepaddle.py) for downloading model files.

#### Save PaddleOCR Model Assets
```
python3 paddlepaddle.py save
```

#### Run Inference
```
python3 paddlepaddle.py run /path/to/example.jpg
```

### Tensorflow
Coming Soon!

### ONNX
Coming Soon!
