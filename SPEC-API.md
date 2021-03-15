### GET /goscope/api/application-name
  
```json
{
  "applicationName": "Data Sync"
}
```
  
### GET /goscope/api/requests?offset=0
  
```json
{
    "applicationName": "Data Sync",
    "data": [
      {
        "method": "POST",
        "path": "/hooks/sync",
        "time": 1596117389,
        "uid": "ff6b8d82-f7aa-4d91-5854-9fc1ad0e3ab6",
        "responseStatus": 200
      }
    ],
    "entriesPerPage": 50
}
```
  
### GET /goscope/api/requests/:uuid
  
```json
{
    "applicationName": "Data Sync",
    "data": {
      "request": {
        "body": "",
        "clientIP": "127.0.0.1",
        "headers": "{\n    \"Accept\": [\n        \"image/webp,*/*\"\n    ],\n    \"Accept-Encoding\": [\n        \"gzip, deflate\"\n    ],\n    \"Accept-Language\": [\n        \"en-US,en;q=0.5\"\n    ],\n    \"Cache-Control\": [\n        \"max-age=0\"\n    ],\n    \"Connection\": [\n        \"keep-alive\"\n    ],\n    \"Cookie\": [\n        \"io=KTQMAozFS7YxI9SDAAAI\"\n    ],\n    \"Referer\": [\n        \"http://localhost:7005/goscope/requests\"\n    ],\n    \"User-Agent\": [\n        \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:80.0) Gecko/20100101 Firefox/80.0\"\n    ]\n}",
        "host": "localhost:7005",
        "method": "GET",
        "path": "/assets/logo.svg",
        "referrer": "http://localhost:7005/goscope/requests",
        "time": 1597502010,
        "uid": "4b52e4e7-79af-44bd-5b18-a7dd9a5801cd",
        "url": "/assets/logo.svg",
        "userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:80.0) Gecko/20100101 Firefox/80.0"
      },
      "response": {
        "body": "",
        "clientIP": "127.0.0.1",
        "headers": "{\n    \"Content-Length\": [\n        \"40289\"\n    ],\n    \"Content-Type\": [\n        \"svg\"\n    ]\n}",
        "path": "/assets/logo.svg",
        "size": 0,
        "status": "200",
        "time": 1597502010,
        "requestUID": "",
        "uid": "9daa100b-256c-4faf-7ab5-ea04f9283ca1"
      }
    }
}
```
  
### GET /goscope/api/logs?offset=0
  
```json
{
    "applicationName": "Data Sync",
    "data": [
      {
        "error": "utils.go:69: invalid character '\u0026' looking for beginning of object key string\n",
        "time": 1597405457,
        "uid": "1529492a-67e0-44c5-4599-c65685a9ba10"
      }
    ],
    "entriesPerPage": 50
}
```
  
### GET /goscope/logs/:uuid
  
```json
{
    "applicationName": "Data Sync",
    "data": {
      "logDetails": {
        "error": "utils.go:69: invalid character '\u0026' looking for beginning of object key string\n",
        "time": 1597405294,
        "uid": "8ad41605-2172-49d7-7f18-414e8c7e4a46"
      }
    }
}
```
  
### GET /goscope/api/info
  
```json
{
    "applicationName": "Data Sync",
    "cpu": {
      "coreCount": "8 Cores",
      "modelName": "Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz"
    },
    "disk": {
      "freeSpace": "297.63 GB",
      "mountPath": "/",
      "partitionType": "apfs",
      "totalSpace": "465.63 GB"
    },
    "host": {
      "hostOS": "darwin",
      "hostPlatform": "darwin",
      "hostname": "averageflow.fritz.box",
      "kernelArch": "x86_64",
      "kernelVersion": "19.6.0",
      "uptime": "0.41 hours"
    },
    "memory": {
      "availableMemory": "6.58 GB",
      "totalMemory": "16.00 GB",
      "usedSwap": "0.00%"
    }
}
```
  
### POST /goscope/api/search/requests?offset=0

To learn more about the search types look in [the correct enum](internal/repository/Types.go)

Request Body:

```json
{
  "query": "test",
  "searchType": "0"
}
```
Response Body:
```json
{
  "applicationName": "Data Sync",
  "data": [
    {
      "method": "POST",
      "path": "/hooks/sync",
      "time": 1596117850,
      "uid": "eda14d71-2671-4868-41a2-6cbfaebab6eb",
      "responseStatus": 200
    }
  ],
  "entriesPerPage": 50
}
```


### POST /goscope/api/search/logs?offset=0

Request Body:
```json
{
  "query": "test"
}
```

Response Body:
```json
{
  "applicationName": "Data Sync",
  "data": [
    {
      "error": "utils.go:69: invalid character '\u0026' looking for beginning of object key string\n",
      "time": 1597405457,
      "uid": "1529492a-67e0-44c5-4599-c65685a9ba10"
    }
  ],
  "entriesPerPage": 50
}
```