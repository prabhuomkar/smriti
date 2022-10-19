run:
	@python3 src/main.py

test: 
	@python3 -m pytest

lint:
	@pylint src/*.py src/components/...

proto:
	@python3 -m grpc_tools.protoc -I../protos --python_out=. --grpc_python_out=. ../protos/*.proto