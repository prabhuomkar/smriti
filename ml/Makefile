install:
	@pip install -r requirements.txt

all:
	@echo "ℹ️ Creating faces models..."
	@cd faces && python3 export_onnx.py save
	@echo "✅ Created faces models"
	@echo "ℹ️ Creating ocr models..."
	@cd ocr && python3 export_paddlepaddle.py save
	@echo "✅ Created ocr models"
	@echo "ℹ️ Creating search models..."
	@cd search && python3 pytorch.py save
	@echo "✅ Created search models"