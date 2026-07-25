# Installation

## Docker

- Deploying with Docker is plain and simple.
  Refer to [docker-compose.yml](https://github.com/prabhuomkar/smriti/blob/master/docker-compose.yaml) for setting up services.
- Run all the services using:

```bah
docker-compose up -d
```

### Deploy using Docker Swarm

- [Docker Swarm](https://docs.docker.com/engine/swarm/) is the preferred approach for now to deploy Smriti.
- Environment variables can be set in [deploy/docker-swarm/docker-compose.yaml](https://github.com/prabhuomkar/smriti/blob/master/deploy/docker-swarm/docker-compose.yaml).
- Run following command to start services:

```bash
docker swarm init
docker swarm join --token <token>
docker compose up -d
```

## Kubernetes

Coming Soon!
