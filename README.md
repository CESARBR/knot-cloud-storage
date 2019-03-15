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
