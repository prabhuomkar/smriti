# Smriti ML Faces
This is a reference for making ML models ready for inference to detect faces in images.  
The idea is for smriti to support multiple runtimes for different types of models.

## Providers
Following types of providers are available for running Face Detection:
- [PyTorch](https://pytorch.org/)

### PyTorch
Refer to [pytorch.py](pytorch.py) for downloading model files.

#### Save PyTorch Model Assets
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
