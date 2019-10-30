# knot-cloud-storage

A web server that receives messages forwarded by devices and stores them in a database. This service is currently being rewritten in golang to become suitable for running on constrained boards. The legacy instructions can be seen [here](LEGACY.md).


## Installation and usage

### Requirements

*   Go version 1.13+.
*   Be sure the local packages binaries path is in the system's `PATH` environment variable:

```bash
$ export PATH=$PATH:<your_go_workspace>/bin
```

### Configuration

You can set the `ENV` environment variable to `development` and update the `internal/config/development.yaml` when necessary. On the other way, you can use environment variables to configure your installation. In case you are running the published Docker image, you'll need to stick with the environment variables.

The configuration parameters are the following (the environment variable name is in parenthesis):

*   `server`
    *   `port` (`SERVER_PORT`) **Number** Server port number. (Default: 80)

### Setup

```bash
make tools
make deps
```

### Compiling and running

```bash
make run
```

### Verify

```bash
curl http://<hostname>:<port>/healthcheck
```
