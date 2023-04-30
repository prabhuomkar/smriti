# Folder Structure

- `/`: 
  - `Makefile`: Helpful commands for interacting with project
  - `docker-compose.yaml`: Containerized services for local integration tests
- `api`: contains the source code for the API server
- `docs`: contains the overall documentation of the project
- `infra`: contains examples for deployments e.g. Docker, Kubernetes, etc.
- `protos`: contains Protocol Buffer definitions for API and Worker services
- `tests`: contains code and test cases for Integration Tests
- `worker`: contains the source code for the Worker server with its individual components for Place Detection, Face Detection & Clustering and Things Detection
