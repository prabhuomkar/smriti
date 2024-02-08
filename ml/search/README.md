# Smriti ML Search
This is a reference for making ML models ready for inference to search based on embeddings.  
The idea is for smriti to support multiple runtimes for different types of models.

## Providers
Following types of providers are available for running Search:
- [PyTorch](https://pytorch.org/)

### PyTorch
Refer to [pytorch.py](pytorch.py) for creating a [TorchScript](https://pytorch.org/docs/stable/jit.html) Module using a [HuggingFace](https://huggingface.co) model.

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
