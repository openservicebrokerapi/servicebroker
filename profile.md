# Open Service Broker API Profile (master - might contain changes that are not yet released)

## Abstract

The [Open Service Broker API specification](spec.md) allows for extensions
and variations based on the environments in which it is being used; this
document contains the suggested usage pattern for some of those variants.

## Table of Contents

- [Notations and Terminology](#notations-and-terminology)
- [Context Object](#context-object)
  - [Context Object Properties](#context-object-properties)
  - [Cloud Foundry](#cloud-foundry)
  - [Kubernetes](#kubernetes)
- [Service Metadata](#service-metadata)

## Notations and Terminology

### Notational Conventions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED",  "MAY", and "OPTIONAL" in this document are to
be interpreted as described in [RFC 2119]( https://tools.ietf.org/html/rfc2119).

### Terminology

The terminology defined below are defined by the
[Open Service Broker API specification](spec.md) and are included here
for convenience. If there are any inconsistencies between their definitions in
this document and the [Open Service Broker API specification](spec.md), then
the [Open Service Broker API specification](spec.md) SHALL take precedence.

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
  interact with Service Brokers to provision, and manage, Service Instances
  and Service Bindings.

- *Service Instance*: An instantiation of a Service offering.

- *Service Binding*: The representation of an association between an
  Application and a Service Instance. Often, Service Bindings, will
  contain the credentials that the Application will use to communicate
  with the Service Instance.

- *Application*: The software that uses a Service Instance via Service Binding.

## Context Object

In the [Open Service Broker API specification](spec.md) there are certain
message flows that include a `context` property. This property is defined
as an opaque JSON object that is meant to contain contextual information
about the environment in which the Platform or Application is executing.
For example, it might include the organizational information (eg. a
Cloud Foundry `organization` GUID) in which the Application is owned.

While the `context` property is defined as an opaque JSON object, in practice,
it is often useful and necessary for there to be an agreed upon set of
properties to ensure a common understanding of this data between the
Platform and the Service Brokers.

While use of this profile is OPTIONAL, an implementation is not compliant
with this profile if it fails to satisfy one or more of the MUST, SHALL
or REQUIRED level requirements defined herein.

### Context Object Properties

The list of properties within the Context Object can vary depending on
which Service Broker API is being invoked and which Platform is being used.
This section will define a set of properties for each platform and specify
when each is meant to be used. Platforms MAY choose to provide additional
properties beyond the ones defined in this document.

Aside from the Platform specific properties, defined in the following
sections, there is one common property called `platform` that
MUST also appear within `context` to indicate which Platform is being used.

The `platform` property MUST be a `string` and serialized as follows:
```
"platform": "platform-string-here"
```

Each section below will define the `platform` value that MUST be used based
on the Platform and the set of additional properties that MUST be present.

#### Cloud Foundry

*`platform` Property Value*: `cloudfoundry`

The following properties are defined for usage within a Cloud Foundry
deployment:

- `organization_guid`<br>
  The organization GUID as defined by the Cloud Foundry
  specification/project. This property MUST be a non-empty string serialized
  as follows:
  ```
  "organization_guid": "organization-guid-here"
  ```
  For example:
  ```
  "organization_guid": "1113aa0-124e-4af2-1526-6bfacf61b111"
  ```

- `space_guid`<br>
  The space GUID as defined by the Cloud Foundry
  specification/project. This property MUST be a non-empty string serialized
  as follows:
  ```
  "space_guid": "space-guid-here"
  ```
  For example:
  ```
  "space_guid": "aaaa1234-da91-4f12-8ffa-b51d0336aaaa"
  ```

The following table specifies which properties MUST appear in each API:

| Request API | Properties |
| --- | --- |
| `PUT /v2/service_instances/:instance_id` | `organization_guid`<br>`space_guid` |
| `PATCH /v2/service_instances/:instance_id` | `organization_guid`<br>`space_guid` |
| `PUT /v2/service_instances/:instance_id/service_bindings/:binding_id` | `organization_guid`<br>`space_guid` |

Example:

The following example shows a `context` property that might appear as
part of a Cloud Foundry API call:
  ```
  "context": {
    "platform": "cloudfoundry",
    "organization_guid": "1113aa0-124e-4af2-1526-6bfacf61b111",
    "space_guid": "aaaa1234-da91-4f12-8ffa-b51d0336aaaa"
  }
  ```

#### Kubernetes

*`platform` Property Value*: `kubernetes`

The following properties are defined for usage within a Kubernetes deployment:

- `namespace`<br>
  The name of the Kubernetes namespace in which the service instance
  will be visible. This property MUST be a non-empty string serialized
  as follows:

  ```
  "namespace": "namespace-name-here"
  ```
  For example:
  ```
  "namespace": "testing"
  ```

The following table specifies which properties MUST appear in each API:

| Request API | Properties |
| --- | --- |
| `PUT /v2/service_instances/:instance_id` | `namespace` |
| `PATCH /v2/service_instances/:instance_id` | `namespace` |
| `PUT /v2/service_instances/:instance_id/service_bindings/:binding_id` | `namespace` |

Example:

The following example shows a `context` property that might appear as
part of a Kubernetes API call:
  ```
  "context": {
    "platform": "kubernetes",
    "namespace": "development"
  }
  ```

## Service Metadata

While the [specification](spec.md) does not mandate the property names used
in the `metadata` objects, it is RECOMMENDED that the following names
be used when possible in an attempt to provide some degree of interoperability
and consistency.

### Service Metadata Fields

| Broker API Field | Type | Description |
| --- | --- | --- |
| metadata.displayName | string | The name of the service to be displayed in graphical clients. |
| metadata.imageUrl | string | The URL to an image. |
| metadata.longDescription | string | Long description. |
| metadata.providerDisplayName | string | The name of the upstream entity providing the actual service. |
| metadata.documentationUrl | string | Link to documentation page for the service. |
| metadata.supportUrl | string | Link to support page for the service. |

### Plan Metadata Fields
| Broker API Field | Type | Description |
| --- | --- | --- |
| metadata.bullets | array-of-strings | Features of this plan, to be displayed in a bulleted-list. |
| metadata.costs | object | An array-of-objects that describes the costs of a service, in what currency, and the unit of measure. If there are multiple costs, all of them could be billed to the user (such as a monthly + usage costs at once). |
| metadata.displayName | string | Name of the plan to be displayed to clients. |

#### Cost Object
This object describes the costs of a service, in what currency, and the unit
of measure.

| Field | Type | Description |
| --- | --- | --- |
| amount* | object | An array of pricing in various currencies for the cost type as key-value pairs where key is currency code and value (as a float) is currency amount. |
| unit* | string | Display name for type of cost, e.g. Monthly, Hourly, Reques, GB. |

\* Fields with an asterisk are REQUIRED.

For example:
```
"costs": [
  {
    "amount": {
      "usd": 649.0
    },
    "unit": "MONTHLY"
  }
]
```

### Example Broker Response Body

The example below contains a catalog of one service, having one service plan.
Of course, a broker can offering a catalog of many services, each having
many plans.

```
{
   "services":[
      {
      "id":"766fa866-a950-4b12-adff-c11fa4cf8fdc",
         "name":"cloudamqp",
         "description":"Managed HA RabbitMQ servers in the cloud",
         "requires":[

         ],
         "tags":[
            "amqp",
            "rabbitmq",
            "messaging"
         ],
         "metadata":{
            "displayName":"CloudAMQP",
            "imageUrl":"https://d33na3ni6eqf5j.cloudfront.net/app_resources/18492/thumbs_112/img9069612145282015279.png",
            "longDescription":"Managed, highly available, RabbitMQ clusters in the cloud",
            "providerDisplayName":"84codes AB",
            "documentationUrl":"http://docs.cloudfoundry.com/docs/dotcom/marketplace/services/cloudamqp.html",
            "supportUrl":"http://www.cloudamqp.com/support.html"
         },
         "dashboard_client":{
            "id": "p-mysql-client",
            "secret": "p-mysql-secret",
            "redirect_uri": "http://p-mysql.example.com/auth/create"
         },
         "plans":[
            {
               "id":"024f3452-67f8-40bc-a724-a20c4ea24b1c",
               "name":"bunny",
               "description":"A mid-sided plan",
               "metadata":{
                  "bullets":[
                     "20 GB of messages",
                     "20 connections"
                  ],
                  "costs":[
                     {
                        "amount":{
                           "usd":99.0
                        },
                        "unit":"MONTHLY"
                     },
                     {
                        "amount":{
                           "usd":0.99
                        },
                        "unit":"1GB of messages over 20GB"
                     }
                  ],
                  "displayName":"Big Bunny"
               }
            }
         ]
      }
   ]
}
```
