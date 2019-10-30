# knot-cloud-storage

A web server that receives messages forwarded by devices and stores them in a database.

## Installation and usage

This service is part of the KNoT Cloud and requires the a subset of its service to work.

### Configuration

Either create a [configuration file](https://github.com/lorenwest/node-config/wiki/Configuration-Files) in the `config` (`./config/local.json` is ignored by Git in this repository) or use environment variables to configure your installation. In case you are running the published Docker image, you'll need to stick with the environment variables.

The configuration parameters are the following (the environment variable name is in parenthesis):

* `database` (`DB`) **String** Database service type. Currently only **MONGO**. (Default: **MONGO**)
* `databaseConfig`
  * `hostname` (`DB_HOSTNAME`) **String** Mongo hostname.
  * `port` (`DB_PORT`) **Number** Mongo port.
  * `name` (`DB_NAME`) **String** Mongo database name. (Default: **knot_storage**)
* `server`
  * `port` (`PORT`) **Number** Server port number. (Default: 80)
  * `publicKey` (`PUBLIC_KEY_BASE64`) **String** Publisher public key (See POST /data documentation for details)
* `meshblu`
  * `aliasLookupServerUri` (`MESHBLU_ALIAS_LOOKUP_SERVER_URI`) **String** Alias lookup service base URI.
  * `cacheRedisUri` (`MESHBLU_CACHE_REDIS_URI`) **String** URI of Redis server used by Meshblu cache.
  * `namespace` (`MESHBLU_NAMESPACE`) **String** Meshblu's namespace on Redis. (Default: **meshblu**)
  * `messagesNamespace` (`MESHBLU_MESSAGES_NAMESPACE`) **String** Meshblu's messages namespace on Redis. (Default: **messages**)
  * `redisUri` (`MESHBLU_REDIS_URI`) **String** URI of Redis server used by Meshblu.
  * `jobTimeoutSeconds` (`MESHBLU_JOB_TIMEOUT_SECONDS`) **Number** Job timeout in seconds (Default: 30)
  * `jobLogSampleRate` (`MESHBLU_JOB_LOG_SAMPLE_RATE`) **Number** Job sample rate (Default: 0)
  * `requestQueueName` (`MESHBLU_REQUEST_QUEUE_NAME`) **String** Meshblu's request queue name (Default: **v2:request:queue**)
  * `responseQueueName` (`MESHBLU_RESPONSE_QUEUE_NAME`) **String** Meshblu's response queue name (Default: **v2:response:queue**)

Only change the `MESHBLU_` parameters that have a default value if you know what you are doing.

### Build and run (local)

First, install the dependencies:

```
npm install --production
```

Then:

```
npm run build
npm start
```

### Build and run (local, development)

First, install the dependencies:

```
npm install
```

Then, start the server with auto-reload:

```
npm run start:watch
```

or, start the server in debug mode:

```
npm run start:debug
```

### Build and run (Docker, development)

A development container is specified at `Dockerfile-dev`. To use it, execute the following steps:

1. Build the image:

    ```
    docker build . -f Dockerfile-dev -t knot-cloud-storage-dev
    ```

1. Create a file containining the configuration as environment variables.
1. Run the container:

    ```
    docker run --env-file storage.env -p 4000:80 -v `pwd`:/usr/src/app -ti knot-cloud-storage-dev
    ```

The first argument to `-v` must be the root of this repository, so if you are running from another folder, replace `` `pwd` `` with the corresponding path.

This will start the server with auto-reload.

### Run (Docker)

Containers built from the master branch and the published tags in this repository are available on [DockerHub](https://hub.docker.com/r/cesarbr/knot-cloud-storage/).

1. Create a file containining the configuration as environment variables.
1. Run the container:

```
docker run --env-file storage.env -p 4000:80 -ti cesarbr/knot-cloud-storage
```

To verify if the service is running properly, execute:

```
curl http://<hostname>:<port>/healthcheck
```


## API

    POST /data

Stores the device message received from the webhook service.

### Parameters
#### Header Parameters

This endpoint uses the [HTTP Signature Scheme](https://github.com/joyent/node-http-signature/blob/master/http_signing.md#signature-authentication-scheme)
to authenticate the request.

Field | Required | Description
--- | --- | ---
x-meshblu-route | Y | Message route (shows the device that has sent it).
date | Y | Date on which the message was sent.
Authorization | Y | Contains a signature that can be verified using the public key in the configuration.

#### Body Parameters

Field | Required | Description
--- | --- | ---
data | Y | The message data sent by the device. This data includes topic, payload and recipient devices.

### Example
#### Request

    POST https://storage.knot.cloud/data

##### Request Body
```json
{
  "metadata": {
    "route": [
      {
        "from": "188824f0-28c4-475b-ab36-2505402bebcb",
        "to": "cf2497d2-7426-46c4-a229-ad789063bf88",
        "type": "broadcast.sent"
      },
      {
        "from": "188824f0-28c4-475b-ab36-2505402bebcb",
        "to": "cf2497d2-7426-46c4-a229-ad789063bf88",
        "type": "broadcast.received"
      },
      {
        "from": "cf2497d2-7426-46c4-a229-ad789063bf88",
        "to": "cf2497d2-7426-46c4-a229-ad789063bf88",
        "type": "broadcast.received"
      }
    ]
  },
  "data": {
    "topic": "data",
    "devices": ["*"],
    "payload": {
       "sensorId": 2,
       "value": 234
    }
  }
}
```

#### Response
    200 OK

---

    GET /data

Get all the device messages.

### Parameters
#### Header Parameters
In order to get the device messages, you need to be authenticated as its owner (gateway or user).

Field | Required | Description
--- | --- | ---
auth_id | Y | Device ID.
auth_token | Y | Device token.

#### URI Parameters
Field | Required | Description
--- | --- | ---
orderBy | N | The field used to order.
order | N | Ascending (1) or descending (-1) order, default=1.
skip | N | The number of data to skip (returns skip + 1), default=0.
take | N | The maximum number of data that you want from skip + 1 (the number is limited to 100), default=10.
startDate | N | The start date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).
finishDate | N | The finish date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).

### Example
#### Request

    GET https://storage.knot.cloud/data?take=15&orderBy=timestamp&order=1

#### Response
```json
[
    {
        "from": "188824f0-28c4-475b-ab36-2505402bebcb",
        "payload": {
            "sensorId": 2,
            "value": 234
        },
        "timestamp": "2019-03-18T12:48:05.569Z"
    },
    {
        "from": "188824f0-28c4-475b-ab36-2505402bebcb",
        "payload": {
            "sensorId": 1,
            "value": true
        },
        "timestamp": "2019-03-18T14:42:03.192Z"
    }
]
```

---

    GET /data/{device_id}

Get the messages received from the webhook service by a specific device.

### Parameters
#### Header Parameters
In order to get the device messages, you need to be authenticated as its owner (gateway or user).

Field | Required | Description
--- | --- | ---
auth_id | Y | Device ID.
auth_token | Y | Device token.

#### URI Parameters
Field | Required | Description
--- | --- | ---
orderBy | N | The field used to order.
order | N | Ascending (1) or descending (-1) order, default=1.
skip | N | The number of data to skip (returns skip + 1), default=0.
take | N | The maximum number of data that you want from skip + 1 (the number is limited to 100), default=10.
startDate | N | The start date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).
finishDate | N | The finish date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).

### Example
#### Request

    GET https://storage.knot.cloud/data/cc5429a29afcd158?startDate=2019-03-18 13:00:00

#### Response
```json
[
    {
        "from": "188824f0-28c4-475b-ab36-2505402bebcb",
        "payload": {
            "sensorId": 1,
            "value": true
        },
        "timestamp": "2019-03-18T14:42:03.192Z"
    }
]
```

---

    GET /data/{device_id}/sensor/${sensor_id}

Get the messages received from the webhook service by a specific device sensor.

### Parameters
#### Header Parameters
In order to get the device messages, you need to be authenticated as its owner (gateway or user).

Field | Required | Description
--- | --- | ---
auth_id | Y | Device ID.
auth_token | Y | Device token.

#### URI Parameters
Field | Required | Description
--- | --- | ---
orderBy | N | The field used to order.
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
        "timestamp": "2019-03-18T12:48:05.569Z"
    }
]
```
