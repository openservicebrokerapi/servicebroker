# Service Broker POC 

[![Build Status](https://travis-ci.org/cncf/servicebroker.svg?branch=master)](https://travis-ci.org/cncf/servicebroker)
[![Go Report Card](https://goreportcard.com/badge/github.com/cncf/servicebroker)](https://goreportcard.com/report/github.com/cncf/servicebroker)

This repo has the PoC code the CNCF Service Broker WG.

All PRs must be signed with a DCO.

## Building

To build everything just run: `make` and that should leave you with a
`service_controller` executable in the `k8s/service_controller/` directory
along with a Docker image called`service_controller`.

## Running

`docker run -ti -p 10000:10000 service_controller` should bring up a Service Controller
listening on port 10000. So, using:
```
curl localhost:10000
```
should be able to hit it.

## Testing

`make test-unit`
