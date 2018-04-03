# Open Service Broker API (master - might contain changes that are not yet released)

## Table of Contents
  - [API Overview](#api-overview)
  - [Notations and Terminology](#notations-and-terminology)
  - [Changes](#changes)
    - [Change Policy](#change-policy)
    - [Changes Since v2.12](#changes-since-v212)
  - [API Version Header](#api-version-header)
  - [Platform to Service Broker Authentication](#platform-to-service-broker-authentication)
  - [URL Properties](#url-properties)
  - [Originating Identity](#originating-identity)
  - [Service Broker Errors](#service-broker-errors)
  - [Catalog Management](#catalog-management)
    - [Adding a Service Broker to the Platform](#adding-a-service-broker-to-the-platform)
  - [Service Dependencies](#service-dependencies)
  - [Synchronous and Asynchronous Operations](#synchronous-and-asynchronous-operations)
    - [Synchronous Operations](#synchronous-operations)
    - [Asynchronous Operations](#asynchronous-operations)
  - [Blocking Operations](#blocking-operations)
  - [Polling Last Operation](#polling-last-operation)
    - [Polling Interval and Duration](#polling-interval-and-duration)
  - [Provisioning](#provisioning)
  - [Updating a Service Instance](#updating-a-service-instance)
  - [Binding](#binding)
    - [Types of Binding](#types-of-binding)
  - [Unbinding](#unbinding)
  - [Deprovisioning](#deprovisioning)
  - [Orphans](#orphans)

## API Overview

The Open Service Broker API defines an HTTP(S) interface between Platforms and
Service Brokers.

The Service Broker is the component of the service that implements the Service
Broker API, for which a Platform is a client. Service Brokers are responsible
for advertising a catalog of Service Offerings and Service Plans to the
Platform, and acting on requests from the Platform for provisioning, binding,
unbinding, and deprovisioning.

In general, provisioning reserves a resource on a service; we call this
reserved resource a Service Instance. What a Service Instance represents can
vary by service. Examples include a single database on a multi-tenant server,
a dedicated cluster, or an account on a web application.

What a binding represents MAY also vary by service. In general creation of a
binding either generates credentials necessary for accessing the resource or
provides the Service Instance with information for a configuration change.

A Platform MAY expose services from one or many Service Brokers, and an
individual Service Broker MAY support one or many Platforms using different URL
prefixes and credentials.

## Notations and Terminology

### Notational Conventions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to
be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119).

### Terminology

This specification defines the following terms:

- *Application*: Often the entity using a Service Instance will be a piece of
  software, however, this does not need to be the case. For the purposes of
  this specification, the term "Application" will be used to represent all
  entities that might make use of, and be bound to, a Service Instance.

- *Platform*: The software that will manage the cloud environment into which
  Applications are provisioned and Service Brokers are registered. Users will
  not directly provision Services from Service Brokers, rather they will ask
  the Platform to manage Services and interact with the Service Brokers for
  them.

- *Service*: A managed software offering that can be used by an Application.
  Typically, Services will expose some API that can be invoked to perform
  some action. However, there can also be non-interactive Services that can
  perform the desired actions without direct prompting from the Application.

- *Service Binding*: Represents the request to use a Service Instance. As part
  of this request there might be a reference to the entity, also known as the
  Application, that will use the Service Instance. Service Bindings will often
  contain the credentials that can then be used to communicate with the Service
  Instance.

- *Service Broker*: Service Brokers manage the lifecycle of Services. Platforms
  interact with Service Brokers to provision, and manage, Service Instances
  and Service Bindings.

- *Service Instance*: An instantiation of a Service Offering.

- *Service Offering*: The advertisement of a service that a Service Broker
  supports.

- *Service Plan*: The representation of the costs and benefits for a given
  variant of the service, potentially as a tier.


## Changes

### Change Policy

* Existing endpoints and fields MUST NOT be removed or renamed.
* New OPTIONAL endpoints, or new HTTP methods for existing endpoints, MAY be
added to enable support for new features.
* New fields MAY be added to existing request/response messages.
These fields MUST be OPTIONAL and SHOULD be ignored by clients and servers
that do not understand them.

### Changes Since v2.12

* Added [`schemas`](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#schema-object)
  field to services in the catalog that Service Brokers can use to declare the
  configuration parameters their service accepts for creating a Service
  Instance, updating a Service Instance and creating a Service Binding.
* Added [`context`](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#binding)
  field to request body for creating a Service Binding.
* Added [warning text](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#url-properties)
  about using characters outside of the "Unreserved Characters" set in IDs.
* Added information about
  [`volume_mounts`](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#volume-mounts-object)
  objects.
* `instance_id` and `binding_id` MUST be globally unique non-empty strings.
* Allow [non-BasicAuth authentication mechanisms](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#authentication).
* Added a [Getting Started](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/gettingStarted.md)
  page including sample Service Brokers.
* Define what a [CLI-friendly string](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#catalog-management)
  is.
* Add [service/plan metadata conventions](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/profile.md#service-metadata).
* Add [originating identity HTTP header](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#originating-identity).

For changes in older versions, see the [release notes](https://github.com/openservicebrokerapi/servicebroker/blob/master/release-notes.md).

## API Version Header

Requests from the Platform to the Service Broker MUST contain a header that
declares the version number of the  Open Service Broker API that the Platform
will use:

`X-Broker-API-Version: 2.13`

The version numbers are in the format `MAJOR.MINOR` using semantic versioning.

This header allows Service Brokers to reject requests from Platforms for
versions they do not support. While minor API revisions will always be
additive, it is possible that Service Brokers depend on a feature from a newer
version of the API that is supported by the Platform. In this scenario the
Service Broker MAY reject the request with `412 Precondition Failed` and
provide a message that informs the operator of the API version that is to be
used instead.

## Vendor Extension Fields

Senders of messages defined by this specification MAY include additional
fields within the JSON objects. When adding new fields, unique prefixes
SHOULD be used for the field names to reduce the chances of conflicts with
with future specification defined fields or other extensions.

Receivers of a messages defined by this specification that contain unknown
extension fields MUST ignore those fields and MUST NOT halt processing
of those messages due to the presence of those fields. Receiver are under
no obligation to understand or process unknown extension fields.

## Platform to Service Broker Authentication

While the communication between a Platform and Service Broker MAY be unsecure,
it is RECOMMENDED that all communications between a Platform and a Service
Broker are secured via TLS and authenticated.

Unless there is some out of band communication and agreement between a
Platform and a Service Broker, the Platform MUST authenticate with the
Service Broker using HTTP basic authentication (the `Authorization:` header)
on every request. This specification does not specify how Platform and Service
Brokers agree on other methods of authentication.

Platforms and Service Brokers MAY agree on an authentication mechanism other
than basic authentication, but the specific agreements are not covered by this
specification. Please see the 
[Platform Features authentication mechanisms wiki document](https://github.com/openservicebrokerapi/servicebroker/wiki/Platform-Features)
for details on these mechanisms.

If authentication is used, the Service Broker MUST authenticate the request
using the predetermined authentication mechanism, and MUST return a `401 Unauthorized`
response if the authentication fails.

Additionally, the Service Broker MUST secure communications with TLS. The Platform
and Service Broker SHOULD agree whether the Service Broker will use a root-signed
certificate or a self-signed certificate.

Note: Using an authentication mechanism that is agreed to via out of band
communications could lead to interoperability issues with other Platforms.

## URL Properties

This specification defines the following properties that might appear within
URLs:
- `service_id`
- `plan_id`
- `instance_id`
- `binding_id`
- `operation`

While this specification places no restriction on the set of characters used
within these strings, it is RECOMMENDED that these properties only contain
characters from the "Unreserved Characters" as defined by
[RFC3986](https://tools.ietf.org/html/rfc3986#section-2.3). In other words:
uppercase and lowercase letters, decimal digits, hyphen, period, underscore
and tilde.

If a character outside of the "Unreserved Characters" set is used, then it
SHOULD be percent-encoded prior to being used as part of the HTTP request, per
[RFC3986](https://tools.ietf.org/html/rfc3986#section-2.1).

## Originating Identity

Often a Service Broker will need to know the identity of the user that
initiated the request from the Platform. For example, this might be needed for
auditing or authorization purposes. In order to facilitate this, the Platform
will need to provide this identification information to the Service Broker on
each request. Platforms MAY support this feature, and if they do, they MUST
adhere to the following:
- For any OSBAPI request that is the result of an action taken by a Platform's
  user, there MUST be an associated `OriginatingIdentity` header on that HTTP
  request.
- Any OSBAPI request that is not associated with an action from a Platform's
  user, such as the Platform refetching the catalog, MAY exclude the header from
  that HTTP request.
- If present on a request, the `OriginatingIdentity` header MUST contain the
  identify information for the Platform's user that took the action to cause the
  request to be sent.

The format of the header MUST be:

```
X-Broker-API-Originating-Identity: Platform value
```

`Platform` MUST be a non-empty string indicating the Platform from which
the request is being sent. The specific value SHOULD match the values
defined in the [profile](profile.md) document for the `context.platform`
property. When `context` is sent as part of a message, this value MUST
be the same as the `context.platform` value.

`value` MUST be a Base64 encoded string. The string MUST be a serialized
JSON object. The specific properties will be Platform specific - see
the [profile](profile.md) document for more information.

For example:
```
X-Broker-API-Originating-Identity: cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIiwNCiAgInVzZXJfbmFtZSI6ICJqb2VAZXhhbXBsZS5jb20iDQp9
```

Where the `value`, when decoded, is:
```
{
  "user_id": "683ea748-3092-4ff4-b656-39cacc4d5360"
}
```

Note that not all messages sent to a Service Broker are initiated by an
end-user of the Platform. For example, during orphan mitigation or during the
querying of the Service Broker's catalog, the Platform might not have an
end-user with which to associate the request, therefore in those cases the
originating identity header would not be included in those messages.

## Service Broker Errors

When a request to a Service Broker fails, the Service Broker MUST return an
appropriate HTTP response code. Where the specification defines the expected
response code, that response code MUST be used.

For error responses, the following fields are defined:

| Response Field | Type | Description |
| --- | --- | --- |
| error | string | A single word in camel case that uniquely identifies the error condition. If present, MUST be a non-empty string. |
| description | string | A user-facing error message explaining why the request failed. If present, MUST be a non-empty string. |

### Error Codes

There are failure scenarios described throughout this specification for which
the `error` field MUST contain a specific string. Service Broker authors MUST
use these error codes for the specified failure scenarios.

| Error | Reason | Expected Action |
| --- | --- | --- |
| AsyncRequired | This request requires client support for asynchronous service operations. | The query parameter `accepts_incomplete=true` MUST be included the request. |
| ConcurrencyError | The Service Broker does not support concurrent requests that mutate the same resource. | Clients MUST wait until pending requests have completed for the specified resources. |
| RequiresApp | The request body is missing the `app_guid` field. | The `app_guid` MUST be included in the request body. |

## Catalog Management

The first endpoint that a Platform will interact with on the Service Broker is
the service catalog (`/v2/catalog`). This endpoint returns a list of all
services available on the Service Broker. Platforms query this endpoint from
all Service Brokers in order to present an aggregated user-facing catalog.

Periodically, a Platform MAY re-query the service catalog endpoint for a
Service Broker to see if there are any changes to the list of services.
Service Brokers MAY add, remove or modify (metadata, plans, etc.) the list of
services from previous queries.

When determining what, if anything, has changed on a Service Broker, the
Platform will use the `id` of the resources (services or plans) as the only
immutable property and MUST use that to locate the same resource as was
returned from a previous query. Likewise, a Service Broker MUST NOT change the
`id` of a resource across queries, otherwise a Platform will treat it as a
different resource.

When a Platform receives different `id` values for the same type of resource,
even if all of the other metadata in those resources are the exact same, it
MUST treat them as separate instances of that resource.

Service Broker authors are expected to be cautious when removing services and
plans from their catalogs, as Platforms might have provisioned Service
Instances of these plans. For example, Platforms might restrict the actions
that users can perform on existing Service Instances if the associated service
or plan is deleted. Consider your deprecation strategy.


The following sections describe catalog requests and responses in the Service
Broker API.

### Request

#### Route
`GET /v2/catalog`

#### Headers

The following HTTP Headers are defined for this operation:

| Header | Type | Description |
| --- | --- | --- |
| X-Broker-API-Version* | string | See [API Version Header](#api-version-header). |

\* Headers with an asterisk are REQUIRED.

#### cURL
```
$ curl http://username:password@service-broker-url/v2/catalog -H "X-Broker-API-Version: 2.13"
```

### Response

| Status Code | Description |
| --- | --- |
| 200 OK | MUST be returned upon successful processing of this request. The expected response body is below. |

#### Body

CLI clients will typically have restrictions on how names, such as service
and plan names, will be provided by users. Therefore, this specification
defines a "CLI-friendly" string as a short string that MUST only use
alphanumeric characters, periods, and hyphens, with no spaces. This will make it
easier for users when they have to type it as an argument on the command line.
For comparison purposes, service and plan names MUST be treated as
case-sensitive strings.

Note: In previous versions of the specification service and plan names
were not allowed to use uppercase characters. However, this requirement was
not enforced and therefore to ensure backwards compatibility with existing
Service Brokers that might use uppercase characters the specification
has been changed.

For backwards compatibility reasons, this specification does not preclude
the use of CLI-unfriendly strings that might cause problems for command line
parsers (or that are not very meaningful to users), such as `-`.
It is therefore RECOMMENDED that implementations avoid such strings.

| Response Field | Type | Description |
| --- | --- | --- |
| services* | array of [Service](#service-object) objects | Schema of service objects defined below. MAY be empty. |

\* Fields with an asterisk are REQUIRED.

##### Service Object

| Response Field | Type | Description |
| --- | --- | --- |
| name* | string | A CLI-friendly name of the service. MUST only contain alphanumeric characters, periods, and hyphens (no spaces). MUST be unique across all service objects returned in this response. MUST be a non-empty string. |
| id* | string | An identifier used to correlate this service in future requests to the Service Broker. This MUST be globally unique such that Platforms (and their users) MUST be able to assume that seeing the same value (no matter what Service Broker uses it) will always refer to this service. MUST be a non-empty string. Using a GUID is RECOMMENDED. |
| description* | string | A short description of the service. MUST be a non-empty string. |
| tags | array of strings | Tags provide a flexible mechanism to expose a classification, attribute, or base technology of a service, enabling equivalent services to be swapped out without changes to dependent logic in applications, buildpacks, or other services. E.g. mysql, relational, redis, key-value, caching, messaging, amqp. |
| requires | array of strings | A list of permissions that the user would have to give the service, if they provision it. The only permissions currently supported are `syslog_drain`, `route_forwarding` and `volume_mount`. |
| bindable* | boolean | Specifies whether Service Instances of the service can be bound to applications. This specifies the default for all plans of this service. Plans can override this field (see [Plan Object](#plan-object)). |
| metadata | object | An opaque object of metadata for a Service Offering. It is expected that Platforms will treat this as a blob. Note that there are [conventions](profile.md#service-metadata) in existing Service Brokers and Platforms for fields that aid in the display of catalog data. |
| dashboard_client | [DashboardClient](#dashboard-client-object) | Contains the data necessary to activate the Dashboard SSO feature for this service. |
| plan_updateable | boolean | Whether the service supports upgrade/downgrade for some plans. Please note that the misspelling of the attribute `plan_updatable` as `plan_updateable` was done by mistake. We have opted to keep that misspelling instead of fixing it and thus breaking backward compatibility. Defaults to false. |
| plans* | array of [Plan](#plan-object) objects | A list of plans for this service, schema is defined below. MUST contain at least one plan. |

\* Fields with an asterisk are REQUIRED.

Note: Platforms will typically use the service name as an input parameter
from their users to indicate which service they want to instantiate. Therefore,
it is important that these values be unique for all services that have been
registered with a Platform. To achieve this goal service providers often will
prefix their service names with some unique value (such as the name of their
company). Additionally, some Platforms might modify the service names before
presenting them to their users. This specification places no requirements on
how Platforms might expose these values to their users.

##### Dashboard Client Object

| Response Field | Type | Description |
| --- | --- | --- |
| id | string | The id of the Oauth client that the dashboard will use. If present, MUST be a non-empty string. |
| secret | string | A secret for the dashboard client. If present, MUST be a non-empty string. |
| redirect_uri | string | A URI for the service dashboard. Validated by the OAuth token server when the dashboard requests a token. |


##### Plan Object

| Response Field | Type | Description |
| --- | --- | --- |
| id* | string | An identifier used to correlate this plan in future requests to the Service Broker. This MUST be globally unique such that Platforms (and their users) MUST be able to assume that seeing the same value (no matter what Service Broker uses it) will always refer to this plan and for the same service. MUST be a non-empty string. Using a GUID is RECOMMENDED. |
| name* | string | The CLI-friendly name of the plan. MUST only contain alphanumeric characters, periods, and hyphens (no spaces). MUST be unique within the service. MUST be a non-empty string. |
| description* | string | A short description of the plan. MUST be a non-empty string. |
| metadata | object | An opaque object of metadata for a Service Plan. It is expected that Platforms will treat this as a blob. Note that there are [conventions](profile.md#service-metadata) in existing Service Brokers and Platforms for fields that aid in the display of catalog data. |
| free | boolean | When false, Service Instances of this plan have a cost. The default is true. |
| bindable | boolean | Specifies whether Service Instances of the Service Plan can be bound to applications. This field is OPTIONAL. If specified, this takes precedence over the `bindable` attribute of the service. If not specified, the default is derived from the service. |
| schemas | [Schemas](#schemas-object) | Schema definitions for Service Instances and bindings for the plan. |
| dependencies | [Dependencies](#dependencies) | Dependencies |

\* Fields with an asterisk are REQUIRED.

##### Schemas Object

| Response Field | Type | Description |
| --- | --- | --- |
| service_instance | [ServiceInstanceSchema](#service-instance-schema-object) | The schema definitions for creating and updating a Service Instance. |
| service_binding | [ServiceBindingSchema](#service-binding-schema-object) | The schema definition for creating a Service Binding. Used only if the Service Plan is bindable. |

##### Service Instance Schema Object

| Response Field | Type | Description |
| --- | --- | --- |
| create | [InputParametersSchema](#input-parameters-schema-object) | The schema definition for creating a Service Instance. |
| update | [InputParametersSchema](#input-parameters-schema-object) | The schema definition for updating a Service Instance. |

##### Service Binding Schema Object

| Response Field | Type | Description |
| --- | --- | --- |
| create | [InputParametersSchema](#input-parameters-schema-object) | The schema definition for creating a Service Binding. |

##### Input Parameters Schema Object

| Response Field | Type | Description |
| --- | --- | --- |
| parameters | JSON schema object | The schema definition for the input parameters. Each input parameter is expressed as a property within a JSON object. |

The following rules apply if `parameters` is included anywhere in the catalog:
* Platforms MUST support at least
[JSON Schema draft v4](http://json-schema.org/).
* Platforms SHOULD be prepared to support later versions of JSON schema.
* The `$schema` key MUST be present in the schema declaring the version of JSON
schema being used.
* Schemas MUST NOT contain any external references.
* Schemas MUST NOT be larger than 64kB.

```
{
  "services": [{
    "name": "fake-service",
    "id": "acb56d7c-XXXX-XXXX-XXXX-feb140a59a66",
    "description": "A fake service.",
    "tags": ["no-sql", "relational"],
    "requires": ["route_forwarding"],
    "bindable": true,
    "metadata": {
      "provider": {
        "name": "The name"
      },
      "listing": {
        "imageUrl": "http://example.com/cat.gif",
        "blurb": "Add a blurb here",
        "longDescription": "A long time ago, in a galaxy far far away..."
      },
      "displayName": "The Fake Service Broker"
    },
    "dashboard_client": {
      "id": "398e2f8e-XXXX-XXXX-XXXX-19a71ecbcf64",
      "secret": "277cabb0-XXXX-XXXX-XXXX-7822c0a90e5d",
      "redirect_uri": "http://localhost:1234"
    },
    "plan_updateable": true,
    "plans": [{
      "name": "fake-plan-1",
      "id": "d3031751-XXXX-XXXX-XXXX-a42377d3320e",
      "description": "Shared fake Server, 5tb persistent disk, 40 max concurrent connections.",
      "free": false,
      "metadata": {
        "max_storage_tb": 5,
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
        "bullets": [
          "Shared fake server",
          "5 TB storage",
          "40 concurrent connections"
        ]
      },
      "schemas": {
        "service_instance": {
          "create": {
            "parameters": {
              "$schema": "http://json-schema.org/draft-04/schema#",
              "type": "object",
              "properties": {
                "billing-account": {
                  "description": "Billing account number used to charge use of shared fake server.",
                  "type": "string"
                }
              }
            }
          },
          "update": {
            "parameters": {
              "$schema": "http://json-schema.org/draft-04/schema#",
              "type": "object",
              "properties": {
                "billing-account": {
                  "description": "Billing account number used to charge use of shared fake server.",
                  "type": "string"
                }
              }
            }
          }
        },
        "service_binding": {
          "create": {
            "parameters": {
              "$schema": "http://json-schema.org/draft-04/schema#",
              "type": "object",
              "properties": {
                "billing-account": {
                  "description": "Billing account number used to charge use of shared fake server.",
                  "type": "string"
                }
              }
            }
          }
        }
      }
    }, {
      "name": "fake-plan-2",
      "id": "0f4008b5-XXXX-XXXX-XXXX-dace631cd648",
      "description": "Shared fake Server, 5tb persistent disk, 40 max concurrent connections. 100 async.",
      "free": false,
      "metadata": {
        "max_storage_tb": 5,
        "costs":[
            {
               "amount":{
                  "usd":199.0
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
        "bullets": [
          "40 concurrent connections"
        ]
      }
    }]
  }]
}
```


### Adding a Service Broker to the Platform

After implementing the first endpoint `GET /v2/catalog` documented
[above](#catalog-management), the Service Broker will need to be registered
with your Platform to make your services and plans available to end users.

## Service Dependencies

Broker authors can define dependencies of their services to other services.
Dependencies are declared in the catalog at the plan level.
There are three types of dependencies.

* Open dependencies are arbitrary dependencies.
  The service instance accepts any bindings and service instance IDs and
  evaluates them at runtime.
* Service dependencies are specified by a service or a set of plans.
  When the service instance is provisioned, the platform has to provide a
  binding or a service instance ID of a matching service instance.
* Instance dependencies declare dependencies to a specific instance either
  by a service instance ID or by a service instance alias.

Dependencies are injected as parameters when a service instance is provisioned.
The names of the parameters are declared in the catalog.

The parameter values are objects and are different for bindable and
non-bindable services. They all have the following fields in common.

| Field | Type | Description |
| --- | --- | --- |
| service_id* | string | Service ID of the instance. |
| service_name* | string | Service name of the instance. |
| plan_id* | string | Plan ID of the instance. |
| plan_name* | string | Plan name of the instance. |

\* Fields with an asterisk are REQUIRED.

For bindable services, the value MUST contain be the response of a binding request. 

```
{
  "parameters": {
    "database": {
      "service_id": "71ef3933-8520-4f26-b808-f483d23f94a6",
      "service_name": "mysql",
      "plan_id": "1b367d58-6f17-4cb3-9ef8-3a69e752b826",
      "plan_name": "small",
      "credentials": {
        "uri": "mysql://mysqluser:pass@mysqlhost:3306/dbname",
        "username": "mysqluser",
        "password": "pass",
        "host": "mysqlhost",
        "port": 3306,
        "database": "dbname"
      }
    }
  }
}
```

Instances of non-bindable services can only be provided to the same broker.
The platform MUST set the instance ID into the parameter value object.

| Field | Type | Description |
| --- | --- | --- |
| instance_id* | string | Service instance ID. |

\* Fields with an asterisk are REQUIRED.

```
{
  "parameters": {
    "database": {
      "service_id": "90c7c453-c6a3-47ef-8284-5a2c00f4f18b",
      "service_name": "the-database",
      "plan_id": "d9301d66-7621-4526-9b97-63e93ddc2d13",
      "plan_name": "big",
      "instance_id": "0c7ec11f-977b-4027-a4d2-3dd81a9fbeb6"
    }
  }
```

The value of an open dependencies parameter is an array of objects, following
the structures above. Additionally, the platform MAY add a dependency name,
set by the user to each object.

| Field | Type | Description |
| --- | --- | --- |
| name* | string | Dependency name. |

\* Fields with an asterisk are REQUIRED.

```
{
  "parameters": {
    "some_dependencies ": [
      {
        "name": "master-database",
        "service_id": "71ef3933-8520-4f26-b808-f483d23f94a6",
        "service_name": "mysql",
        "plan_id": "1b367d58-6f17-4cb3-9ef8-3a69e752b826",
        "plan_name": "small",
        "credentials": {
          "uri": "mysql://mysqluser:pass@mysqlhost:3306/dbname",
          "username": "mysqluser",
          "password": "pass",
          "host": "mysqlhost",
          "port": 3306,
          "database": "dbname"
        }
      },
      {
        "name": "internal-store",
        "service_id": "90c7c453-c6a3-47ef-8284-5a2c00f4f18b",
        "service_name": "the-database",
        "plan_id": "d9301d66-7621-4526-9b97-63e93ddc2d13",
        "plan_name": "big",
        "instance_id": "0c7ec11f-977b-4027-a4d2-3dd81a9fbeb6"
      }
    ]
  }
}
```

### Dependencies

| Field | Type | Description |
| --- | --- | --- |
| open\_dependencies\_parameter | string | A parameter name used for an open list of dependencies. |
| services | [array-of-service-dependency-objects](#service-dependencies) | A list of service dependencies. |
| instances | [array-of-instance-dependency-objects](#instance-dependencies) | A list of instance dependencies. |

\* Fields with an asterisk are REQUIRED.

### Service Dependencies

| Field | Type | Description |
| --- | --- | --- |
| parameter* | string | Name of the parameter. MUST be a non-empty string and MUST be unique within all dependency parameters. |
| description* | string | A short description of the dependency. MUST be a non-empty string. |
| optional  | boolean | Whether the dependency is optional or not. Defaults to false. |
| service_id | string | Service ID of the instance. MUST be a non-empty string. MUST NOT be used in conjunction with the fields `plan_ids` and `service_name`. |
| plan_ids | array-of-strings | An array of plan IDs. The instance must match one of those plans. MUST be a non-empty array of non-empty strings. MUST NOT be used in conjunction with the fields `service_id` and `service_name`. |
| service_name | string | Service name of the instance. MUST be a non-empty string. MUST NOT be used in conjunction with the fields `service_id` and `plan_ids`. |
| plan_names | array-of-strings | An array of plan names. The instance must match one of those plans. MUST be a non-empty array of non-empty strings. REQUIRES the field `service_name`. |

\* Fields with an asterisk are REQUIRED.

### Instance Dependencies

| Field | Type | Description |
| --- | --- | --- |
| parameter* | string | Name of the parameter. MUST be a non-empty string and MUST be unique within all dependency parameters. |
| description* | string | A short description of the dependency. MUST be a non-empty string. |
| optional  | boolean | Whether the dependency is optional or not. Defaults to false. |
| instance\_id\_parameter | string | The name of the parameter that contains the service instance ID of an existing service. MUST be a non-empty string. MUST NOT be used in conjunction with the field `alias_parameter`. |
| alias_parameter | string | The name of the parameter that contains the alias of an existing service. MUST be a non-empty string. MUST NOT be used in conjunction with the field `instance_id_parameter`. |

\* Fields with an asterisk are REQUIRED.

```
{
  "dependencies": {
    "open_dependencies_parameter": "deps",
    "services": [
      {
        "parameter": "uuid-service",
        "description": "UUID Service to create UUIs.",
        "service_id": "4ec72acd-b8fa-4be0-bdc2-ffb2d1104185",
        "optional": true
      },
      {
        "parameter": "random",
        "description": "Random numbers.",
        "plan_ids": [
          "80d96caa-8ae7-41f7-9e80-349de02397df"
        ],
        "optional": false
      },
      {
        "parameter": "memcache",
        "description": "Cache server.",
        "service_name": "memcache"
      },
      {
        "parameter": "database",
        "description": "Main database instance.",
        "service_name": "mysql",
        "plan_names": [
          "big",
          "huge"
        ]
      }
    ],
    "instances": [
      {
        "parameter": "queue",
        "description": "Message queue.",
        "alias_parameter": "queue_alias"
      }
    ]
  }
}
```

## Synchronous and Asynchronous Operations

Platforms expect prompt responses to all API requests in order to provide
users with fast feedback. Service Broker authors SHOULD implement their
Service Brokers to respond promptly to all requests but will need to decide
whether to implement synchronous or asynchronous responses. Service Brokers
that can guarantee completion of the requested operation with the response
SHOULD return the synchronous response. Service Brokers that cannot guarantee
completion of the operation with the response SHOULD implement the
asynchronous response.

Providing a synchronous response for a provision, update, or bind operation
before actual completion causes confusion for users as their service might not
be usable and they have no way to find out when it will be. Asynchronous
responses set expectations for users that an operation is in progress and can
also provide updates on the status of the operation.

Support for synchronous or asynchronous responses MAY vary by Service
Offering, even by Service Plan.

### Synchronous Operations

To execute a request synchronously, the Service Broker need only return the
usual status codes: `201 Created` for provision and bind, and `200 OK` for
update, unbind, and deprovision.

Service Brokers that support synchronous responses for provision, update, and
delete can ignore the `accepts_incomplete=true` query parameter if it is
provided by the client.

### Asynchronous Operations

Note: Asynchronous operations are currently supported only for provision,
update, and deprovision.

For a Service Broker to return an asynchronous response, the query parameter
`accepts_incomplete=true` MUST be included the request. If the parameter is not
included or is set to `false`, and the Service Broker cannot fulfil the request
synchronously (guaranteeing that the operation is complete on response), then
the Service Broker SHOULD reject the request with the status code `422
Unprocessable Entity` and a response body containing error code
`"AsyncRequired"` (see [Service Broker Errors](#service-broker-errors)). The
error response MAY include a helpful error message in the `description` field
such as `"This Service Plan requires client support for asynchronous service
operations."`.

If the query parameter described above is present, and the Service Broker
executes the request asynchronously, the Service Broker MUST return the
asynchronous response `202 Accepted`. The response body SHOULD be the same as
if the Service Broker were serving the request synchronously.

An asynchronous response triggers the Platform to poll the endpoint
`GET /v2/service_instances/:instance_id/last_operation` until the Service Broker
indicates that the requested operation has succeeded or failed. Service Brokers
MAY include a status message with each response for the `last_operation`
endpoint that provides visibility to end users as to the progress of the
operation.

## Blocking Operations

Service Brokers do not have to support concurrent requests that mutate the same
resource. If a Service Broker receives a request that it is not able to process
due to other activity being done on that resource then the Service Broker MUST
reject the request with a HTTP `422 Unprocessable Entity` and a response body
containing error code `"ConcurrencyError"` (see
[Service Broker Errors](#service-broker-errors)). The error response MAY include
a helpful error message in the `description` field such as `"Another operation
for this Service Instance is in progress."`.

Note that per the [Orphans](#orphans) section, this error response does not
cause orphan mitigation to be initiated. Therefore, Platforms receiving this
error response SHOULD resend the request at a later time.

Brokers MAY choose to treat the creation of a binding as a mutation of the
corresponding Service Instance - it is an implementation choice. Doing so would
cause Platforms to serialize multiple binding creation requests when they are
directed at the same Service Instance if concurrent updates are not supported.

## Polling Last Operation

When a Service Broker returns status code `202 Accepted` for
[Provision](#provisioning), [Update](#updating-a-service-instance), or
[Deprovision](#deprovisioning), the Platform will begin polling the
`/v2/service_instances/:instance_id/last_operation` endpoint to obtain the
state of the last requested operation. The Service Broker response MUST contain
the field `state` and MAY contain the field `description`.

Valid values for `state` are `in progress`, `succeeded`, and `failed`. The
Platform will poll the `last_operation` endpoint as long as the Service Broker
returns `"state": "in progress"`. Returning `"state": "succeeded"` or `"state":
"failed"` will cause the Platform to cease polling. The value provided for
`description` will be passed through to the Platform API client and can be used
to provide additional detail for users about the progress of the operation.

### Request

#### Route
`GET /v2/service_instances/:instance_id/last_operation`

`:instance_id` MUST be a globally unique non-empty string.

#### Parameters

The request provides these query string parameters as useful hints for Service Brokers.

| Query-String Field | Type | Description |
| --- | --- | --- |
| service_id | string | If present, it MUST be the ID of the service being used. |
| plan_id | string | If present, it MUST be the ID of the plan for the service being use. |
| operation | string | A Service Broker-provided identifier for the operation. When a value for `operation` is included with asynchronous responses for [Provision](#provisioning), [Update](#updating-a-service-instance), and [Deprovision](#deprovisioning) requests, the Platform MUST provide the same value using this query parameter as a percent-encoded string. If present, MUST be a non-empty string. |

Note: Although the request query parameters `service_id` and `plan_id` are not
mandatory, the Platform SHOULD include them on all `last_operation` requests
it makes to Service Brokers.

#### Headers

The following HTTP Headers are defined for this operation:

| Header | Type | Description |
| --- | --- | --- |
| X-Broker-API-Version* | string | See [API Version Header](#api-version-header). |

\* Headers with an asterisk are REQUIRED.

#### cURL
```
$ curl http://username:password@service-broker-url/v2/service_instances/:instance_id/last_operation -H "X-Broker-API-Version: 2.13"
```

### Response

| Status Code | Description |
| --- | --- |
| 200 OK | MUST be returned upon successful processing of this request. The expected response body is below. |
| 400 Bad Request | MUST be returned if the request is malformed or missing mandatory data. |
| 410 Gone | Appropriate only for asynchronous delete operations. The Platform MUST consider this response a success and forget about the resource. The expected response body is `{}`. Returning this while the Platform is polling for create or update operations SHOULD be interpreted as an invalid response and the Platform SHOULD continue polling. |

Responses with any other status code SHOULD be interpreted as an error or
invalid response. The Platform SHOULD continue polling until the Service Broker
returns a valid response or the
[maximum polling duration](#polling-interval-and-duration) is reached.

#### Body

For success responses, the following fields are defined:

| Response Field | Type | Description |
| --- | --- | --- |
| state* | string | Valid values are `in progress`, `succeeded`, and `failed`. While `"state": "in progress"`, the Platform SHOULD continue polling. A response with `"state": "succeeded"` or `"state": "failed"` MUST cause the Platform to cease polling. |
| description | string | A user-facing message displayed to the Platform API client. Can be used to tell the user details about the status of the operation. If present, MUST be a non-empty string. |

\* Fields with an asterisk are REQUIRED.

```
{
  "state": "in progress",
  "description": "Creating service (10% complete)."
}
```

If the successful response includes a `state` of `failed` then the Platform
MUST send a deprovision request to the Service Broker to prevent an orphan
being created on the Service Broker. However, while the Platform will attempt
to send a deprovision request, Service Brokers MAY automatically delete
any resources associated with the failed bind request on their own.

### Polling Interval and Duration

The frequency and maximum duration of polling MAY vary by Platform client. If
a Platform has a max polling duration and this limit is reached, the Platform
MUST cease polling and the operation state MUST be considered `failed`.

## Provisioning

When the Service Broker receives a provision request from the Platform, it
MUST take whatever action is necessary to create a new resource. What
provisioning represents varies by service and plan, although there are several
common use cases. For a MySQL service, provisioning could result in an empty
dedicated database server running on its own VM or an empty schema on a shared
database server. For non-data services, provisioning could just mean an
account on an multi-tenant SaaS application.

### Request

#### Route
`PUT /v2/service_instances/:instance_id`

`:instance_id` MUST be a globally unique non-empty string.
This ID will be used for future requests (bind and deprovision), so the
Service Broker will use it to correlate the resource it creates.

#### Parameters
| Parameter Name | Type | Description |
| --- | --- | --- |
| accepts_incomplete | boolean | A value of true indicates that the Platform and its clients support asynchronous Service Broker operations. If this parameter is not included in the request, and the Service Broker can only provision a Service Instance of the requested plan asynchronously, the Service Broker MUST reject the request with a `422 Unprocessable Entity` as described below. |

#### Headers

The following HTTP Headers are defined for this operation:

| Header | Type | Description |
| --- | --- | --- |
| X-Broker-API-Version* | string | See [API Version Header](#api-version-header). |
| X-Broker-API-Originating-Identity | string | See [Originating Identity](#originating-identity). |

\* Headers with an asterisk are REQUIRED.

#### Body
| Request Field | Type | Description |
| --- | --- | --- |
| service_id* | string | MUST be the ID of a service from the catalog for this Service Broker. |
| plan_id* | string | MUST be the ID of a plan from the service that has been requested. |
| context | object | Platform specific contextual information under which the Service Instance is to be provisioned. Although most Service Brokers will not use this field, it could be helpful in determining data placement or applying custom business rules. `context` will replace `organization_guid` and `space_guid` in future versions of the specification - in the interim both SHOULD be used to ensure interoperability with old and new implementations. |
| organization_guid* | string | Deprecated in favor of `context`. The Platform GUID for the organization under which the Service Instance is to be provisioned. Although most Service Brokers will not use this field, it might be helpful for executing operations on a user's behalf. MUST be a non-empty string. |
| space_guid* | string | Deprecated in favor of `context`. The identifier for the project space within the Platform organization. Although most Service Brokers will not use this field, it might be helpful for executing operations on a user's behalf. MUST be a non-empty string. |
| parameters | object | Configuration options for the Service Instance. Service Brokers SHOULD ensure that the client has provided valid configuration parameters and values for the operation. |
| alias | string | Alias for this instance. |

\* Fields with an asterisk are REQUIRED.

```
{
  "service_id": "service-id-here",
  "plan_id": "plan-id-here",
  "context": {
    "platform": "cloudfoundry",
    "some_field": "some-contextual-data"
  },
  "organization_guid": "org-guid-here",
  "space_guid": "space-guid-here",
  "parameters": {
    "parameter1": 1,
    "parameter2": "foo"
  }
}
```

#### cURL
```
$ curl http://username:password@service-broker-url/v2/service_instances/:instance_id?accepts_incomplete=true -d '{
  "service_id": "service-id-here",
  "plan_id": "plan-id-here",
  "context": {
    "platform": "cloudfoundry",
    "some_field": "some-contextual-data"
  },
  "organization_guid": "org-guid-here",
  "space_guid": "space-guid-here",
  "parameters": {
    "parameter1": 1,
    "parameter2": "foo"
  }
}' -X PUT -H "X-Broker-API-Version: 2.13" -H "Content-Type: application/json"
```

### Response

| Status Code | Description |
| --- | --- |
| 200 OK | MUST be returned if the Service Instance already exists, is fully provisioned, and the requested parameters are identical to the existing Service Instance. The expected response body is below. |
| 201 Created | MUST be returned if the Service Instance was provisioned as a result of this request. The expected response body is below. |
| 202 Accepted | MUST be returned if the Service Instance provisioning is in progress. The `operation` string MUST match that returned for the original request. This triggers the Platform to poll the [Service Instance Last Operation Endpoint](#polling-last-operation) for operation status. Note that a re-sent `PUT` request MUST return a `202 Accepted`, not a `200 OK`, if the Service Instance is not yet fully provisioned. |
| 400 Bad Request | MUST be returned if the request is malformed or missing mandatory data. |
| 409 Conflict | MUST be returned if a Service Instance with the same id already exists but with different attributes. |
| 422 Unprocessable Entity | MUST be returned if the Service Broker only supports asynchronous provisioning for the requested plan and the request did not include `?accepts_incomplete=true`. The response body MUST contain a response body containing error code `"AsyncRequired"` (see [Service Broker Errors](#service-broker-errors)). The error response MAY include a helpful error message in the `description` field such as `"This Service Plan requires client support for asynchronous service operations."`. |

Responses with any other status code will be interpreted as a failure and a
deprovision request MUST be sent to the Service Broker to prevent an orphan
being created on the Service Broker. However, while the platform will attempt
to send a deprovision request, Service Brokers MAY automatically delete
any resources associated with the failed provisioning request on their own.

#### Body

For success responses, the following fields are defined:

| Response Field | Type | Description |
| --- | --- | --- |
| dashboard_url | string | The URL of a web-based management user interface for the Service Instance; we refer to this as a service dashboard. The URL MUST contain enough information for the dashboard to identify the resource being accessed (`9189kdfsk0vfnku` in the example below). Note: a Service Broker that wishes to return `dashboard_url` for a Service Instance MUST return it with the initial response to the provision request, even if the service is provisioned asynchronously. If present, MUST be a non-empty string. |
| operation | string | For asynchronous responses, Service Brokers MAY return an identifier representing the operation. The value of this field MUST be provided by the Platform with requests to the [Last Operation](#polling-last-operation) endpoint in a percent-encoded query parameter. If present, MUST be a non-empty string. |

```
{
  "dashboard_url": "http://example-dashboard.example.com/9189kdfsk0vfnku",
  "operation": "task_10"
}
```

## Updating a Service Instance

By implementing this endpoint, Service Broker authors can enable users to
modify two attributes of an existing Service Instance: the Service Plan and
parameters. By changing the Service Plan, users can upgrade or downgrade their
Service Instance to other plans. By modifying parameters, users can change
configuration options that are specific to a service or plan.

To enable support for the update of the plan, a Service Broker MUST declare
support per service by including `"plan_updateable": true` in its [catalog
endpoint](#catalog-management).

If `"plan_updateable": true` is declared for a plan in the
[Catalog](#catalog-management) endpoint, the Platform MAY request an update of
a Service Instance using the given plan. Otherwise, Platforms MUST NOT make an
update request to the Service Broker for any Service Instance using the given
plan.

Not all permutations of plan changes are expected to be supported. For
example, a service might support upgrading from plan "shared small" to "shared
large" but not to plan "dedicated". It is up to the Service Broker to validate
whether a particular permutation of plan change is supported. If a particular
plan change is not supported, the Service Broker SHOULD return a meaningful
error message in response.

### Request

#### Route
`PATCH /v2/service_instances/:instance_id`

`:instance_id` MUST be the ID a previously provisioned Service Instance.

#### Parameters
| Parameter Name | Type | Description |
| --- | --- | --- |
| accepts_incomplete | boolean | A value of true indicates that the Platform and its clients support asynchronous Service Broker operations. If this parameter is not included in the request, and the Service Broker can only provision a Service Instance of the requested plan asynchronously, the Service Broker SHOULD reject the request with a `422 Unprocessable Entity` as described below. |

#### Headers

The following HTTP Headers are defined for this operation:

| Header | Type | Description |
| --- | --- | --- |
| X-Broker-API-Version* | string | See [API Version Header](#api-version-header). |
| X-Broker-API-Originating-Identity | string | See [Originating Identity](#originating-identity). |

\* Headers with an asterisk are REQUIRED.

#### Body

| Request Field | Type | Description |
| --- | --- | --- |
| context | object | Contextual data under which the Service Instance is created. |
| service_id* | string | MUST be the ID of a service from the catalog for this Service Broker. |
| plan_id | string | If present, MUST be the ID of a plan from the service that has been requested. If this field is not present in the request message, then the Service Broker MUST NOT change the plan of the instance as a result of this request. |
| parameters | object | Configuration options for the Service Instance. Service Brokers SHOULD ensure that the client has provided valid configuration parameters and values for the operation. See "Note" below. |
| previous_values | [PreviousValues](#previous-values-object) | Information about the Service Instance prior to the update. |

\* Fields with an asterisk are REQUIRED.

##### Previous Values Object

| Request Field | Type | Description |
| --- | --- | --- |
| service_id | string | Deprecated; determined to be unnecessary as the value is immutable. If present, it MUST be the ID of the service for the Service Instance. |
| plan_id | string | If present, it MUST be the ID of the plan prior to the update. |
| organization_id | string | Deprecated as it was redundant information. Organization for the Service Instance MUST be provided by Platforms in the top-level field `context`. If present, it MUST be the ID of the organization specified for the Service Instance. |
| space_id | string | Deprecated as it was redundant information. Space for the Service Instance MUST be provided by Platforms in the top-level field `context`. If present, it MUST be the ID of the space specified for the Service Instance. |

Note: The `parameters` specified are expected to be the values specified
by an end-user. Whether the user chooses to include the complete set of
configuration options or just a subset (or even none) is their choice. How a
Service Broker interprets these parameters (including the absence of any
particular parameter) is out of scope of this specification - with the
exception that if this field is not present in the request then the
Service Broker MUST NOT change the parameters of the instance as a result of
this request.

Since some Service Instances will provide a `dashboard_url`, it is possible
that a user has modified some of these parameters via the dashboard and
therefore the Platform might not be aware of these changes. For this reason,
Platforms SHOULD NOT include any parameters on the request that
the user did not explicitly specify in their request for the update.

\* Fields with an asterisk are REQUIRED.

```
{
  "context": {
    "platform": "cloudfoundry",
    "some_field": "some-contextual-data"
  },
  "service_id": "service-id-here",
  "plan_id": "plan-id-here",
  "parameters": {
    "parameter1": 1,
    "parameter2": "foo"
  },
  "previous_values": {
    "plan_id": "old-plan-id-here",
    "service_id": "service-id-here",
    "organization_id": "org-guid-here",
    "space_id": "space-guid-here"
  }
}
```

#### cURL
```
$ curl http://username:password@service-broker-url/v2/service_instances/:instance_id?accepts_incomplete=true -d '{
  "context": {
    "platform": "cloudfoundry",
    "some_field": "some-contextual-data"
  },
  "service_id": "service-id-here",
  "plan_id": "plan-id-here",
  "parameters": {
    "parameter1": 1,
    "parameter2": "foo"
  },
  "previous_values": {
    "plan_id": "old-plan-id-here",
    "service_id": "service-id-here",
    "organization_id": "org-guid-here",
    "space_id": "space-guid-here"
  }
}' -X PATCH -H "X-Broker-API-Version: 2.13" -H "Content-Type: application/json"
```

### Response

| Status Code | Description |
| --- | --- |
| 200 OK | MUST be returned if the request's changes have been applied. The expected response body is `{}`. |
| 202 Accepted | MUST be returned if the Service Instance update is in progress. The `operation` string MUST match that returned for the original request. This triggers the Platform to poll the [Last Operation](#polling-last-operation) for operation status. Note that a re-sent `PATCH` request MUST return a `202 Accepted`, not a `200 OK`, if the requested update has not yet completed. |
| 400 Bad Request | MUST be returned if the request is malformed or missing mandatory data. |
| 422 Unprocessable entity | MUST be returned if the requested change is not supported or if the request cannot currently be fulfilled due to the state of the Service Instance (e.g. Service Instance utilization is over the quota of the requested plan). Additionally, a `422 Unprocessable Entity` can also be returned if the Service Broker only supports asynchronous update for the requested plan and the request did not include `?accepts_incomplete=true`; in this case the response body MUST contain a error code `"AsyncRequired"` (see [Service Broker Errors](#service-broker-errors)). The error response MAY include a helpful error message in the `description` field such as `"This Service Plan requires client support for asynchronous service operations."`. |

Responses with any other status code will be interpreted as a failure.
When the response includes a 4xx status code, the Service Broker MUST NOT
apply any of the requested changes to the Service Instance.

#### Body

For success responses, the following fields are defined:

| Response Field | Type | Description |
| --- | --- | --- |
| operation | string | For asynchronous responses, Service Brokers MAY return an identifier representing the operation. The value of this field MUST be provided by the Platform with requests to the [Last Operation](#polling-last-operation) endpoint in a percent-encoded query parameter. If present, MUST be a non-empty string. |

```
{
  "operation": "task_10"
}
```


## Binding

If `"bindable": true` is declared for a service or plan in the
[Catalog](#catalog-management) endpoint, the Platform MAY request generation
of a Service Binding. Otherwise, Platforms MUST NOT make a binding request to
the Service Broker for any Service Instance using the given service or plan.

Note: Not all services need to be bindable --- some deliver value just from
being provisioned. Service Brokers that offer services that are bindable MUST
declare them as such using `"bindable": true` in the
[Catalog](#catalog-management). Service Brokers that do not offer any bindable
services do not need to implement the endpoint for bind requests.

### Types of Binding

#### Credentials

Credentials are a set of information used by an Application or a user to
utilize the Service Instance. If the Service Broker supports generation of
credentials it MUST return `credentials` in the response for a request to
create a Service Binding. Credentials SHOULD be unique whenever possible, so
access can be revoked for each binding without affecting consumers of other
bindings for the Service Instance.

#### Log Drain

There are a class of Service Offerings that provide aggregation, indexing, and
analysis of log data. To utilize these services an application that generates
logs needs information for the location to which it will stream logs. A create
binding response from a Service Broker that provides one of these services
MUST include a `syslog_drain_url`. The Platform MUST use the
`syslog_drain_url` value when sending logs to the service.

Service Brokers MUST NOT include a `syslog_drain_url` in a create binding
response if the associated [Catalog](#catalog-management) entry for the
service did not include a `"requires":["syslog_drain"]` property.

#### Route Services
Route services are a class of Service Offerings that intermediate requests to
applications, performing functions such as rate limiting or authorization. To
indicate support for route services, the catalog entry for the Service MUST
include the `"requires":["route_forwarding"]` property.

When creating a route service type of Service Binding, a Platform MUST send
a routable address, or endpoint, for the application along with the request to
create a Service Binding using `"bind_resource":{"route":"some-address.com"}`.
Service Brokers MAY support configuration specific to an address using
parameters; exposing this feature to users would require a Platform to support
binding multiple routable addresses to the same Service Instance.

If a service is deployed in a configuration to support this behavior, the
Service Broker MUST return a `route_service_url` in the response for a request
to create a binding, so that the Platform knows where to proxy the application
request. If the service is deployed such that the network configuration to
proxy application requests through instances of the service is managed
out-of-band, the Service Broker MUST NOT return `route_service_url` in the
response.

Service Brokers MUST NOT include a `route_service_url` in a create binding
response if the associated [Catalog](#catalog-management) entry for the
service did not include a `"requires":["route_forwarding"]` property.

#### Volume Services

There are a class of services that provide network storage to applications
via volume mounts in the application container. A create binding response
from one of these services MUST include `volume_mounts`.

Service Brokers MUST NOT include `volume_mounts` in a create binding response
if the associated [Catalog](#catalog-management) entry for the service
did not include a `"requires":["volume_mount"]` property.

### Request

#### Route
`PUT /v2/service_instances/:instance_id/service_bindings/:binding_id`

`:instance_id` MUST be the ID of a previously provisioned Service Instance.

`:binding_id` MUST be a globally unique non-empty string.
This ID will be used for future unbind requests, so the Service Broker will use
it to correlate the resource it creates.

#### Headers

The following HTTP Headers are defined for this operation:

| Header | Type | Description |
| --- | --- | --- |
| X-Broker-API-Version* | string | See [API Version Header](#api-version-header). |
| X-Broker-API-Originating-Identity | string | See [Originating Identity](#originating-identity). |

\* Headers with an asterisk are REQUIRED.

#### Body

| Request Field | Type | Description |
| --- | --- | --- |
| context | object | Contextual data under which the Service Binding is created. |
| service_id* | string | MUST be the ID of the service that is being used. |
| plan_id* | string | MUST be the ID of the plan from the service that is being used. |
| app_guid | string | Deprecated in favor of `bind_resource.app_guid`. GUID of an application associated with the binding to be created. If present, MUST be a non-empty string. |
| bind_resource | [BindResource](#bind-resource-object) | A JSON object that contains data for Platform resources associated with the binding to be created. See [Bind Resource Object](#bind-resource-object) for more information. |
| parameters | object | Configuration options for the Service Binding. Service Brokers SHOULD ensure that the client has provided valid configuration parameters and values for the operation. |

\* Fields with an asterisk are REQUIRED.

##### Bind Resource Object

The `bind_resource` object contains Platform specific information related to
the context in which the service will be used. In some cases the Platform
might not be able to provide this information at the time of the binding
request, therefore the `bind_resource` and its fields are OPTIONAL.

Below are some common fields that MAY be used. Platforms MAY choose to add
additional ones as needed (see
[Bind Resource Object](profile.md#bind-resource-object) conventions).

| Request Field | Type | Description |
| --- | --- | --- |
| app_guid | string | GUID of an application associated with the binding. For [credentials](#types-of-binding) bindings. MUST be unique within the scope of the Platform. |
| route | string | URL of the application to be intermediated. For [route services](#route-services) bindings. |

`app_guid` represents the scope to which the binding will apply within
the Platform. For example, in Cloud Foundry it might map to an "application"
while in Kubernetes it might map to a "namespace". The scope of what a
Platform maps the `app_guid` to is Platform specific and MAY vary across
binding requests.

```
{
  "context": {
    "platform": "cloudfoundry",
    "some_field": "some-contextual-data"
  },
  "service_id": "service-id-here",
  "plan_id": "plan-id-here",
  "bind_resource": {
    "app_guid": "app-guid-here"
  },
  "parameters": {
    "parameter1-name-here": 1,
    "parameter2-name-here": "parameter2-value-here"
  }
}
```


#### cURL
```
$ curl http://username:password@service-broker-url/v2/service_instances/:instance_id/service_bindings/:binding_id -d '{
  "context": {
    "platform": "cloudfoundry",
    "some_field": "some-contextual-data"
  },
  "service_id": "service-id-here",
  "plan_id": "plan-id-here",
  "bind_resource": {
    "app_guid": "app-guid-here"
  },
  "parameters": {
    "parameter1-name-here": 1,
    "parameter2-name-here": "parameter2-value-here"
  }
}' -X PUT -H "X-Broker-API-Version: 2.13" -H "Content-Type: application/json"
```

### Response

| Status Code | Description |
| --- | --- |
| 200 OK | MUST be returned if the binding already exists and the requested parameters are identical to the existing binding. The expected response body is below. |
| 201 Created | MUST be returned if the binding was created as a result of this request. The expected response body is below. |
| 400 Bad Request | MUST be returned if the request is malformed or missing mandatory data. |
| 409 Conflict | MUST be returned if a Service Binding with the same id, for the same Service Instance, already exists but with different parameters. |
| 422 Unprocessable Entity | MUST be returned if the Service Broker requires that `app_guid` be included in the request body. The response body MUST contain error code `"RequiresApp"` (see [Service Broker Errors](#service-broker-errors)). The error response MAY include a helpful error message in the `description` field such as `"This Service supports generation of credentials through binding an application only."`. Additionally, if the Service Broker rejects the request due to a concurrent request to create a binding for the same Service Instance, then this error MUST be returned (see [Blocking Operations](#blocking-operations)). |

Responses with any other status code will be interpreted as a failure and an
unbind request MUST be sent to the Service Broker to prevent an orphan being
created on the Service Broker. However, while the platform will attempt
to send an unbind request, Service Brokers MAY automatically delete
any resources associated with the failed bind request on their own.

#### Body

For success responses, the following fields are defined:

| Response Field | Type | Description |
| --- | --- | --- |
| credentials | object | A free-form hash of credentials that can be used by applications or users to access the service. |
| syslog_drain_url | string | A URL to which logs MUST be streamed. `"requires":["syslog_drain"]` MUST be declared in the [Catalog](#catalog-management) endpoint or the Platform MUST consider the response invalid. |
| route_service_url | string | A URL to which the Platform MUST proxy requests for the address sent with `bind_resource.route` in the request body. `"requires":["route_forwarding"]` MUST be declared in the [Catalog](#catalog-management) endpoint or the Platform can consider the response invalid. |
| volume_mounts | array of [VolumeMount](#volume-mount-object) objects | An array of configuration for remote storage devices to be mounted into an application container filesystem. `"requires":["volume_mount"]` MUST be declared in the [Catalog](#catalog-management) endpoint or the Platform can consider the response invalid. |

##### Volume Mount Object

| Response Field | Type | Description |
| --- | --- | --- |
| driver* | string | Name of the volume driver plugin which manages the device. |
| container_dir* | string | The path in the application container onto which the volume will be mounted. This specification does not mandate what action the Platform is to take if the path specified already exists in the container. |
| mode* | string | "r" to mount the volume read-only or "rw" to mount it read-write. |
| device_type* | string | A string specifying the type of device to mount. Currently the only supported value is "shared". |
| device* | [Device](#device-object) | Device object containing device_type specific details. Currently only shared devices are supported. |

\* Fields with an asterisk are REQUIRED.

##### Device Object

Currently only shared devices are supported; a distributed file system which
can be mounted on all app instances simultaneously.

| Field | Type | Description |
| --- | --- | --- |
| volume_id* | string | ID of the shared volume to mount on every app instance. |
| mount_config | object | Configuration object to be passed to the driver when the volume is mounted. |

\* Fields with an asterisk are REQUIRED.

```
{
  "credentials": {
    "uri": "mysql://mysqluser:pass@mysqlhost:3306/dbname",
    "username": "mysqluser",
    "password": "pass",
    "host": "mysqlhost",
    "port": 3306,
    "database": "dbname"
  }
}
```

```
{
  "volume_mounts": [{
    "driver": "cephdriver",
    "container_dir": "/data/images",
    "mode": "r",
    "device_type": "shared",
    "device": {
      "volume_id": "bc2c1eab-05b9-482d-b0cf-750ee07de311",
      "mount_config": {
        "key": "value"
      }
    }
  }]
}
```

## Unbinding

Note: Service Brokers that do not provide any bindable services or plans do
not need to implement this endpoint.

When a Service Broker receives an unbind request from a Platform, it MUST
delete any resources associated with the binding. In the case where
credentials were generated, this might result in requests to the Service
Instance failing to authenticate.

### Request

#### Route

`DELETE /v2/service_instances/:instance_id/service_bindings/:binding_id`

`:instance_id` MUST be the ID of a previously provisioned Service Instance.

`:binding_id` MUST be the the ID of a previously provisioned binding for that
Service Instance.

#### Parameters

The request provides these query string parameters as useful hints for Service
Brokers.

| Query-String Field | Type | Description |
| --- | --- | --- |
| service_id* | string | MUST be the ID of the service associated with the binding being deleted. |
| plan_id* | string | MUST be the ID of the plan associated with the binding being deleted. |

\* Query parameters with an asterisk are REQUIRED.

#### Headers

The following HTTP Headers are defined for this operation:

| Header | Type | Description |
| --- | --- | --- |
| X-Broker-API-Version* | string | See [API Version Header](#api-version-header). |
| X-Broker-API-Originating-Identity | string | See [Originating Identity](#originating-identity). |

\* Headers with an asterisk are REQUIRED.

#### cURL

```
$ curl 'http://username:password@service-broker-url/v2/service_instances/:instance_id/
  service_bindings/:binding_id?service_id=service-id-here&plan_id=plan-id-here' -X DELETE -H "X-Broker-API-Version: 2.13"
```

### Response

| Status Code | Description |
| --- | --- |
| 200 OK | MUST be returned if the binding was deleted as a result of this request. The expected response body is `{}`. |
| 400 Bad Request | MUST be returned if the request is malformed or missing mandatory data. |
| 410 Gone | MUST be returned if the binding does not exist. The expected response body is `{}`. |

Responses with any other status code will be interpreted as a failure and the
Platform MUST continue to remember the Service Binding.

#### Body

For a success response, the expected response body is `{}`.

## Deprovisioning

When a Service Broker receives a deprovision request from a Platform, it MUST
delete any resources it created during the provision. Usually this means that
all resources are immediately reclaimed for future provisions.

Platforms MUST delete all bindings for a service prior to attempting to
deprovision the service. This specification does not specify what a Service
Broker is to do if it receives a deprovision request while there are still
bindings associated with it.

### Request

#### Route

`DELETE /v2/service_instances/:instance_id`

`:instance_id` MUST be the ID of a previously provisioned Service Instance.

#### Parameters

The request provides these query string parameters as useful hints for Service
Brokers.


| Query-String Field | Type | Description |
| --- | --- | --- |
| service_id* | string | MUST be the ID of the Service Instance being deleted. |
| plan_id* | string | MUST be the ID of the plan associated with the Service Instance being deleted. |
| accepts_incomplete | boolean | A value of true indicates that both the Platform and the requesting client support asynchronous deprovisioning. If this parameter is not included in the request, and the Service Broker can only deprovision a Service Instance of the requested plan asynchronously, the Service Broker MUST reject the request with a `422 Unprocessable Entity` as described below. |

\* Query parameters with an asterisk are REQUIRED.

#### Headers

The following HTTP Headers are defined for this operation:

| Header | Type | Description |
| --- | --- | --- |
| X-Broker-API-Version* | string | See [API Version Header](#api-version-header). |
| X-Broker-API-Originating-Identity | string | See [Originating Identity](#originating-identity). |

\* Headers with an asterisk are REQUIRED.

#### cURL
```
$ curl 'http://username:password@service-broker-url/v2/service_instances/:instance_id?accepts_incomplete=true
  &service_id=service-id-here&plan_id=plan-id-here' -X DELETE -H "X-Broker-API-Version: 2.13"
```

### Response

| Status Code | Description |
| --- | --- |
| 200 OK | MUST be returned if the Service Instance was deleted as a result of this request. The expected response body is `{}`. |
| 202 Accepted | MUST be returned if the Service Instance deletion is in progress. The `operation` string MUST match that returned for the original request. This triggers the Platform to poll the [Service Instance Last Operation Endpoint](#polling-last-operation) for operation status. Note that a re-sent `DELETE` request MUST return a `202 Accepted`, not a `200 OK`, if the delete request has not completed yet. |
| 400 Bad Request | MUST be returned if the request is malformed or missing mandatory data. |
| 410 Gone | MUST be returned if the Service Instance does not exist. |
| 422 Unprocessable Entity | MUST be returned if the Service Broker only supports asynchronous deprovisioning for the requested plan and the request did not include `?accepts_incomplete=true`. The response body MUST contain error code `"AsyncRequired"` (see [Service Broker Errors](#service-broker-errors)). The error response MAY include a helpful error message in the `description` field such as `"This Service Plan requires client support for asynchronous service operations."`. |

Responses with any other status code will be interpreted as a failure and the
Platform MUST remember the Service Instance.

#### Body

For success responses, the following fields are defined:

| Response Field | Type | Description |
| --- | --- | --- |
| operation | string | For asynchronous responses, Service Brokers MAY return an identifier representing the operation. The value of this field MUST be provided by the Platform with requests to the [Last Operation](#polling-last-operation) endpoint in a percent-encoded query parameter. If present, MUST be a non-empty string. |

```
{
  "operation": "task_10"
}
```

## Orphans

The Platform is the source of truth for Service Instances and Service Bindings.
Service Brokers are expected to have successfully provisioned all of the Service
Instances and Service Bindings that the Platform knows about, and none that it
doesn't.

Orphans can result if the Service Broker does not return a response before a
request from the Platform times out (typically 60 seconds). For example, if a
Service Broker does not return a response to a provision request before the
request times out, the Service Broker might eventually succeed in provisioning
a Service Instance after the Platform considers the request a failure. This
results in an orphan Service Instance on the Service Broker's side.

To mitigate orphan Service Instances and Service Bindings, the Platform SHOULD
attempt to delete resources it cannot be sure were successfully created, and
SHOULD keep trying to delete them until the Service Broker responds with a
success.

Platforms SHOULD initiate orphan mitigation in the following scenarios:

| Status Code Of Service Broker Response | Platform Interpretation Of Response | Platform Initiates Orphan Mitigation? |
| --- | --- | --- |
| 200 | Success | No |
| 200 with malformed response | Failure | No |
| 201 | Success | No |
| 201 with malformed response | Failure | Yes |
| All other 2xx | Failure | Yes |
| 408 | Client timeout failure (request not received at the server) | No |
| All other 4xx | Request rejected | No |
| 5xx | Service Broker error | Yes |
| Timeout | Failure | Yes |

If the Platform encounters an internal error provisioning a Service Instance or
Service Binding (for example, saving to the database fails), then it MUST at
least send a single delete or unbind request to the Service Broker to prevent
the creation of an orphan.
