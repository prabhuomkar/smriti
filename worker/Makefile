run:
	@python3 src/main.py

test: 
	@python3 -m pytest

cover: 
	@coverage run --source src -m pytest
	@coverage report --omit=*_pb2*.py*

cover-html: cover
	@coverage html --omit=*_pb2*.py*

lint:
	@pylint src/*.py src/components/...

proto:
	@python3 -m grpc_tools.protoc -I../protos --python_out=./src --pyi_out=./src --grpc_python_out=./src ../protos/*.proto