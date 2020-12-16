# GPS Tracking Service

The GPS Tracking microservice is responsible for knowing the GPS coordinates
of the active fleet at any given time.

Data Sources report fleet telemetry to the tracking service.

## Flags

There are a few flags that control the execute of the service.

| Flag | Description | Default |
|---|---|---|
|object-ttl|Sets the expiration time for recorded objects|60 seconds|
|datastore|The datastore to use for objects|inmemdb|
|addr|interface and port to bind the service too|'0.0.0.0:5000'


## Container

The GPS Tracking service is packaged as a [container](https://hub.docker.com/r/scbunn/gps-tracking-service) on DockerHub.

### Tags

`latest` : Tracks the `master branch`

### Execution

```
$ docker run --rm -it \
    --name gps-tracking-service \
    scbunn/gps-tracking-service:tag -datastore inmemdb -object-ttl 30s

```

## Routes

| Method | Route | Description |
|---|---|---|
|GET|/metrics|Prometheus Exposition Formatted Metrics|
|GET|/health/liveness|Health check to determine if the container is alive|
|GET|/health/readiness|Health check to determine if the container is ready to take traffic|
|GET|/api/v1/location/:id|Retrieve the telemetry of a specific fleet object by id|
|GET|/api/v1/location/|Retrive a list of all fleet object's telemetry|
|POST|/api/v1/location/|Add/Update a fleet objects telemetry|

Fleet objects are ephemeral.  When the service recieves new telemetry about
an object it will either update the existing data or add a new object if one
does not exist.

### TTL Expiration

Fleet telemetry is expired after a given duration.  The GPS Tracking Service is
designed to track active fleet members only.  Objects that have not refreshed
their current telemetry will be expired from the service.

### Example Input Payload

```json
{
  "source": "sensor-collector-1",
  "objectId": "unique-id-to-source",
  "status": "object status (optional)",
  "posistion": {
    "latitude": 127.123,
    "longitude": -42.567,
    "elevation": 0, (optional)
  }
}
```

## Intended Bugs and Breaks

This service is designed to be used as a traning and testing tool.  A number of
ineffencies and bugs have been introduced in specific versions to demonstrate
failure scenerios or a specific troubleshooting or performance tool.

Below is a list of knows issues.

TBD
