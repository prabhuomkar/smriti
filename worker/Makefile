install:
	@pip install -r requirements.txt

run:
	@python3 -m src.main

test-install:
	@pip install pytest pytest-cov requests-mock pytest-asyncio

test: 
	@python3 -m pytest

cover: 
	@coverage run --source src -m pytest
	@coverage report --omit=src/protos/*,*__.py
	@coverage xml

cover-html: cover
	@coverage html --omit=src/protos/*,*__.py

lint:
	@pylint *

proto:
	@mkdir -p ../src/protos
	@cp -r ../protos/* ../src/protos/
	@python3 -m grpc_tools.protoc -I../ --python_out=./ --pyi_out=./ --grpc_python_out=./ ../src/protos/*.proto
	@rm -rf ../src