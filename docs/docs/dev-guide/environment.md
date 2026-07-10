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

- Install [Golang 1.26](https://go.dev/dl/) or above
- Install [golangci-lint 2.12.2](https://golangci-lint.run/) or above
- Run the following command to finish your API setup

```
make setup-api
```

### Worker

- Install [CMake 3.28+](https://cmake.org/download/) or above
- Install [cpplint 2.0.0](https://pypi.org/project/cpplint/) or above
- Install [clang-tidy 19.1.7](https://clang.llvm.org/extra/clang-tidy/) or above
- Install [clang-format 19.1.7](https://clang.llvm.org/docs/ClangFormat.html) or above

```
make setup-worker
```

### Docs

- Install [Node 18.10](https://nodejs.org/en/download/) or above
- Run the following command to finish your Docs setup

```
make setup-docs
```

### Tests

- Install [Python 3.12](https://www.python.org/downloads/) or above
- Run the following command to finish your Tests setup

```
make setup-tests
```
