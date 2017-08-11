# gotime

### Requirements

* Go

### Usage

> If you have a Github token, provide it using query parameter `?token=<token>` . 

* `go run main.go`
* Get the latest activity in a Github repo:

```
curl -X GET \
  'http://localhost:8000/latest-activity/github/baopham/gotime' \
  -H 'content-type: application/json'
```

* Get the average response time for a Github repo:

```
curl -X GET \
  'http://localhost:8000/response-time/github/baopham/gotime' \
  -H 'content-type: application/json'
```
