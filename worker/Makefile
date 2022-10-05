run:
	@python main.py

test: 
	@pytest .

lint:
	@pylint *.py components/...

proto:
	@python3 -m grpc_tools.protoc -I../protos --python_out=. --grpc_python_out=. ../protos/*.proto