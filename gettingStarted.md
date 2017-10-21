# Table of Contents

- [Sample Service Brokers](#sample-service-brokers)
- [Libraries](#libraries)

# Sample Service Brokers

The following example service broker implementations have been developed
as a starting point if you are developing your own service broker.

The Open Service Broker API does not make any statement as to the
validity, stability or compliance of any of them.

If you would like to add additional brokers to this list open a
a pull request against 
[this repository](https://github.com/openservicebrokerapi/servicebroker)
and edit [this file](gettingStarted.md).

## Ruby

[GitHub Repository service](https://github.com/cloudfoundry-samples/github-service-broker-ruby):
This is designed to be an easy-to-read example of a service broker, with
complete documentation, and comes with a demo app that uses the service.
The broker can be deployed as an application to any Cloud Foundry instance
or hosted elsewhere. The service broker uses GitHub as the service back end.

[MySQL database service](https://github.com/cloudfoundry/cf-mysql-release):
This broker and its accompanying MySQL server are designed to be deployed
together as a BOSH release. [BOSH](https://github.com/cloudfoundry/bosh) is
used to deploy or upgrade the release, monitors the health of running
components, and restarts or recreates unhealthy VMs. The broker code alone
can be found [here](https://github.com/cloudfoundry/cf-mysql-broker).

## Java

[Spring Cloud - Cloud Foundry Service Broker](https://github.com/spring-cloud/spring-cloud-cloudfoundry-service-broker):
This implements the REST contract for service brokers and the artifacts are
published to the Spring Maven repository. This greatly simplifies development:
include a single dependency in Gradle, implement interfaces, and configure. A
sample implementation has been provided for
[MongoDB](https://github.com/spring-cloud-samples/cloudfoundry-service-broker).

[MySQL Java Broker](https://github.com/cloudfoundry-community/cf-mysql-java-broker):
A Java port of the Ruby-based
[MySQL broker](https://github.com/cloudfoundry/cf-mysql-broker) above.

## Go

[Asynchronous Service Broker for AWS EC2](https://github.com/cloudfoundry-samples/go_service_broker):
This broker implements support for the 
[Asynchronous Service Operations](https://docs.cloudfoundry.org/services/api.html#asynchronous-operations),
and calls AWS APIs to provision EC2 VMs.

[Storage Service Operations](https://github.com/leonwanghui/opensds-broker/tree/master),
for OpenSDS to provision storage as a service.

# Libraries

## Go

[Go Client Library](https://github.com/pmorie/go-open-service-broker-client):
This library is a golang client for communicating with service brokers,
useful for Platform developers.
