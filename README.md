# knot-cloud-storage

A web server that receives messages forwarded by devices and stores them in a database. The API provided by this service is specified below.

    POST /data

Stores the device message received from the webhook service.

## Parameters
### Header Parameters
Field | Required | Description
--- | --- | ---
x-meshblu-route | Y | Message route (shows the device that has sent it).
date | Y | Date on which the message was sent.

### Body Parameters

Field | Required | Description
--- | --- | ---
data | Y | The message data sent by the device. This data includes topic, payload and recipient devices.

## Example
### Request

    POST https://storage.knot.cloud/data

#### Request Body
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

### Response
    200 OK

---

    GET /data

Get all the device messages.

## Parameters
### Header Parameters
In order to get the device messages, you need to be authenticated as its owner (gateway or user).

Field | Required | Description
--- | --- | ---
auth_id | Y | Device ID.
auth_token | Y | Device token.

### URI Parameters
Field | Required | Description
--- | --- | ---
orderBy | N | The field used to order.
order | N | Ascending (1) or descending (-1) order, default=1.
skip | N | The number of data to skip (returns skip + 1), default=0.
take | N | The maximum number of data that you want from skip + 1 (the number is limited to 100), default=10.
startDate | N | The start date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).
finishDate | N | The finish date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).

## Example
### Request

    GET https://storage.knot.cloud/data?take=15&orderBy=timestamp&order=1

### Response
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

## Parameters
### Header Parameters
In order to get the device messages, you need to be authenticated as its owner (gateway or user).

Field | Required | Description
--- | --- | ---
auth_id | Y | Device ID.
auth_token | Y | Device token.

### URI Parameters
Field | Required | Description
--- | --- | ---
orderBy | N | The field used to order.
order | N | Ascending (1) or descending (-1) order, default=1.
skip | N | The number of data to skip (returns skip + 1), default=0.
take | N | The maximum number of data that you want from skip + 1 (the number is limited to 100), default=10.
startDate | N | The start date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).
finishDate | N | The finish date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).

## Example
### Request

    GET https://storage.knot.cloud/data/cc5429a29afcd158?startDate=2019-03-18 13:00:00

### Response
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

## Parameters
### Header Parameters
In order to get the device messages, you need to be authenticated as its owner (gateway or user).

Field | Required | Description
--- | --- | ---
auth_id | Y | Device ID.
auth_token | Y | Device token.

### URI Parameters
Field | Required | Description
--- | --- | ---
orderBy | N | The field used to order.
order | N | Ascending (1) or descending (-1) order, default=1.
skip | N | The number of data to skip (returns skip + 1), default=0.
take | N | The maximum number of data that you want from skip + 1 (the number is limited to 100), default=10.
startDate | N | The start date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).
finishDate | N | The finish date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).

## Example
### Request

    GET https://storage.knot.cloud/data/cc5429a29afcd158/sensor/2

### Response
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
