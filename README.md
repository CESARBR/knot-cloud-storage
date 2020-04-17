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

## API

```bash
POST /data
```

Stores the device messages.

### Parameters
#### Head Parameters

Field | Required | Description
--- | --- | ---
auth_token | Y | User's authentication token

### Body Parameters

Field | Required | Description
--- | --- | ---
data | Y | The message data sent by the device. It includes both sender identification and payload.

### Example
#### Request

```bash
POST https://storage.knot.cloud/data
```

##### Request Body
```json
{
    "from": "188824f0-28c4-475b-ab36-2505402bebcb",
    "payload": {
        "sensor_id": 2,
        "value": 234
    }
}
```

##### Response
```
201 Created
```

---

```
GET /data
```
Get all the device messages.

### Parameters
#### Header Parameters
In order to get the device messages, you need to be authenticated as its owner.

Field | Required | Description
--- | --- | ---
auth_token | Y | User's authentication token

#### URI Parameters
Field | Required | Description
--- | --- | ---
order | N | Ascending (1) or descending (-1) order, default=1.
skip | N | The number of data to skip (returns skip + 1), default=0.
take | N | The maximum number of data that you want from skip + 1 (the number is limited to 100), default=10.
startDate | N | The start date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).
finishDate | N | The finish date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).

### Example
#### Request
```
GET https://storage.knot.cloud/data?take=15&order=1
```

#### Response
```json
[
    {
        "from": "188824f0-28c4-475b-ab36-2505402bebcb",
        "payload": {
            "sensorId": 2,
            "value": 234
        },
        "timestamp": "2020-04-10T10:32:26.456753-03:00"
    },
    {
        "from": "188824f0-28c4-475b-ab36-2505402bebcb",
        "payload": {
            "sensorId": 1,
            "value": true
        },
        "timestamp": "2020-04-13T13:56:23.596301-03:00"
    }
]
```

---

```
GET /data/{deviceID}
```
Get the messages by a specific device.

### Parameters
#### Header Parameters
In order to get the device messages, you need to be authenticated as its owner.

Field | Required | Description
--- | --- | ---
auth_token | Y | User's authentication token

#### URI Parameters
Field | Required | Description
--- | --- | ---
order | N | Ascending (1) or descending (-1) order, default=1.
skip | N | The number of data to skip (returns skip + 1), default=0.
take | N | The maximum number of data that you want from skip + 1 (the number is limited to 100), default=10.
startDate | N | The start date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).
finishDate | N | The finish date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).


### Example
#### Request
```
GET https://storage.knot.cloud/data/cc5429a29afcd158?startDate=2020-04-13 13:00:00
```

#### Response
```json
[
  {
        "from": "188824f0-28c4-475b-ab36-2505402bebcb",
        "payload": {
            "sensorId": 1,
            "value": true
        },
        "timestamp": "2020-04-13T13:56:23.596301-03:00"
    }
]
```

---

```
GET /data/{deviceID}/sensor/{sensorID}
```
Get the messages by a specific device sensor.

### Parameters
#### Header Parameters
In order to get the device messages, you need to be authenticated as its owner.

Field | Required | Description
--- | --- | ---
auth_token | Y | User's authentication token

#### URI Parameters
Field | Required | Description
--- | --- | ---
order | N | Ascending (1) or descending (-1) order, default=1.
skip | N | The number of data to skip (returns skip + 1), default=0.
take | N | The maximum number of data that you want from skip + 1 (the number is limited to 100), default=10.
startDate | N | The start date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).
finishDate | N | The finish date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).

### Example
#### Request

    GET https://storage.knot.cloud/data/cc5429a29afcd158/sensor/2

#### Response
```json
[
    {
        "from": "188824f0-28c4-475b-ab36-2505402bebcb",
        "payload": {
            "sensorId": 2,
            "value": 234
        },
        "timestamp": "2020-04-10T10:32:26.456753-03:00"
    }
]
```

---

```
DELETE /data
```
Delete all the device messages.

### Parameters
#### Header Parameters
In order to delete the messages, you need to be authenticated as a valid user. Only data from the things which are owned by this user will be removed.

Field | Required | Description
--- | --- | ---
auth_token | Y | User's authentication token

### Example
#### Request

    DELETE https://storage.knot.cloud/data

#### Response
```
200 OK
```

---

```
DELETE /data/{deviceID}
```
Delete all messages from a specific device.

### Parameters
#### Header Parameters
In order to delete the device messages, you need to be authenticated as its owner.

Field | Required | Description
--- | --- | ---
auth_token | Y | User's authentication token

### Example
#### Request

    DELETE https://storage.knot.cloud/data/cc5429a29afcd158

#### Response
```
200 OK
```