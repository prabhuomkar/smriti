setup:
	@echo "[setup]: Installing requirements..."; \
		pip install -r requirements.txt; \
		echo "[setup]: Installing drivers..."
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
test-all:
	@echo "testing all features"
	behave