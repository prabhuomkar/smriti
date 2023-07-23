setup: setup-api setup-worker setup-docs setup-tests
setup-api:
	@cd api; \
		echo "[setup-api]: Verifying go modules..."; \
		go mod verify; \
		echo "[setup-api]: Running linter..."; \
		make lint; \
		echo "[setup-api]: Running unit tests..."; \
		make test
setup-worker:
	@cd worker; \
		echo "[setup-worker]: Installing requirements"; \
		make install; \
		make test-install; \
		echo "[setup-worker]: Running linter..."; \
		make lint; \
		echo "[setup-worker]: Running unit tests..."; \
		make test
setup-docs:
	@cd docs; \
		echo "[setup-docs]: Installing dependencies"; \
		npm install
setup-tests:
	@cd tests; \
		echo "[setup-tests]: Installing dependencies"; \
		make setup
start:
	@docker compose up -d
stop:
	@docker compose down