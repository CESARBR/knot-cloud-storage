# knot-cloud-storage

A web server that receives messages forwarded by devices and stores them in a database. The API provided by this service is specified below.

    POST /data

Stores the device message received from the webhook service.

## Parameters
### Body Parameters

Field | Required | Description
--- | --- | ---
metadata | Y | Required to identify which device has sent the message.
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

Get all the device messages received from the webhook service.

## Parameters
### URI Parameters
Field | Required | Description
--- | --- | ---
orderBy | N | The field used to order.
order | N | Ascending (1) or descending (-1) order, default=1.
limit | N | The maximum number of data that you want, default=10.
start | N | The start date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).
finish | N | The finish date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).

## Example
### Request

    GET https://storage.knot.cloud/data?limit=15 orderBy=timestamp&order=1

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
### URI Parameters
Field | Required | Description
--- | --- | ---
orderBy | N | The field used to order.
order | N | Ascending (1) or descending (-1) order, default=1.
limit | N | The maximum number of data that you want, default=10.
start | N | The start date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).
finish | N | The finish date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).

## Example
### Request

    GET https://storage.knot.cloud/data/188824f0-28c4-475b-ab36-2505402bebcb?start=2019-03-18 13:00:00

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
### URI Parameters
Field | Required | Description
--- | --- | ---
orderBy | N | The field used to order.
order | N | Ascending (1) or descending (-1) order, default=1.
limit | N | The maximum number of data that you want, default=10.
start | N | The start date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).
finish | N | The finish date that you want your set of data (format=YYYY-MM-DD HH:MM:SS).

## Example
### Request

    GET https://storage.knot.cloud/data/188824f0-28c4-475b-ab36-2505402bebcb/sensor/2

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
