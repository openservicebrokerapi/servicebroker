# Table of Contents

- [Introduction](#introduction)
- [Platform Implementations](#platform-implementations)
- [Quickstarts](#quickstarts)
- [Service Broker Libraries](#service-broker-libraries)
- [Other Libraries](#other-libraries)
- [Example Service Brokers](#example-and-community-service-brokers)
- [Related community components](#related-community-components)

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

# Quickstarts

[`osb-starter-pack`](https://github.com/pmorie/osb-starter-pack):
A go project that lets you easily deploy and iterate on a new service broker.
Uses the [`osb-broker-lib`](https://github.com/pmorie/osb-broker-lib) and
[`go-open-service-broker-client`](https://github.com/pmorie/go-open-service-broker-client)
projects.

# Service Broker Libraries

[`brokerapi`](https://github.com/pivotal-cf/brokerapi):
A Go package for building Open Service Broker API Service Brokers.

[Spring Cloud Open Service Broker](https://spring.io/projects/spring-cloud-open-service-broker):
Spring Cloud Open Service Broker provides a framework based on Spring Boot that
enables you to quickly create a service broker for your own managed service on
platform that support the Open Service Broker API.

[`osb-broker-lib`](https://github.com/pmorie/osb-broker-lib):
A go library that provides the REST API implementation for the OSB API. Users
implement an interface that uses the types from the
[`go-open-service-broker-client`](https://github.com/pmorie/go-open-service-broker-client).

[`Open Service Broker API for .NET`](https://github.com/AXOOM/OpenServiceBroker):
.NET libraries for client and server implementations of the Open Service Broker API. The client library allows you to call Service Brokers that implement the API using idiomatic C# interfaces and type-safe DTOs. The server library implements the API for you using ASP.NET Core. You simply need to provide implementations for a few interfaces, shielded from the HTTP-related details.

[spring-cloud-app-broker](https://github.com/spring-cloud/spring-cloud-app-broker)
Spring Cloud App Broker is a framework for building Spring Boot applications that implement the Open Service Broker API to dynamically deploy Cloud Foundry applications.

[Cloud service broker](https://github.com/pivotal/cloud-service-broker/)
This service broker uses Terraform to provision and bind services.

[Open Broker API](https://openbrokerapi.readthedocs.io/en/latest/)
A Python library that provides the REST API implementation for the OSB API. Users
implement an interface.

# Other Libraries

[`go-open-service-broker-client`](https://github.com/pmorie/go-open-service-broker-client):
This library is a golang client for communicating with service brokers,
useful for Platform developers.

# Example and Community Service Brokers

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

[On Demand Service Broker](https://github.com/pivotal-cf/on-demand-service-broker):
This is a generic service broker for BOSH deployed services. The broker
deploys any BOSH release on demand. It is used by the
[Redis for PCF](https://www.cloudfoundry.org/the-foundry/redis-for-pcf/), 
[MySQL for PCF](https://pivotal.io/platform/services-marketplace/data-management/mysql), 
[RabbitMQ for PCF](https://www.cloudfoundry.org/the-foundry/rabbitmq-for-pcf/)
and 
[Pivotal Cloud Cache](https://pivotal.io/platform/services-marketplace/data-management/pivotal-cloud-cache) 
service brokers. The On Demand Broker is open source, and typically deployed via
[its BOSH release](https://github.com/pivotal-cf/on-demand-service-broker-release).

[Open Service Broker for Huawei Cloud](https://github.com/huaweicloud/huaweicloud-service-broker):
This Service Broker implements support for Huawei cloud services.

[World's Simplest Service Broker](https://github.com/cloudfoundry-community/worlds-simplest-service-broker) 
This service broker shares the same binding credentials with everyone - for Kubernetes and Cloud Foundry 

[Logs-service-broker](https://github.com/orange-cloudfoundry/logs-service-broker)
A logs service broker to forward CloudFoundry [syslog drains](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#log-drain) logs 
to one or more syslog backends (http, tcp, or udp). This supports log parsing and transformation. 

[Cf-redis-broker](https://github.com/pivotal-cf/cf-redis-broker/)
A service broker for a shared redis cluster.

[cf-rabbitmq-multitenant-broker](https://github.com/pivotal-cf/cf-rabbitmq-multitenant-broker-release/)
a multi-tenant RabbitMQ service broker for Cloud Foundry.

[mongodb-open-service-broker](https://github.com/orange-cloudfoundry/mongodb-boshrelease/tree/master/src/mongodb-open-service-broker)
A service broker for a mongodb cluster

## Java

[MySQL Java Broker](https://github.com/cloudfoundry-community/cf-mysql-java-broker):
A Java port of the Ruby-based
[MySQL broker](https://github.com/cloudfoundry/cf-mysql-broker).

[Swisscom Open Service Broker](https://github.com/swisscom/open-service-broker):
enables platforms such as Cloud Foundry & Kubernetes to provision and manage
services. It is built in a modular way and one can host multiple services.
Open Service Broker offers extra functionality regarding billing,
backup/restore on top of the Open Service Broker API.

[Static credentials Broker](https://github.com/orange-cloudfoundry/static-creds-broker/) 
This service broker serves statically configured data (catalog and service bindings)

[Cassandra broker](https://github.com/orange-cloudfoundry/cassandra-boshrelease/tree/master/src/cassandra-open-service-broker)
A service broker creating service instances as cassandra keyspaces and service bindings as cassandra roles.

## Ruby

[cf-mysql-broker](https://github.com/cloudfoundry-attic/cf-mysql-broker)
A service broker for a shared mariadb cluster

# Related community components

[OSB CMDB](https://github.com/orange-cloudfoundry/osb-cmdb)
A configuration management database for Service Brokers.  
This enables sharing of service brokers among multiple OSB client platforms by providing inventory, events, quotas, analytics, etc...

[OSB Reverse proxy](https://github.com/orange-cloudfoundry/osb-reverse-proxy)
A reverse proxy for open service broker endpoints, providing recent remote access to logs of recent requests 

[Overview broker](https://github.com/cloudfoundry/overview-broker)
For the purpose of testing OSB client platforms, a service broker that provides an overview of its 
service instances and bindings, its dashboard provides the full OSB API calls received.  

[Peripli Service Manager](https://peripli.github.io/)
The Service Manager is a component that manages Open Service Broker API compatible service brokers. 
It can enforce polices on service brokers, instances and binding and enables cross-platform capabilities such as cross-platform service instance sharing.

[Sec-group-broker-filter](https://github.com/orange-cloudfoundry/sec-group-broker-filter)
A service broker designed to be chained in front-of other service brokers and dynamically open Cloud Foundry security groups 
to let apps bound to service instance emmit outgoing traffic to IP addresses returned in the chained service instances credentials.

[Subway](https://github.com/cloudfoundry-community/cf-subway)
Subway is a multiplexing service broker that allows you to scale out single node brokers 
