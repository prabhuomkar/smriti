# Smriti ML Faces

This is a reference for making ML models ready for inference to detect and recgonize faces in images. The idea is for
smriti to support multiple runtimes or libraries for different types of models.

## Providers

Following types of providers are available for running Face Detection:

- [ONNX](https://onnxruntime.ai/)

### ONNX

Refer to [export_onnx.py](export_onnx.py) for downloading, exporting, saving and running the model.

#### Save PyTorch Model Assets

```
python3 pytorch.py save
```

#### Run Inference

```
python3 pytorch.py run /path/to/example.jpg
```
