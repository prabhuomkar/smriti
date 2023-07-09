# Smriti ML Classification
This is a reference for making ML models ready for inference.  
The idea is for smriti to support multiple runtimes for different types of models.

## Providers
Following types of providers are available for running ML Classification:
- [PyTorch](https://pytorch.org/)
- [Tensorflow](https://www.tensorflow.org/)
- [ONNX](https://onnxruntime.ai/)

### PyTorch
Refer to [pytorch.py](pytorch.py) for creating a [TorchScript](https://pytorch.org/docs/stable/jit.html) Module using a [TorchVision](https://pytorch.org/vision/stable/index.html) model.

#### Save TorchScript Module
```
python3 pytorch.py save
```

#### Run Inference
```
python3 pytorch.py run /path/to/example.jpg
```

### Tensorflow
Coming Soon!

### ONNX
Coming Soon!
