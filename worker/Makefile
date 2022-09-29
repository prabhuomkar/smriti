run:
	@python main.py

test: 
	@python test ./...

lint:
	@pylint *.py components/...

proto:
	@python3 -m grpc_tools.protoc -I../proto --python_out=. --grpc_python_out=. ../proto/*.proto