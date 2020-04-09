# knot-cloud-storage

A web server that receives messages forwarded by devices and stores them in a database. This service is currently being rewritten in golang to become suitable for running on constrained boards. The legacy instructions can be seen [here](legacy/LEGACY.md).

## Contents

- [Basic installation and usage](#basic-installation-and-usage)
  - [Requirements](#requirements)
  - [Configuration](#configuration)
  - [Setup](#setup)
  - [Compiling and Running](#compiling-and-running)
- [Docker installation and usage](#docker-installation-and-usage)
  - [Building and Running](#building-and-running)
    - [Production](#production)
    - [Development](#development)
- [Verify service health](#verify-service-health)

## Basic installation and usage

### Requirements

- Go version 1.13+.
- Be sure the local packages binaries path is in the system's `PATH` environment variable:

```bash
export PATH=$PATH:<your_go_workspace>/bin
```

### Configuration

You can set the `ENV` environment variable to `development` and update the `internal/config/development.yaml` when necessary. On the other way, you can use environment variables to configure your installation. In case you are running the published Docker image, you'll need to stick with the environment variables.

The configuration parameters are the following (the environment variable name is in parenthesis):

- `server`
  - `port` (`SERVER_PORT`) **Number** Server port number. (Default: 80)

### Setup

```bash
make tools
```

### Compiling and running

```bash
make run
```

> You can use the `make watch` command to run the application on watching mode, allowing it to be restarted automatically when the code changes.

## Docker installation and usage

Make sure you have [Docker Engine](<https://docs.docker.com/install/>) installed.

### Building and Running

#### Production

A container is specified at `docker/Dockerfile`. To use it, execute the following steps:

1. Build the image:

    ```bash
    docker build . -f docker/Dockerfile -t cesarbr/knot-cloud-storage:dev-go
    ```

2. Create a file containing the configuration as environment variables.

3. Run the container:

    ```bash
    docker run --env-file knot-cloud-storage.env -ti cesarbr/knot-cloud-storage:dev-go
    ```

#### Development

A development container is specified at `docker/Dockerfile-dev`. To use it, execute the following steps:

1. Build the image:

    ```bash
    docker build . -f docker/Dockerfile-dev -t cesarbr/knot-cloud-storage-go:dev-go
    ```

2. Create a file containing the configuration as environment variables.

3. Run the container:

    ```bash
    docker run --env-file knot-cloud-storage.env -p 8080:80 -v `pwd`:/usr/src/app -ti cesarbr/knot-cloud-storage-go:dev-go
    ```

The first argument to -v must be the root of this repository, so if you are running from another folder, replace `pwd` with the corresponding path.

This will start the server with auto-reload.

## Verify service health

### Verify

```bash
curl http://<hostname>:<port>/healthcheck
```
