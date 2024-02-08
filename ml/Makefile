install:
	@pip install -r requirements.txt

all:
	@echo "ℹ️ Creating classification models..."
	@cd classification && python3 pytorch.py save
	@echo "✅ Created classification models"
	@echo "ℹ️ Creating faces models..."
	@cd faces && python3 pytorch.py save
	@echo "✅ Created faces models"
	@echo "ℹ️ Creating ocr models..."
	@cd ocr && python3 paddlepaddle.py save
	@echo "✅ Created ocr models"
	@echo "ℹ️ Creating search models..."
	@cd search && python3 pytorch.py save
	@echo "✅ Created search models"