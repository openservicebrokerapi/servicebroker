# Table of Contents

- [Introduction](#introduction)
- [Platform Implementations](#platform-implementations)
- [Service Broker Libraries](#service-broker-libraries)
- [Other Libraries](#other-libraries)
- [Example Service Brokers](#example-service-brokers)

# Introduction

The Open Service Broker API project allows developers, ISVs, and SaaS vendors a
single, simple, and elegant way to deliver services to applications running
within cloud native platforms. To build a Service Broker, you must implement the
required endpoints as defined in the [API specification](spec.md).

This getting started guide contains some useful links, libraries and example
implementations to help you get started as quickly as possible. If you have
any questions, please feel free to join the
[Weekly Call](https://github.com/openservicebrokerapi/servicebroker/wiki/Weekly-Call)
where the community will be happy to help.

If you would like to add additional information, implementations or libraries to
this guide, please open a Pull Request.

Please note that the Open Service Broker API community does not make any
statement as to the validity, stability or compliance of any projects or tools
linked to in this guide.

# Platform Implementations

## Cloud Foundry

The [Cloud Controller](https://github.com/cloudfoundry/cloud_controller_ng)
component in Cloud Foundry is responsible for registering Service Brokers and
interacting with them on behalf of users. To deploy Cloud Foundry, please follow
the deployment guide for the
[cf-deployment](https://github.com/cloudfoundry/cf-deployment) project.

## Kubernetes

The [Service Catalog](https://github.com/kubernetes-incubator/service-catalog)
project is responsible for integrating Service Brokers to the Kubernetes
ecosystem. The project has its own
[Special Interest Group (SIG)](https://github.com/kubernetes/community/tree/master/sig-service-catalog)
and can be deployed by following
[this walkthrough guide](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/walkthrough.md).

# Service Broker Libraries

[brokerapi](https://github.com/pivotal-cf/brokerapi):
A Go package for building Open Service Broker API Service Brokers.

[Spring Cloud - Cloud Foundry Service Broker](https://github.com/spring-cloud/spring-cloud-cloudfoundry-service-broker):
This implements the REST contract for service brokers and the artifacts are
published to the Spring Maven repository. This greatly simplifies development:
include a single dependency in Gradle, implement interfaces, and configure. A
sample implementation has been provided for
[MongoDB](https://github.com/spring-cloud-samples/cloudfoundry-service-broker).

# Other Libraries

[go-open-service-broker-client](https://github.com/pmorie/go-open-service-broker-client):
This library is a golang client for communicating with service brokers,
useful for Platform developers.

# Example Service Brokers

## Go

[Asynchronous Service Broker for AWS EC2](https://github.com/cloudfoundry-samples/go_service_broker):
This Service Broker implements support for the
[Asynchronous Service Operations](https://docs.cloudfoundry.org/services/api.html#asynchronous-operations),
and calls AWS APIs to provision EC2 VMs.

[Storage Service Operations](https://github.com/opensds/nbp/tree/master/service-broker),
for OpenSDS to provision storage as a service.

[Open Service Broker for Azure](https://github.com/Azure/open-service-broker-azure):
This Service Broker implements support for Azure cloud services.

[GitHub Repository service](https://github.com/cloudfoundry-samples/github-service-broker-ruby):
This is designed to be an easy-to-read example of a service broker, with
complete documentation, and comes with a demo app that uses the service.
The Service Broker can be deployed as an application to any Cloud Foundry instance
or hosted elsewhere. The service broker uses GitHub as the service back end.

[MySQL database service](https://github.com/cloudfoundry/cf-mysql-release):
This Service Broker and its accompanying MySQL server are designed to be deployed
together as a BOSH release. [BOSH](https://github.com/cloudfoundry/bosh) is
used to deploy or upgrade the release, monitors the health of running
components, and restarts or recreates unhealthy VMs. The Service Broker code alone
can be found [here](https://github.com/cloudfoundry/cf-mysql-broker).

## Java

[MySQL Java Broker](https://github.com/cloudfoundry-community/cf-mysql-java-broker):
A Java port of the Ruby-based
[MySQL broker](https://github.com/cloudfoundry/cf-mysql-broker).

[Swisscom Open Service Broker](https://github.com/swisscom/open-service-broker):
enables platforms such as Cloud Foundry & Kubernetes to provision and manage
services. It is built in a modular way and one can host multiple services.
Open Service Broker offers extra functionality regarding billing,
backup/restore on top of the Open Service Broker API.
