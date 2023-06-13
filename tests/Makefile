setup:
	@echo "[setup]: Installing requirements..."; \
		pip install -r requirements.txt; \
		echo "[setup]: Installing drivers..."; \
		echo "[setup]: Starting services..."; \
		docker compose -f ../docker-compose.yaml up -d
test-users:
	@echo "testing feature: users"
	behave -i users.feature
test-auth:
	@echo "testing feature: auth"
	behave -i auth.feature
test-mediaitems:
	@echo "testing feature: mediaitems"
	behave -i mediaitems.feature
test-albums:
	@echo "testing feature: albums"
	behave -i albums.feature
test-library:
	@echo "testing feature: library"
	behave -i library.feature
test-explore:
	@echo "testing feature: explore"
	behave -i explore.feature
test-search:
	@echo "testing feature: search"
	@echo "feature not developed yet!"
test-raw:
	@echo "testing feature: raw mediaitems"
	behave -i raw.feature
test-all:
	@echo "testing all features"
	behave --tags=~long
clear:
	@echo "clearing uploaded files"
	@sudo rm -rf ../storage/originals/*
	@sudo rm -rf ../storage/previews/*
	@sudo rm -rf ../storage/thumbnails/*
	@sudo find ../storage/* -maxdepth 0 ! -name 'originals' ! -name 'previews' ! -name 'thumbnails' -exec rm -rf {} +
