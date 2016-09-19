# Open Service Broker API - `Context` Profile

## Abstract

This document contains the recommended usage pattern for the `context`
property that is defined as part of the
[Open Service Broker API specification](spec.md).
While the `context` property is defined as an opaque JSON object, in practice
it is often useful and necessary for there to be an agreed upon set of
properties within the `context` object to share information between the
Platform and Service Brokers in order to enable the Service Brokers to perform
their job.  It is RECOMMENDED that implementations follow the guidelines
defined in this document in order to aide in interoperability across
various Platform deployments and Service Brokers.

## Table of Contents

- [Overview](#overview)
- [Notations and Terminology](#notations-and-terminology)
- [Context Property](#context-property)
- [Cloud Foundry](#cloud-foundry)
- [Kubernetes](#kubernetes)

## Overview

In the [Open Service Broker API specification](spec.md) there are certain
message flows that include a `context` property. This property is defined
as an opaque JSON object that is meant to contain contextual information
about the usage of the Service being acted upon. For example, it may include
organizational information (e.g. a Cloud Foundry `organization` GUID), or
the ID of the Application to which the Service will be bound.

While Platforms MAY include any information within the `context` property,
or none at all, is it strongly RECOMMENDED that they adhere to the guidelines
defined within this document.

While use of this profile is OPTIONAL, an implementation is not compliant
with this profile if it fails to satisfy one or more of the MUST, SHALL
or REQUIRED level requirements defined herein.

## Notations and Terminology

### Notational Conventions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED",  "MAY", and "OPTIONAL" in this document are to
be interpreted as described in [RFC 2119]( https://tools.ietf.org/html/rfc2119).

### Terminology

The terminology defined below are defined by the
[Open Service Broker API specification](spec.md) and are included here
for convinience. If there are any inconsistencies between their definitions in
this document and the [Open Service Broker API specification](spec.md), then
the [Open Service Broker API specification](spec.md) SHALL take precidence.

- *Platform*: The software that will manage the cloud environment into which
  Applications and Service Brokers are provisioned.  Users will not directly
  provision Services from Service Brokers, rather they will ask the Platform
  (ie. their cloud provider) to manage Services and interact with the
  Service Brokers for them.

- *Service*: A managed software offering that can be used by an Application.
  Typically, Services will expose some API that can be invoked to perform
  some action. However, there can also be non-interactive Services that can
  perform the desired actions without direct prompting from the Application.

- *Service Broker*: Service Brokers manage the lifecycle of Services. Platforms
  interact with Service Brokers to provision, bind and delete Service
  Instances.

- *Application*: The software that will make use of, be bound to, a Service.

## Context Property

The `context` property is defined as an opaque JSON object. Depending on
which Service Broker API is being invoked different properties MUST appear
within the `context` object.
The following sections define the set of properties that each Platform
defines and in which API call they MUST appear.

Aside from the Platform specific properties, there is one common property
called `platform` that MUST also appear within `context` to indicate which
Platform is being used.

The `platform` property MUST be a `string` and serialized as follows:<br>
```
"platform" : "platform-string-here"
```

Each section below will define the `platform` value that MUST be used.

### Cloud Foundry

*Platform Property Value*: `cloudfoundry`

The following properties are defined for usage within a Cloud Foundry
deployment:

- `organization_guid`<br>
  The organization GUID as defined by the Cloud Foundry
  specification/project. This property MUST be a string serialize as follows:
  ```
  "organization_guid" : "organization-guid-here"
  ```
  For example:
  ```
  "organization_guid" : "1113aa0-124e-4af2-1526-6bfacf61b111"
  ```

- `space_guid`<br>
  The space GUID as defined by the Cloud Foundry
  specification/project. This property MUST be a string serialized as follows:
  ```
  "space_guid": "space-guid-here"
  ```
  For example:
  ```
  "space_guid" : "aaaa1234-da91-4f12-8ffa-b51d0336aaaa"
  ```

The following table specifies which properties MUST appear in each API:

Request API | Properties
----------- | ----------
`POST /v2/service_plan_visibilities` | `organization_guid`
`PUT /v2/service_instances/:instance_id` | `organization_guid`<br>`space_guid`
`PATCH /v2/service_instances/:instance_id` | `organization_guid`<br>`space_guid`

Example:

The following example shows a `context` property that might appear as
part of a Cloud Foundry API call:
  ```
  "context" : {
    "platform" : "cloudfoundry",
    "organization_guid" : "1113aa0-124e-4af2-1526-6bfacf61b111",
    "space_guid" : "aaaa1234-da91-4f12-8ffa-b51d0336aaaa"
  }
  ```

### Kubernetes

*Platform Property Value*: `kubernetes`

- `namespace`: ...

