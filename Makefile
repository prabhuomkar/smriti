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
		pip install -r requirements.txt; \
		echo "[setup-worker]: Running linter..."; \
		make lint; \
		echo "[setup-worker]: Running unit tests..."; \
		make test
