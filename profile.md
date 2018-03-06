# Open Service Broker API Profile (master - might contain changes that are not yet released)

## Abstract

The [Open Service Broker API specification](spec.md) allows for extensions
and variations based on the environments in which it is being used; this
document contains the suggested usage pattern for some of those variants.

While use of this profile is OPTIONAL, an implementation is not compliant
with this profile if it fails to satisfy one or more of the MUST, SHALL
or REQUIRED level requirements defined herein.

## Table of Contents

- [Notations and Terminology](#notations-and-terminology)
  - [Notational Conventions](#notational-conventions)
  - [Terminology](#terminology)
- [Originating Identity Header](#originating-identity-header)
  - [Cloud Foundry Originating Identity Header](#cloud-foundry-originating-identity-header)
  - [Kubernetes Originating Identity Header](#kubernetes-originating-identity-header)
- [Context Object](#context-object)
  - [Context Object Properties](#context-object-properties)
  - [Cloud Foundry Context Object](#cloud-foundry-context-object)
  - [Kubernetes Context Object](#kubernetes-context-object)
- [Bind Resource Object](#bind-resource-object)
  - [Cloud Foundry Bind Resource Object](#cloud-foundry-bind-resource-object)
- [Service Metadata](#service-metadata)
  - [Cloud Foundry Service Metadata](#cloud-foundry-service-metadata)

## Notations and Terminology

### Notational Conventions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED",  "MAY", and "OPTIONAL" in this document are to
be interpreted as described in [RFC 2119]( https://tools.ietf.org/html/rfc2119).

### Terminology

Please refer to terminology defined by the
[Open Service Broker API specification](spec.md#terminology).

## Originating Identity Header

In the [Open Service Broker API specification](spec.md) it defines an
additional HTTP Header that can be included in messages from the Platform
to identify the user that requested the action to be taken.

The header consists of two parts: a `platform` string and `value` string,
where the `value` is a Base64 encoded serialized JSON object.
Both parts will vary based on the Platform which is being used. The
following sections define the values to be used based on the Platform
and which properties are expected to appear in the `value` JSON.

Note that when both the originating identity HTTP Header and the Context
object appear in the same message the `platform` value MUST be the same
for both.

### Cloud Foundry Originating Identity Header

*`platform` Value*: `cloudfoundry`

The following properties MUST appear within the JSON encoded `value`:

| Property | Type | Description |
| --- | --- | --- |
| user_id | string | The `user_id` value from the Cloud Foundry JWT token. |

Platforms MAY include additional properties.

For example, a `value` of:
```json
{
  "user_id": "683ea748-3092-4ff4-b656-39cacc4d5360"
}
```
would appear in the HTTP Header as:
```
X-Broker-API-Originating-Identity: cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIiwNCiAgInVzZXJfbmFtZSI6ICJqb2VAZXhhbXBsZS5jb20iDQp9
```

### Kubernetes Originating Identity Header

*`platform` Value*: `kubernetes`

The following properties MUST appear within the JSON encoded `value`:

| Property | Type | Description |
| --- | --- | --- |
| username | string | The `username` property from the Kubenernetes `user.info` object. |
| uid | string | The `uid` property from the Kubenernetes `user.info` object. |
| groups | array-of-strings | The `groups` property from the Kubenernetes `user.info` object. |
| extra | map-of-array-of-strings | The `extra` property from the Kubernetes `user.info` object. |

Platforms MAY include additional properties.

For example, a `value` of:
```json
{
  "username": "duke",
  "uid": "c2dde242-5ce4-11e7-988c-000c2946f14f",
  "groups": [ "admin", "dev" ],
  "extra": {
    "mydata": [ "data1", "data3" ]
  }
}
```
would appear in the HTTP Header as:
```
X-Broker-API-Originating-Identity: kubernetes eyANCiAgInVzZXJuYW1lIjogImR1a2UiLA0KICAidWlkIjogImMyZGRlMjQyLTVjZTQtMTFlNy05ODhjLTAwMGMyOTQ2ZjE0ZiIsDQogICJncm91cHMiOiB7ICJhZG1pbiIsICJkZXYiIH0NCn0=
```


## Context Object

In the [Open Service Broker API specification](spec.md) there are certain
message flows that include a `context` property. This property is defined
as an opaque JSON object that is meant to contain contextual information
about the environment in which the Platform or Application is executing.

While the `context` property is defined as an opaque JSON object, in practice,
it is often useful and necessary for there to be an agreed upon set of
properties to ensure a common understanding of this data between the
Platform and the Service Brokers.

### Context Object Properties

The list of properties within the Context Object can vary depending on
which Service Broker API is being invoked and which Platform is being used.
This section will define a set of properties for each Platform and specify
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

Note that when both the originating identity HTTP Header and the Context
object appear in the same message the `platform` value MUST be the same
for both.

### Cloud Foundry Context Object

*`platform` Property Value*: `cloudfoundry`

The following properties are defined for usage within a Cloud Foundry
deployment:

- `organization_guid`  
  The GUID of the organization that a Service Instance is associated with.
  This property MUST be a non-empty string serialized as follows:
  ```
  "organization_guid": "organization-guid-here"
  ```
  For example:
  ```
  "organization_guid": "1113aa0-124e-4af2-1526-6bfacf61b111"
  ```

- `space_guid`  
  The GUID of the space that a Service Instance is associated with.
  This property MUST be a non-empty string serialized as follows:
  ```
  "space_guid": "space-guid-here"
  ```
  For example:
  ```
  "space_guid": "aaaa1234-da91-4f12-8ffa-b51d0336aaaa"
  ```

The following table specifies the REQUIRED properties for requests from a
Platform to a Service Broker:

| Request API | Properties |
| --- | --- | --- |
| `PUT /v2/service_instances/:instance_id` | `organization_guid`, `space_guid` |
| `PATCH /v2/service_instances/:instance_id` | `organization_guid`, `space_guid` |
| `PUT /v2/service_instances/:instance_id/service_bindings/:binding_id` | `organization_guid`, `space_guid` |

The following example shows a `context` object that might appear as part of a
Cloud Foundry API call:
  ```
  "context": {
    "platform": "cloudfoundry",
    "organization_guid": "1113aa0-124e-4af2-1526-6bfacf61b111",
    "space_guid": "aaaa1234-da91-4f12-8ffa-b51d0336aaaa"
  }
  ```

### Kubernetes Context Object

*`platform` Property Value*: `kubernetes`

The following properties are defined for usage within a Kubernetes deployment:

- `namespace`<br>
  The name of the Kubernetes namespace in which the Service Instance
  will be visible. This property MUST be a non-empty string serialized
  as follows:

  ```
  "namespace": "namespace-name-here"
  ```
  For example:
  ```
  "namespace": "testing"
  ```

- `clusterID`<br>
  The unique identifier for the Kubernetes cluster from which the request
  was sent. This property MUST be a non-empty string serialized as follows:

  ```
  "clusterID": "id-goes-here"
  ```
  For example:
  ```
  "clusterID": "644e1dd7-2a7f-18fb-b8ed-ed78c3f92c2b"
  ```

The following table specifies which properties MUST appear in each API:

| Request API | Properties |
| --- | --- |
| `PUT /v2/service_instances/:instance_id` | `namespace`, `clusterID` |
| `PATCH /v2/service_instances/:instance_id` | `namespace`, `clusterID` |
| `PUT /v2/service_instances/:instance_id/service_bindings/:binding_id` | `namespace`, `clusterID` |

Example:

The following example shows a `context` property that might appear as
part of a Kubernetes API call:
  ```
  "context": {
    "platform": "kubernetes",
    "namespace": "development",
    "clusterID": "8263feba-9b8a-23ae-99ed-abcd1234feda"
  }
  ```

## Bind Resource Object

In the [Open Service Broker API specification](spec.md), requests to
[create a Service Binding](spec.md#binding) can contain a `bind_resource`
object in which Platforms MAY choose to add additional fields.

### Cloud Foundry Bind Resource Object

The following properties are defined for usage within a Cloud Foundry
deployment:

- `space_guid`  
  This OPTIONAL property is the GUID of a space that a Service Binding is
  associated with. If present, this property MUST be a non-empty string
  serialized as follows:
  ```
  "space_guid": "space-guid-here"
  ```
  For example:
  ```
  "space_guid": "15823972-c216-4ba5-9f3f-e75b0b891297"
  ```

The following example shows a `bind_resource` object that might appear as part
of a Cloud Foundry API call:
  ```
  "bind_resource": {
    "app_guid": "5e76c9bf-d5e3-46bf-9877-6d8dddfc8a45",
    "space_guid": "15823972-c216-4ba5-9f3f-e75b0b891297"
  }
  ```

## Service Metadata

While the [specification](spec.md) does not mandate the property names used
in the `metadata` objects, it is RECOMMENDED that the following names
be used when possible in an attempt to provide some degree of interoperability
and consistency.

#### Service Metadata Fields

| Service Broker API Field | Type | Description |
| --- | --- | --- |
| metadata.displayName | string | The name of the service to be displayed in graphical clients. |
| metadata.imageUrl | string | The URL to an image. |
| metadata.longDescription | string | Long description. |
| metadata.providerDisplayName | string | The name of the upstream entity providing the actual service. |
| metadata.documentationUrl | string | Link to documentation page for the service. |
| metadata.supportUrl | string | Link to support page for the service. |

#### Plan Metadata Fields
| Service Broker API Field | Type | Description |
| --- | --- | --- |
| metadata.bullets | array-of-strings | Features of this plan, to be displayed in a bulleted-list. |
| metadata.costs | array-of-objects | An array-of-objects that describes the costs of a service, in what currency, and the unit of measure. If there are multiple costs, all of them could be billed to the user (such as a monthly + usage costs at once). |
| metadata.displayName | string | Name of the plan to be displayed to clients. |

#### Cost Object
This object describes the costs of a service, in what currency, and the unit
of measure.

| Field | Type | Description |
| --- | --- | --- |
| amount* | object | An array of pricing in various currencies for the cost type as key-value pairs where key is currency code and value (as a float) is currency amount. |
| unit* | string | Display name for type of cost, e.g. Monthly, Hourly, Request, GB. |

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

#### Example Service Broker Response Body

The example below contains a catalog of one service, having one Service Plan.
Of course, a Service Broker can offering a catalog of many services, each having
many plans.

```json
{
  "services":[
    {
      "id":"766fa866-a950-4b12-adff-c11fa4cf8fdc",
      "name":"cloudamqp",
      "description":"Managed HA RabbitMQ servers in the cloud.",
      "requires":[],
      "tags":[
        "amqp",
        "rabbitmq",
        "messaging"
      ],
      "metadata":{
        "displayName":"CloudAMQP",
        "imageUrl":"https://d33na3ni6eqf5j.cloudfront.net/app_resources/18492/thumbs_112/img9069612145282015279.png",
        "longDescription":"Managed, highly available, RabbitMQ clusters in the cloud.",
        "providerDisplayName":"84codes AB",
        "documentationUrl":"http://docs.cloudfoundry.com/docs/dotcom/marketplace/services/cloudamqp.html",
        "supportUrl":"http://www.cloudamqp.com/support.html"
      },
      "dashboard_client":{
        "id":"p-mysql-client",
        "secret":"p-mysql-secret",
        "redirect_uri":"http://p-mysql.example.com/auth/create"
      },
      "plans":[
        {
          "id":"024f3452-67f8-40bc-a724-a20c4ea24b1c",
          "name":"bunny",
          "description":"A mid-sided plan.",
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

### Cloud Foundry Service Metadata

In addition to the metadata described in [Service Metadata](#service-metadata),
Service Brokers MAY also expose the following fields to enable Cloud Foundry
specific behaviour.

#### Service Metadata Fields

| Broker API Field | Type | Description |
| --- | --- | --- |
| metadata.shareable | string | Allows Service Instances to be shared across orgs and spaces. |
