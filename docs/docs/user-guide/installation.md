# Installation Guide

## Docker
- Deploying with Docker is plain and simple. 
Refer to [docker-compose.yml](https://github.com/prabhuomkar/carousel/blob/master/docker-compose.yaml) for setting up services.
- Run all the services using:
```bah
docker-compose up -d
```

### Deploy on Fly.io
- [Fly.io](https://fly.io/) is a promising alternative to Heroku. One can deploy docker containers with configuration of your choice.  
Users can find sample configuration [here](https://github.com/prabhuomkar/carousel/tree/master/infra/deployments/fly) for API and Worker.
- Update the TOML configuration with your specific app names
```toml
app = "tonystark-carousel-api"
```
- Environment variables can be set using [Fly Secrets](https://fly.io/docs/reference/secrets/).
```bash
flyctl secrets set CAROUSEL_WORKER_HOST=tonystark-carousel-worker.fly.dev 
```

## Kubernetes
TODO(omkar)

## Raspberry Pi
TODO(omkar)