# Environment Setup

The following guide will help you in setting up the development environment to work on Smriti.

## Prerequisites

Make sure, before you start any development, following things are installed and available on your system of choice:

- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
- [Protocol Buffers](https://protobuf.dev/installation/)
- [Common Sense](https://en.wikipedia.org/wiki/Common_sense)

## Getting Codebase

Git clone the repository using:

```
git clone git@github.com:prabhuomkar/smriti.git
```

## Setup

### API

- Install [Golang](https://go.dev/dl/) or above
- Install [golangci-lint](https://golangci-lint.run/) or above
- Run the following command to finish your `api` setup

```
make setup-api
```

### Worker

- Install [CMake](https://cmake.org/download/) or above
- Install [cpplint](https://pypi.org/project/cpplint/) or above
- Install [clang-tidy](https://clang.llvm.org/extra/clang-tidy/) or above
- Install [clang-format](https://clang.llvm.org/docs/ClangFormat.html) or above
- Run the following command to finish your `worker` setup

```
make setup-worker
```

### Docs

- Install [Python](https://www.python.org/downloads/) or above
- Run the following command to finish your docs setup

```
make setup-docs
```

### Tests

- Install [Python](https://www.python.org/downloads/) or above
- Run the following command to finish your tests setup

```
make setup-tests
```
