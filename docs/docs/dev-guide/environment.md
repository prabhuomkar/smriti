# Environment
The following guide will help you in setting up the development environment to work on Pensieve.

## Prerequisites
Make sure, before you start any development, following things are installed and available on your system of choice: 
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
- [Common Sense](https://en.wikipedia.org/wiki/Common_sense)

## Getting Codebase
Git clone the repository using: 
```
git clone git@github.com:prabhuomkar/pensieve.git
```

## Setup

### API
- Install [Golang 1.18](https://go.dev/dl/) or above
- Install [golangci-lint 1.49.0](https://golangci-lint.run/) or above
- Run the following command to finish your API setup
```
make setup-api
```

### Worker
- Install [Python 3.10](https://www.python.org/downloads/) or above
- Install [pylint 2.15.0](https://pypi.org/project/pylint/) or above
- Run the following command to finish your Worker setup
```
make setup-worker
```