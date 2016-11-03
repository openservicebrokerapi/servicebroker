# Service Broker Specification

# Abstract

<!-- **TODO 'managed services' terminology is different from everywhere else** -->


The documentation in this section is intended for developers and operators interested in creating Managed Services for Cloud Foundry. Managed Services are defined as having been integrated with Cloud Foundry via APIs, and enable end users to provision reserved resources and credentials on demand. For documentation targeted at end users, such as how to provision services and integrate them with applications, see [Services Overview](../devguide/services/index.html).

To develop Managed Services for Cloud Foundry, you'll need a Cloud Foundry instance to test your service broker with as you are developing it. You must have admin access to your CF instance to manage service brokers and the services marketplace catalog. For local development, we recommend using [BOSH Lite](https://github.com/cloudfoundry/bosh-lite) to deploy your own local instance of Cloud Foundry.

**Table of Contents**

- [Overview](#overview)
- [API](#api)
- [Catalog Metadata](#catalog-metadata)
- [Binding Credentials](#binding-credentials)
- [Examples](#examples)
- [Managing Service Brokers](#managing-service-brokers)
- [Access Control](#access-control)
- [Dashboard Single Sign-On](#dashboard-single-sign-on)
- [Application Log Streaming](#application-log-streaming)
- [Route Services](#route-services)
- [Manage Application Requests with Route Services](#manage-application-requests-with-route-services)
- [Supporting Multiple Cloud Foundry Instances](#supporting-multiple-cloud-foundry-instances)
- [Volume Services (Experimental)](#volume-services-experimental)
- [Volume Services (Experimental/Obsolete)](#volume-services-experimentalobsolete)
  
# Overview

<!-- **TODO delete this** -->

Services are integrated with Cloud Foundry by implementing a documented API for which the cloud controller is the client; we call this the Service Broker API. This should not be confused with the cloud controller API, often used to refer to the version of Cloud Foundry itself; when one refers to "Cloud Foundry v2" they are referring to the version of the cloud controller API. The services API is versioned independently of the cloud controller API.

<!-- **TODO keep the first sentence, ignore gateway references** -->


Service Broker is the term we use to refer to a component of the service which implements the service broker API. This component was formerly referred to as a Service Gateway, however as traffic between applications and services does not flow through the broker we found the term gateway caused confusion. Although gateway still appears in old code, we use the term broker in conversation, in new code, and in documentation.

<!-- **TODO keep paragragh** -->


Service brokers advertise a catalog of service offerings and service plans, as well as interpreting calls for provision (create), bind, unbind, and deprovision (delete). What a broker does with each call can vary between services; in general, 'provision' reserves resources on a service and 'bind' delivers information to an application necessary for accessing the resource. We call the reserved resource a Service Instance. What a service instance represents can vary by service; it could be a single database on a multi-tenant server, a dedicated cluster, or even just an account on a web application.

<image src="images/managed-services.png">

## Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL
      NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED",  "MAY", and
      "OPTIONAL" in this document are to be interpreted as described in
      [RFC 2119]( https://tools.ietf.org/html/rfc2119).

## <a id='implementation-deployment'></a>Implementation & Deployment ##

<!-- **TODO keep** -->


How a service is implemented is up to the service provider/developer. Cloud Foundry only requires that the service provider implement the service broker API. A broker can be implemented as a separate application, or by adding the required http endpoints to an existing service.

<!-- **TODO TODO delete this** -->


Because Cloud Foundry only requires that a service implements the broker API in order to be available to Cloud Foundry end users, many deployment models are possible. The following are examples of valid deployment models.

* Entire service packaged and deployed by BOSH alongside Cloud Foundry
* Broker packaged and deployed by BOSH alongside Cloud Foundry, rest of the service deployed and maintained by other means
* Broker (and optionally service) pushed as an application to Cloud Foundry user space
* Entire service, including broker, deployed and maintained outside of Cloud Foundry by other means

# API

<!-- **TODO This section is formatted poorly. The described resources
need to be listed at the top. The resource needs to be very boldly
declared at each section. I think the majority of this section would
better be served by a swagger definition. We should enforce a common
format for all sections, and move all descriptive text to a higher
level and link to it from each resource. -->


## <a id='changelog'></a>Document Changelog ##

**TODO we don't need an old changelog. We may want to save the change policy**


[v2 API Change Log](v2-api-changelog.md)

## <a id='changes'></a>Changes ##

### <a id='change-policy'></a>Change Policy ###

<!-- **TODO keep if we believe in this** -->


* Existing endpoints and fields will not be removed or renamed.
* New optional endpoints, or new HTTP methods for existing endpoints, may be
added to enable support for new features.
* New fields may be added to existing request/response messages.
These fields must be optional and should be ignored by clients and servers
that do not understand them.

### <a id='api-changes-since-v2-8'></a>Changes Since v2.8 ###

<!-- **TODO drop** -->


1. Querying `last_operation` now supports `service_id` and `plan_id` query parameters.
1. Provision, Update, Deprovision responses now accepts an optional `operation` json param for async responses. This is used to by service brokers to return an state related to the operation. Provided back to the service broker via the `last_operation` call.
1. Querying `last_operation` now supports `operation` param back to the service broker.

## <a id='dependencies'></a>Dependencies ##

<!-- **TODO delete this** -->


v2.9 of the services API has been supported since:

* [Final build 238](https://github.com/cloudfoundry/cf-release/tree/v238) of [cf-release](https://github.com/cloudfoundry/cf-release/)
* v2.57.0 of the Cloud Controller API
* CLI [v6.14.0](https://github.com/cloudfoundry/cli/releases/tag/v6.14.0)

## <a id='api-overview'></a>API Overview ##

<!-- **TODO keep** -->


The Cloud Foundry services API defines the contract between the Cloud
Controller and the service broker.
The broker is expected to implement several HTTP (or HTTPS) endpoints
underneath a URI prefix.
One or more services can be provided by a single broker, and load balancing
enables horizontal scalability of redundant brokers.
Multiple Cloud Foundry instances can be supported by a single broker using
different URL prefixes and credentials.

<image src="images/v2services-new.png" width="960" height="720" style='background-color:#fff'>

## <a id='api-version-header'></a>API Version Header ##

Requests from the Cloud Controller to the broker contain a header that defines
the version number of the Broker API that Cloud Controller will use.
This header will be useful in future minor revisions of the API to allow
brokers to reject requests from Cloud Controllers that they do not understand.
While minor API revisions will always be additive, it is possible that brokers
will come to depend on a feature that was added after 2.0, so they may use this
header to reject the request.
Error messages from the broker in this situation should inform the operator of
what the required and actual version numbers are so that an operator can go
upgrade Cloud Controller and resolve the issue.
A broker should respond with a `412 Precondition Failed` message when rejecting
a request.

The version numbers are in the format `MAJOR.MINOR`, using semantic versioning
such that 2.9 comes before 2.10.
An example of this header as of publication time is:

`X-Broker-Api-Version: 2.10`

## <a id='authentication'></a>Authentication ##

Cloud Controller (final release v145+) authenticates with the Broker using HTTP
basic authentication (the `Authorization:` header) on every request and will
reject any broker registrations that do not contain a username and password.
The broker is responsible for checking the username and password and returning
a `401 Unauthorized` message if credentials are invalid.
Cloud Controller supports connecting to a broker using SSL if additional
security is desired.

## <a id='catalog-mgmt'></a>Catalog Management ##

<!-- **TODO drop. this document never describes what 'validating a catalog' means** -->


The first endpoint that a broker must implement is the service catalog.
Cloud Controller will initially fetch this endpoint from all brokers and make
adjustments to the user-facing service catalog stored in the Cloud Controller
database.
If the catalog fails to initially load or validate, Cloud Controller will not
allow the operator to add the new broker and will give a meaningful error
message.
Cloud Controller will also update the catalog whenever a broker is updated, so
you can use `update-service-broker` with no changes to force a catalog refresh.

<!-- **TODO drop** -->


When Cloud Controller fetches a catalog from a broker, it will compare the
broker's id for services and plans with the `unique_id` values for services and
plans in the  Cloud Controller database.
If a service or plan in the broker catalog has an id that is not present
amongst the `unique_id` values in the database, a new record will be added to
the database.
If services or plans in the database are found with `unique_id`s that match the
broker catalog's id, Cloud Controller will update the records to match
the broker’s catalog.

<!-- **TODO drop** -->


If the database has plans which are not found in the broker catalog, and there
are no associated service instances, Cloud Controller will remove these plans
from the database. Cloud Controller will then delete services that do not have associated plans
from the database.
If the database has plans which are not found in the broker catalog, and there
**are** provisioned instances, the plan will be marked “inactive” and will
no longer be visible in the marketplace catalog or be provisionable.


### Request ###

#### Route ####
`GET /v2/catalog`

#### cURL ####
<pre class="terminal">
 $ curl -H "X-Broker-API-Version: 2.9" http://username:password@broker-url/v2/catalog
</pre>

### Response ###

<table border="1" class="nice">
<thead>
<tr>
  <th>Status Code</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>200 OK</td>
  <td>The expected response body is below.</td>
</tr>
</tbody>
</table>

#### Body - Schema of Service Objects ####

<!-- **TODO drop this 'string' definition thing.**  -->


CLI and web clients have different needs with regard to service and plan names.
A CLI-friendly string is all lowercase, with no spaces.
Keep it short -- imagine a user having to type it as an argument for a longer
command.
A web-friendly display name is camel-cased with spaces and punctuation supported.

<table border="1" class="nice">
<thead>
<tr>
  <th>Response field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>services*</td>
  <td>array-of-service-objects</td>
  <td>Schema of service objects defined below.</td>
</tr>
</tbody>
</table>

<h5> Service Objects </h5>

<table border="1" class="nice">
<thead>
<tr>
  <th>Response field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>&nbsp;&nbsp;&nbsp;name*</td>
  <td>string</td>
  <td>The CLI-friendly name of the service that will appear in the catalog. All lowercase, no spaces.</td>
</tr>
<tr>
  <td>&nbsp;&nbsp;&nbsp;id*</td>
  <td>string</td>
  <td>An identifier used to correlate this service in future requests to the catalog. This must be unique within Cloud Foundry, using a GUID is recommended. </td>
</tr>
<tr>
  <td>&nbsp;&nbsp;&nbsp;description*</td>
  <td>string</td>
  <td>A short description of the service that will appear in the catalog.</td>
</tr>
<tr>
  <td>&nbsp;&nbsp;&nbsp;tags</td>
  <td>array-of-strings</td>
  <td>Tags provide a flexible mechanism to expose a classification, attribute, or base technology of a service, enabling equivalent services to be swapped out without changes to dependent logic in applications, buildpacks, or other services. E.g. mysql, relational, redis, key-value, caching, messaging, amqp.</td>
</tr>
<tr>
  <td>&nbsp;&nbsp;&nbsp;requires</td>
  <td>array-of-strings</td>
  <td>A list of permissions that the user would have to give the service, if they provision it. The only permissions currently supported are <tt>syslog_drain</tt>, <tt>route_forwarding</tt> and <tt>volume_mount</tt>; for more info see <a href="app-log-streaming.md">Application Log Streaming</a>, <a href="route-services.md">Route Services</a> and <a href="volume-services.md">Volume Services</a>.</td>
<tr>
  <td>&nbsp;&nbsp;&nbsp;max_db_per_node:</td>
  <td>strings</td>
  <td></td>
<tr>
  <td>&nbsp;&nbsp;&nbsp;bindable*</td>
  <td>boolean</td>
  <td>Whether the service can be bound to applications.</td>
</tr>
<tr>
  <td>&nbsp;&nbsp;&nbsp;metadata</td>
  <td>object</td>
  <td>A list of metadata for a service offering. For more information, see <a href="catalog-metadata.md">Service Metadata</a>.</td>
</tr>
<tr>
  <td><a href="#DObject">&nbsp;&nbsp;&nbsp;dashboard_client</a></td>
  <td>object</td>
  <td>Contains the data necessary to activate the <a href="dashboard-sso.md">Dashboard SSO feature</a> for this service</td>
</tr>
<tr>
  <td>&nbsp;&nbsp;&nbsp;plan_updateable</td>
  <td>boolean</td>
  <td>
    Whether the service supports upgrade/downgrade for some plans.
    <br/>
    Please note that the misspelling of the attribute <i>plan_updatable</i> to <i>plan_updateable</i> was done by mistake. We have opted to keep that misspelling instead of fixing it and thus breaking backward compatibility.
  </td>
</tr>
<tr>
  <td><a href="#PObject">&nbsp;&nbsp;&nbsp;plans*</a></td>
  <td>array-of-objects</td>
  <td>A list of plans for this service, schema is defined below.</td>
</tr>
</tbody>
</table>

<h5> Dashboard Client Object <a name="DObject"></a> </h5>

<table border="1" class="nice">
<thead>
<tr>
  <th>Response field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>&nbsp;&nbsp;&nbsp;id</td>
  <td>string</td>
  <td>The id of the Oauth2 client that the service intends to use. The name may be taken, in which case the updat will return an error to the operator</td>
</tr>
<tr>
  <td>&nbsp;&nbsp;&nbsp;secret</td>
  <td>string</td>
  <td>A secret for the dashboard client</td>
</tr>
<tr>
  <td>&nbsp;&nbsp;&nbsp;redirect_uri</td>
  <td>string</td>
  <td>A domain for the service dashboard that will be whitelisted by the UAA to enable SSO</td>
</tr>
</tbody>
</table>


<h5> Plan Object <a name="PObject"></a> </h5>

<table border="1" class="nice">
<thead>
<tr>
  <th>Response field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>&nbsp;&nbsp;&nbsp;id*</td>
  <td>string</td>
  <td>An identifier used to correlate this plan in future requests to the catalog. This must be unique within Cloud Foundry, using a GUID is recommended.</td>
</tr>
<tr>
  <td>&nbsp;&nbsp;&nbsp;name*</td>
  <td>string</td>
  <td>The CLI-friendly name of the plan that will appear in the catalog. All lowercase, no spaces.</td>
</tr>
<tr>
  <td>&nbsp;&nbsp;&nbsp;description*</td>
  <td>string</td>
  <td>A short description of the service that will appear in the catalog.</td>
</tr>
<tr>
  <td>&nbsp;&nbsp;&nbsp;metadata</td>
  <td>object</td>
  <td>A list of metadata for a service plan. For more information, see <a href="catalog-metadata.md">Service Metadata</a>.</td>
</tr>
<tr>
  <td>&nbsp;&nbsp;&nbsp;free</td>
  <td>boolean</td>
  <td>This field allows the plan to be limited by the non_basic_services_allowed field in a Cloud Foundry Quota, see <a href="http://docs.cloudfoundry.org/running/managing-cf/quota-plans.md">Quota Plans</a>. Default: true</td>
</tr>
</tbody>
</table>

\* Fields with an asterisk are required.

<pre>
{
  "services": [{
    "name": "fake-service",
    "id": "acb56d7c-XXXX-XXXX-XXXX-feb140a59a66",
    "description": "fake service",
    "tags": ["no-sql", "relational"],
    "requires": ["route_forwarding"],
    "max_db_per_node": 5,
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
      "displayName": "The Fake Broker"
    },
    "dashboard_client": {
      "id": "398e2f8e-XXXX-XXXX-XXXX-19a71ecbcf64",
      "secret": "277cabb0-XXXX-XXXX-XXXX-7822c0a90e5d",
      "redirect_uri": "http://localhost:1234"
    },
    "plan_updateable": true,
    "plans": [{
      "name": "fake-plan",
      "id": "d3031751-XXXX-XXXX-XXXX-a42377d3320e",
      "description": "Shared fake Server, 5tb persistent disk, 40 max concurrent connections",
      "max_storage_tb": 5,
      "metadata": {
        "cost": 0,
        "bullets": [{
          "content": "Shared fake server"
        }, {
          "content": "5 TB storage"
        }, {
          "content": "40 concurrent connections"
        }]
      }
    }, {
      "name": "fake-async-plan",
      "id": "0f4008b5-XXXX-XXXX-XXXX-dace631cd648",
      "description": "Shared fake Server, 5tb persistent disk, 40 max concurrent connections. 100 async",
      "max_storage_tb": 5,
      "metadata": {
        "cost": 0,
        "bullets": [{
          "content": "40 concurrent connections"
        }]
      }
    }, {
      "name": "fake-async-only-plan",
      "id": "8d415f6a-XXXX-XXXX-XXXX-e61f3baa1c77",
      "description": "Shared fake Server, 5tb persistent disk, 40 max concurrent connections. 100 async",
      "max_storage_tb": 5,
      "metadata": {
        "cost": 0,
        "bullets": [{
          "content": "40 concurrent connections"
        }]
      }
    }]
  }]
}
</pre>


### <a id='create-broker'></a>Adding a Broker to Cloud Foundry ###

Once you've implemented the first endpoint `GET /v2/catalog` above, you'll want
to [register the broker with CF](managing-service-brokers.md#register-broker)
to make your services and plans available to end users.

## <a id='asynchronous-operations'></a>Asynchronous Operations ##

<!-- **TODO this should be a separate higher-level section and not shoved inside a resource definition**  -->


Previously, Cloud Foundry only supported synchronous integration with service brokers. Brokers must return a valid response within 60 seconds and if the response is `201 CREATED`, users expect a service instance to be usable. This limits the services brokers can offer to those that can be provisioned in 60 seconds; brokers could return a success prematurely, but this leaves users wondering why their service instance is not usable and when it will be.

With support for Asynchronous Operations, brokers still must respond within 60 seconds but may now return a `202 ACCEPTED`, indicating that the requested operation has been accepted but is not complete. This triggers Cloud Foundry to poll a new endpoint `/v2/service_instances/:guid/last_operation` until the broker indicates that the requested operation has succeeded or failed. During the intervening time, end users are able to discover the state of the requested operation using Cloud Foundry API clients such as the CLI.

For an operation to be executed asynchronously, all three components (CF API client, CF, and broker) must support the feature. The parameter `accepts_incomplete=true` must be passed in a request by the CF API client, triggering CF to include the same parameter in a request to the broker. The broker can then choose to execute the request synchronously or asynchronously.

If the broker executes the request asynchronously, the response must use the status code `202 ACCEPTED`; the response body should be the same as if the broker were serving the request synchronously.

<p class='note'><strong>Note:</strong> Asynchronous Operations are currently supported only for provision, update, and deprovision. Bind and unbind will be added once the feature is considered stable.</p>

If the `accepts_incomplete=true` parameter is not included, and the broker cannot fulfill the request synchronously (guaranteeing that the operation is complete on response), then the broker should reject the request with the status code `422 UNPROCESSABLE ENTITY` and the following body:

<pre class="terminal">
{
  "error": "AsyncRequired",
  "description": "This service plan requires client support for asynchronous service operations."
}
</pre>

To execute a request synchronously, the broker need only return the usual status codes; `201 CREATED` for create, and `200 OK` for update and delete.

### <a id='sequence-diagram'></a>Sequence Diagram ###
<a href='images/async-service-broker-flow.png' target='_blank'>
  <image src="images/async-service-broker-flow.png" width="1250" height="823" style='background-color:#fff'>
</a>

### <a id='blocking'></a>Blocking Operations ###

The Cloud Controller ensures that service brokers do not receive requests for an instance while an asynchronous operation is in progress. For example, if a broker is in the process of provisioning an instance asynchronously, the Cloud Controller will not allow any update, bind, unbind, or deprovision requests to be made through the platform. A user who attempts to perform one of these actions while an operation is already in progress will get an HTTP 400 with error message "Another operation for this service instance is in progress."

### <a id='when-to-use-async'></a>When to use Asynchronous Service Operations ###

Service brokers should respond to all Cloud Controller requests within 60 seconds. Brokers that can guarantee completion of the requested operation with the response may return the synchronous response (e.g. `201 CREATED` for a provision request). Brokers that cannot guarantee completion of the operation with the response should implement support for asynchronous provisioning. Support for synchronous or asynchronous responses may vary by service offering, even by service plan.

## <a id='polling'></a>Polling Last Operation (async only) ##

<!-- **TODO move up to higher-level section**  -->


When a broker returns status code `202 ACCEPTED` for [provision](#provisioning), [update](#updating_service_instance), or [deprovision](#deprovisioning), Cloud Foundry will begin to poll the `/v2/service_instances/:guid/last_operation` endpoint to obtain the state of the last requested operation. The broker response must contain the field `state` and an optional field `description`.

Valid values for `state` are `in progress`, `succeeded`, and `failed`. Cloud Foundry will poll the `last_operation` endpoint as long as the broker returns `"state": "in progress"`. Returning `"state": "succeeded"` or `"state": "failed"` will cause Cloud Foundry to cease polling. The value provided for `description` will be passed through to the CF API client and can be used to provide additional detail for users about the state of the operation.

### Request ###

##### Route #####
`GET /v2/service_instances/:instance_id/last_operation`

##### Parameters #####

The request provides these query string parameters as useful hints for brokers.

<table border="1" class="nice">
<thead>
<tr>
  <th>Query-String Field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>service_id</td>
  <td>string</td>
  <td>ID of the service from the catalog.</td>
</tr>
<tr>
  <td>plan_id</td>
  <td>string</td>
  <td>ID of the plan from the catalog.</td>
</tr>
<tr>
  <td>operation</td>
  <td>string</td>
  <td>The field optionally returned by the service broker on Provision, Update, Deprovision async responses. Represents any state the service broker responsed with as a URL encoded string.</a>.</td>
</tr>
</tbody>
</table>

<p class="note"><strong>Note:</strong> Although the request query parameters <code>service_id</code> and <code>plan_id</code> are not required, Cloud Controller includes them on all <code>last_operation</code> requests it makes to service brokers.</p>

##### cURL #####
<pre class="terminal">
$ curl http://username:password@broker-url/v2/service_instances/:instance_id/last_operation
</pre>

### Response ###

<table border="1" class="nice">
<thead>
<tr>
  <th>Status Code</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>200 OK</td>
  <td>The expected response body is below.</td>
</tr>
<tr>
  <td>410 GONE</td>
  <td>Appropriate only for asynchronous delete requests. Cloud Foundry will consider this response a success and remove the resource from its database. The expected response body is <code>{}</code>. Returning this while Cloud Foundry is polling for create or update operations will be interpreted as an invalid response and Cloud Foundry will continue polling.</td>
</tr>
</tbody>
</table>

<!-- **TODO this response definition needs to be moved up to a top level section and referenced.**  -->


Responses with any other status code will be interpreted as an error or invalid response; Cloud Foundry will continue polling until the broker returns a valid response or the [maximum polling duration](#max-polling-duration) is reached. Brokers may use the `description` field to expose user-facing error messages about the operation state; for more info see [Broker Errors](api.md#broker-errors).

##### Body #####

For success responses, the following fields are valid.

<table border="1" class="nice">
<thead>
<tr>
  <th>Response field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>state*</td>
  <td>string</td>
  <td>Valid values are <code>in progress</code>, <code>succeeded</code>, and <code>failed</code>. While <code>"state": "in progress"</code>, Cloud Foundry will continue polling. A response with <code>"state": "succeeded"</code> or <code>"state": "failed"</code> will cause Cloud Foundry to cease polling.</td>
</tr>
<tr>
  <td>description</td>
  <td>string</td>
  <td>Optional field. A user-facing message displayed to the Cloud Foundry API client. Can be used to tell the user details about the status of the operation.</td>
</tr>
</tboby>
</table>

\* Fields with an asterisk are required.

<pre class="terminal">
{
  "state": "in progress",
  "description": "Creating service (10% complete)."
}
</pre>

### <a id='polling-interval'></a> Polling Interval ###
When a broker responds asynchronously to a request from Cloud Foundry containing the `accepts_incomplete=true` parameter, Cloud Foundry will poll the broker for the operation state at a configured interval. The Cloud Foundry operator can configure this interval in the BOSH deployment manifest using the property `properties.cc.broker_client_default_async_poll_interval_seconds` (defaults to 60 seconds). The maximum supported polling interval is 86400 seconds (24 hours).

### <a id='max-polling-duration'></a>Maximum Polling Duration ###
When a broker responds asynchronously to a request from Cloud Foundry containing the `accepts_incomplete=true` parameter, Cloud Foundry will poll the broker for the operation state until the broker response includes `"state":"succeeded"` or `"state":"failed"`, or until a maximum polling duration is reached. If the max polling duration is reached, Cloud Foundry will cease polling and the operation state will be considered `failed`. The Cloud Foundry operator can configure this max polling duration in the BOSH deployment manifest using the property `properties.cc.broker_client_max_async_poll_duration_minutes` (defaults to 10080 minutes or 1 week).

### <a id='additional-resources'></a>Additional Resources ###

* An example broker that implements this feature can be found at [Example Service Brokers](examples.md).
* A demo video of the CLI user experience using the above broker can be found [here](https://youtu.be/Ij5KSKrAq9Q).

## <a id='provisioning'></a>Provisioning ##

When the broker receives a provision request from Cloud Controller, it should
synchronously take whatever action is necessary to create a new service
resource for the developer.
The result of provisioning varies by service type, although there are a few
common actions that work for many services.
For a MySQL service, provisioning could result in:

* An empty dedicated `mysqld` process running on its own VM.
* An empty dedicated `mysqld` process running in a lightweight container on a
shared VM.
* An empty dedicated `mysqld` process running on a shared VM.
* An empty dedicated database, on an existing shared running `mysqld`.
* A database with business schema already there.
* A copy of a full database, for example a QA database that is a copy of the
production database.

For non-data services, provisioning could just mean getting an account on an
existing system.

### Request ###

##### Route #####
`PUT /v2/service_instances/:instance_id`

<p class="note"><strong>Note</strong>: The <code>:instance_id</code> of a service instance is provided by the Cloud Controller. This ID will be used for future requests (bind and deprovision), so the broker must use it to correlate the resource it creates.</p>

##### Body #####

<table border="1" class="nice">
<thead>
<tr>
  <th>Request field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>organization_guid*</td>
  <td>string</td>
  <td>The Cloud Controller GUID of the organization under which the service is to be provisioned. Although most brokers will not use this field, it could be helpful in determining data placement or applying custom business rules.</td>
</tr>
<tr>
  <td>plan_id*</td>
  <td>string</td>
  <td>The ID of the plan within the above service (from the catalog endpoint) that the user would like provisioned. Because plans have identifiers unique to a broker, this is enough information to determine what to provision.</td>
</tr>
<tr>
  <td>service_id*</td>
  <td>string</td>
  <td>The ID of the service within the catalog above.</td>
</tr>
<tr>
  <td>space_guid*</td>
  <td>string</td>
  <td>Similar to organization_guid, but for the space.</td>
</tr>
<tr>
  <td>parameters</td>
  <td>JSON object</td>
  <td>Cloud Foundry API clients can provide a JSON object of configuration parameters with their request and this value will be passed through to the service broker. Brokers are responsible for validation.</td>
</tr>
<tr>
  <td>accepts_incomplete</td>
  <td>boolean</td>
  <td>A value of true indicates that both the Cloud Controller and the requesting client support asynchronous provisioning. If this parameter is not included in the request, and the broker can only provision an instance of the requested plan asynchronously, the broker should reject the request with a 422 as described below.</td>
</tr>
</tbody>
</table>

\* Fields with an asterisk are required.

<pre class="terminal">
{
  "organization_guid": "org-guid-here",
  "plan_id":           "plan-guid-here",
  "service_id":        "service-guid-here",
  "space_guid":        "space-guid-here",
  "parameters":        {
    "parameter1": 1,
    "parameter2": "value"
  }
}
</pre>

##### cURL #####
<pre class="terminal">
$ curl http://username:password@broker-url/v2/service_instances/:instance_id -d '{
  "organization_guid": "org-guid-here",
  "plan_id":           "plan-guid-here",
  "service_id":        "service-guid-here",
  "space_guid":        "space-guid-here",
  "parameters":        {
    "parameter1": 1,
    "parameter2": "value"
  }
}' -X PUT -H "X-Broker-API-Version: 2.9" -H "Content-Type: application/json"
</pre>

In this case, `instance_id` refers to the service instance id generated by Cloud
Controller

### Response ###

<table border="1" class="nice">
<thead>
<tr>
  <th>Status Code</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>201 Created</td>
  <td>Service instance has been created. The expected response body is below.</td>
</tr>
<tr>
  <td>200 OK</td>
  <td>
    May be returned if the service instance already exists and the requested parameters are identical to the existing service instance.
    The expected response body is below.
  </td>
</tr>
<tr>
  <td>202 Accepted</td>
  <td>Service instance creation is in progress. This triggers Cloud Controller to poll the <a href="#polling">Service Instance Last Operation Endpoint</a> for operation status.</td>
</tr>
<tr>
  <td>409 Conflict</td>
  <td>Should be returned if the requested service instance already exists. The expected response body is <code>{}</code>.</td>
</tr>
<tr>
  <td>422 Unprocessable Entity</td>
  <td>Should be returned if the broker only supports asynchronous provisioning for the requested plan and the request did not include <code>?accepts_incomplete=true</code>. The expected response body is: <code>{ "error": "AsyncRequired", "description": "This service plan requires client support for asynchronous service operations." }</code>, as described below.</td>
</tr>
</tbody>
</table>

Responses with any other status code will be interpreted as a failure. Brokers can include a user-facing message in the `description` field; for details see [Broker Errors](#broker-errors).

##### Body #####

For success responses, the following fields are supported. Others will be ignored. For error responses, see [Broker Errors](#broker-errors).

<table border="1" class="nice">
<thead>
<tr>
  <th>Response field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>dashboard_url</td>
  <td>string</td>
  <td>The URL of a web-based management user interface for the service instance; we refer to this as a service dashboard. The URL should contain enough information for the dashboard to identify the resource being accessed ("9189kdfsk0vfnku" in the example below). For information on how users can authenticate with service dashboards via SSO, see <a href="dashboard-sso.md">Dashboard Single Sign-On</a>.</td>
</tr>
<tr>
  <td>operation</td>
  <td>string</td>
  <td>For async responses, service brokers can return operation state as a string. This field will be provided back to the service broker on <code>last_operation</code> requests as a URL encoded query param.</a>.</td>
</tr>
</tboby>
</table>
\* Fields with an asterisk are required.
<pre class="terminal">
{
 "dashboard_url": "<span>http</span>://example-dashboard.example.com/9189kdfsk0vfnku",
 "operation": "task_10"
}
</pre>


## <a id='updating_service_instance'></a>Updating a Service Instance ##

By implementing this endpoint, broker authors can enable users to modify two attributes of an existing service instance; the service plan and parameters. By changing the service plan, users can upgrade or downgrade their service instance to other plans. By modifying properties, users can change configuration options that are specific to a service or plan. To see how Cloud Foundry users make these requests, see [Managing Services](../devguide/services/managing-services.html#update_service).

To enable this functionality, a broker declares support for each service by including `plan_updateable: true` in its [catalog endpoint](#catalog-mgmt). If this optional field is not included, Cloud Foundry will return a meaningful error to users for any plan change request, and will not make an API call to the broker. If this field is included and configured as true, Cloud Foundry will make API calls to the broker for all plan change requests, and it is up to the broker to validate whether a particular permutation of plan change is supported. Not all permutations of plan changes are expected to be supported. For example, a service may support upgrading from plan "shared small" to "shared large" but not to plan "dedicated". If a particular plan change is not supported, the broker should return a meaningful error message in response.

### Request ###

##### Route #####
`PATCH /v2/service_instances/:instance_id`

<p class="note"><strong>Note</strong>: <code>:instance_id</code> is the global unique ID of a previously-provisioned service instance.</p>

##### Body #####

<table border="1" class="nice">
<thead>
<tr>
  <th>Request Field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>service_id*</td>
  <td>string</td>
  <td>The ID of the service within the catalog above.</td>
</tr>
<tr>
  <td>plan_id</td>
  <td>string</td>
  <td>ID of the new plan from the catalog.</td>
</tr>
<tr>
  <td>parameters</td>
  <td>JSON object</td>
  <td>Cloud Foundry API clients can provide a JSON object of configuration parameters with their request and this value will be passed through to the service broker. Brokers are responsible for validation.</td>
</tr>
<tr>
  <td>previous_values</td>
  <td>object</td>
  <td>Information about the instance prior to the update.</td>
</tr>
<tr>
  <td>previous_values.plan_id</td>
  <td>string</td>
  <td>ID of the plan prior to the update.</td>
</tr>
<tr>
  <td>previous_values.service_id</td>
  <td>string</td>
  <td>ID of the service for the instance.</td>
</tr>
<tr>
  <td>previous_values.organization_id</td>
  <td>string</td>
  <td>ID of the organization containing the instance.</td>
</tr>
<tr>
  <td>previous_values.space_id</td>
  <td>string</td>
  <td>ID of the space containing the instance..</td>
</tr>
<tr>
  <td>accepts_incomplete</td>
  <td>boolean</td>
  <td>A value of true indicates that both the Cloud Controller and the requesting client support asynchronous update. If this parameter is not included in the request, and the broker can only update an instance of the requested plan asynchronously, the broker should reject the request with a 422 as described below.</td>
</tr>
</tbody>
</table>

\* Fields with an asterisk are required.

<pre class="terminal">
{
  "service_id": "service-guid-here",
  "plan_id": "plan-guid-here",
  "parameters":        {
    "parameter1": 1,
    "parameter2": "value"
  },
  "previous_values": {
    "plan_id": "old-plan-guid-here",
    "service_id": "service-guid-here",
    "organization_id": "org-guid-here",
    "space_id": "space-guid-here"
  }
}
</pre>

##### cURL #####
<pre class="terminal">
$ curl http://username:password@broker-url/v2/service_instances/:instance_id -d '{
  "service_id": "service-guid-here",
  "plan_id": "plan-guid-here",
  "parameters":        {
    "parameter1": 1,
    "parameter2": "value"
  },
  "previous_values": {
    "plan_id": "old-plan-guid-here",
    "service_id": "service-guid-here",
    "organization_id": "org-guid-here",
    "space_id": "space-guid-here"
  }
}' -X PATCH -H "X-Broker-API-Version: 2.9" -H "Content-Type: application/json"
</pre>

### Response ###

<table border="1" class="nice">
<thead>
<tr>
  <th>Status Code</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>200 OK</td>
  <td>New plan is effective. The expected response body is <code>{}</code>.</td>
</tr>
<tr>
  <td>202 Accepted</td>
  <td>Service instance update is in progress. This triggers Cloud Controller to poll the <a href="#polling">Service Instance Last Operation Endpoint</a> for operation status.</td>
</tr>
<tr>
  <td>422 Unprocessable entity</td>
  <td>
    May be returned if the particular plan change requested is not supported or if the request cannot currently be fulfilled due to the state of the instance (eg. instance utilization is over the quota of the requested plan). Broker should include a user-facing message in the body; for details see <a href="#broker-errors">Broker Errors</a>.  Additionally, a <code>422</code> can also be returned if the broker only supports asynchronous update for the requested plan and the request did not include <code>?accepts_incomplete=true</code>. The expected response body is: <code>{ "error": "AsyncRequired", "description": "This service plan requires client support for asynchronous service operations." }</code></a>.
  </td>
</tr>
</tbody>
</table>

Responses with any other status code will be interpreted as a failure. Brokers can include a user-facing message in the `description` field; for details see [Broker Errors](#broker-errors).

##### Body #####

For success responses, the following fields are supported. Others will be ignored. For error responses, see [Broker Errors](#broker-errors).

<table border="1" class="nice">
<thead>
<tr>
  <th>Response field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>operation</td>
  <td>string</td>
  <td>For async responses, service brokers can return operation state as a string. This field will be provided back to the service broker on `last_operation` requests as a URL encoded query param.</a>.</td>
</tr>
</tboby>
</table>
\* Fields with an asterisk are required.
<pre class="terminal">
{
 "operation": "task_10"
}
</pre>


## <a id='binding'></a>Binding ##

<!-- **TODO this note should not be hidden down in the binding definition.**  -->

<p class="note"><strong>Note</strong>: Not all services must be bindable --- some deliver value just from being provisioned. Brokers that offer services that are bindable should declare them as such using <code>bindable: true</code> in the <a href="#catalog-mgmt">Catalog</a>. Brokers that do not offer any bindable services do not need to implement the endpoint for bind requests.</p>

### <a id='binding-types'></a>Types of Binding ###

#### Credentials ####

Credentials are a set of information used by an application or a user to utilize the service instance. If `bindable:true` is declared for a service in the catalog endpoint, users may request generation of credentials either by binding the service instance to an application or by creating a service key. When a service instance is bound to an app, Cloud Foundry will send the app id with the request. When a service key is created, the app id is not included. If the broker supports generation of credentials it should return `credentials` in the response. Credentials should be unique whenever possible, so access can be revoked for one application or user without affecting another. For more information on credentials, see [Binding Credentials](binding-credentials.md).

#### <a id='binding-syslog-drain'></a>Application Log Streaming ####

In response to a bind request for an application (`app_id` included), a broker may also enable streaming of application logs from Cloud Foundry to a consuming service instance by returning `syslog_drain_url`. For details, see [Application Log Streaming](app-log-streaming.md).

#### <a id='binding-route-services'></a>Route Services ####

<!-- **TODO drop**  -->


If a broker has declared `"requires":["route_forwarding"]` for a service in the Catalog endpoint, Cloud Foundry will permit a user to bind a service to a route. When bound to a route, the route itself will be sent with the bind request. A route is an address used by clients to reach apps mapped to the route. In response a broker may return a `route_service_url` which Cloud Foundry will use to proxy any request for the route to the service instance at URL specified by `route_service_url`. A broker may declare `"requires":["route_forwarding"]` but not return `route_service_url`; this enables a broker to dynamically configure a network component already in the request path for the route, requiring no change in the Cloud Foundry router. For more information, see [Route Services](route-services.md).

#### <a id='binding-volume-services'></a>Volume Services (Experimental)####

<!-- **TODO maybe move volume support to a higher level and provide a link here** -->


If a broker has declared `"requires":["volume_mount"]` for a service in the Catalog endpoint, Cloud Foundry will permit a user to bind one or more volumes to an application.  In response to a bind request a volume service broker should return a set of `volume_mount` instructions that Cloud Foundry will ensure are mounted into the application's containers.  For more information, see [Volume Services](volume-services.md)

### Request ###

##### Route #####
`PUT /v2/service_instances/:instance_id/service_bindings/:binding_id`

<p class="note"><strong>Note</strong>: The <code>:binding_id</code> of a service binding is provided by the Cloud Controller.
<code>:instance_id</code> is the ID of a previously-provisioned service instance; <code>:binding_id</code>
will be used for future unbind requests, so the broker must use it to correlate
the resource it creates.</p>

##### Body #####

<table border="1" class="nice">
<thead>
<tr>
  <th>Request Field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>service_id*</td>
  <td>string</td>
  <td>ID of the service from the catalog.</td>
</tr>
<tr>
  <td>plan_id*</td>
  <td>string</td>
  <td>ID of the plan from the catalog.</td>
</tr>
<tr>
  <td>app_guid</td>
  <td>string</td>
  <td>GUID of the application that you want to bind your service to. Will be included when users bind applications to service instances.</td>
</tr>
<tr>
  <td>bind_resource</td>
  <td>JSON object</td>
  <td>A JSON object that contains the required fields of the resource being bound. Currently only <code>app_guid</code> for application bindings and <code>route</code> for route bindings are supported.</td>
</tr>
<tr>
  <td>parameters</td>
  <td>JSON object</td>
  <td>Cloud Foundry API clients can provide a JSON object of configuration parameters with their request and this value will be passed through to the service broker. Brokers are responsible for validation.</td>
</tr>
</tbody>
</table>

\* Fields with an asterisk are required.

<pre class="terminal">
{
  "plan_id":      "plan-guid-here",
  "service_id":   "service-guid-here",
  "app_guid":     "app-guid-here",
  "bind_resource":     {
    "app_guid": "app-guid-here"
  },
  "parameters":        {
    "parameter1-name-here": 1,
    "parameter2-name-here": "parameter2-value-here"
  }
}
</pre>

<pre class="terminal">
{
  "plan_id":      "plan-guid-here",
  "service_id":   "service-guid-here",
  "bind_resource":     {
    "route": "route-url-here"
  },
  "parameters":        {
    "parameter1-name-here": 1,
    "parameter2-name-here": "parameter2-value-here"
  }
}
</pre>

##### cURL #####
<pre class="terminal">
$ curl http://username:password@broker-url/v2/service_instances/:instance_id/service_bindings/:binding_id -d '{
  "plan_id":      "plan-guid-here",
  "service_id":   "service-guid-here",
  "app_guid":     "app-guid-here",
  "bind_resource":     {
    "app_guid": "app-guid-here"
  },
  "parameters":        {
    "parameter1-name-here": 1,
    "parameter2-name-here": "parameter2-value-here"
  }
}' -X PUT
</pre>

<pre class="terminal">
$ curl http://username:password@broker-url/v2/service_instances/:instance_id/service_bindings/:binding_id -d '{
  "plan_id":      "plan-guid-here",
  "service_id":   "service-guid-here",
  "bind_resource":     {
    "route": "route-url-here"
  },
  "parameters":        {
    "parameter1-name-here": 1,
    "parameter2-name-here": "parameter2-value-here"
  }
}' -X PUT
</pre>

In this case, `instance_id` refers to the id of an existing service instance in a previous provisioning, while `binding_id` is service binding id generated by Cloud Controller.

### Response ###

<table border="1" class="nice">
<thead>
<tr>
  <th>Status Code</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>201 Created</td>
  <td>Binding has been created. The expected response body is below.</td>
</tr>
<tr>
  <td>200 OK</td>
  <td>
    May be returned if the binding already exists and the requested parameters are identical to the existing binding.
    The expected response body is below.
  </td>
</tr>
<tr>
  <td>409 Conflict</td>
  <td>Should be returned if the requested binding already exists. The expected response body is <code>{}</code>, though the description field can be used to return a user-facing error message, as described in <a href="#broker-errors">Broker Errors</a>.</td>
</tr>
<tr>
  <td>422 Unprocessable Entity</td>
  <td>Should be returned if the broker requires that <code>app_guid</code> be included in the request body. The expected response body is: <code>{ "error": "RequiresApp", "description": "This service supports generation of credentials through binding an application only." }</code></td>
</tr>
</tbody>
</table>

Responses with any other status code will be interpreted as a failure and an unbind request will be sent to the broker to prevent an orphan being created on the broker. Brokers can include a user-facing message in the `description` field; for details see [Broker Errors](#broker-errors).

##### Body #####

For success responses, the following fields are supported. Others will be ignored. For error responses, see [Broker Errors](#broker-errors).

<table border="1" class="nice">
<thead>
<tr>
  <th>Response Field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>credentials</td>
  <td>object</td>
  <td>A free-form hash of credentials that the bound application can use to access the service. For more information, see <a href="binding-credentials.md">Binding Credentials</a>.</td>
</tr>
<tr>
  <td>syslog_drain_url</td>
  <td>string</td>
  <td>A URL to which Cloud Foundry should drain logs for the bound application. <code>requires:syslog_drain</code> must be declared in the <a href="#catalog-mgmt">catalog endpoint</a> or Cloud Foundry will consider the response invalid. For details, see <a href="app-log-streaming.md">Application Log Streaming</a>.</td>
</tr>
<tr>
  <td>route_service_url</td>
  <td>string</td>
  <td>A URL to which Cloud Foundry should proxy requests for the bound route. <code>requires:route_forwarding</code> must be declared in the <a href="#catalog-mgmt">catalog endpoint</a> or Cloud Foundry will consider the response invalid. For details, see <a href="route-services.md">Route Services</a>.</td>
</tr>
<tr>
  <td>volume_mounts</td>
  <td>array-of-objects</td>
  <td>An array of volume mount instructions.  <code>requires:volume_mount</code> must be declared in the <a href="#catalog-mgmt">catalog endpoint</a> or Cloud Foundry will consider the response invalid.  For more information, see <a href="volume-services.md">Volume Services</a></td>
</tr>
</tbody>
</table>

<pre class="terminal">
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
</pre>

## <a id='unbinding'></a>Unbinding ##

<p class="note"><strong>Note</strong>: Brokers that do not provide any bindable services do not need to implement the endpoint for unbind requests.</p>

When a broker receives an unbind request from Cloud Controller, it should
delete any resources it created in bind.
Usually this means that an application immediately cannot access the resource.

### Request ###




##### Route #####
`DELETE /v2/service_instances/:instance_id/service_bindings/:binding_id`

The `:binding_id` in the URL is the identifier of a previously created binding (the same `:binding_id` passed in the bind request). The request has no body, because DELETE requests generally do not have bodies.

##### Parameters #####

The request provides these query string parameters as useful hints for brokers.

<table border="1" class="nice">
<thead>
<tr>
  <th>Query-String Field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>service_id*</td>
  <td>string</td>
  <td>ID of the service from the catalog.</td>
</tr>
<tr>
  <td>plan_id*</td>
  <td>string</td>
  <td>ID of the plan from the catalog.</td>
</tr>
</tbody>
</table>

\* Query parameters with an asterisk are required.

##### cURL #####
<pre class="terminal">
$ curl 'http://username:password@broker-url/v2/service_instances/:instance_id/
  service_bindings/:binding_id?service_id=service-id-here&plan_id=plan-id-here' -X DELETE -H "X-Broker-API-Version: 2.9"
</pre>

### Response ###

<table border="1" class="nice">
<thead>
<tr>
  <th>Status Code</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>200 OK</td>
  <td>Binding was deleted. The expected response body is <code>{}</code>.</td>
</tr>
<tr>
  <td>410 Gone</td>
  <td>Should be returned if the binding does not exist. The expected response body is <code>{}</code>.</td>
</tr>
</tbody>
</table>

Responses with any other status code will be interpreted as a failure and the binding will remain in the Cloud Controller database. Brokers can include a user-facing message in the `description` field; for details see [Broker Errors](#broker-errors).

##### Body #####

For a success response, the expected response body is `{}`.

## <a id='deprovisioning'></a>Deprovisioning ##

When a broker receives a deprovision request from Cloud Controller, it should
delete any resources it created during the provision.
Usually this means that all resources are immediately reclaimed for future
provisions.

### Request ###

##### Route #####
`DELETE /v2/service_instances/:instance_id`

The `:instance_id` in the URL is the identifier of a previously provisioned instance (the same
`:instance_id` passed in the provision request).  The request has no body, because DELETE
requests generally do not have bodies.

##### Parameters #####

The request provides these query string parameters as useful hints for brokers.

<table border="1" class="nice">
<thead>
<tr>
  <th>Query-String Field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>service_id*</td>
  <td>string</td>
  <td>ID of the service from the catalog.</td>
</tr>
<tr>
  <td>plan_id*</td>
  <td>string</td>
  <td>ID of the plan from the catalog.</td>
</tr>
<tr>
  <td>accepts_incomplete</td>
  <td>boolean</td>
  <td>A value of true indicates that both the Cloud Controller and the requesting client support asynchronous deprovisioning. If this parameter is not included in the request, and the broker can only deprovision an instance of the requested plan asynchronously, the broker should reject the request with a 422 as described below.</td>
</tr>
</tbody>
</table>

\* Query parameters with an asterisk are required.

##### cURL #####
<pre class="terminal">
$ curl 'http://username:password@broker-url/v2/service_instances/:instance_id?service_id=
    service-id-here&plan_id=plan-id-here' -X DELETE -H "X-Broker-API-Version: 2.9"
</pre>

### Response ###

<table border="1" class="nice">
<thead>
<tr>
  <th>Status Code</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>200 OK</td>
  <td>Service instance was deleted. The expected response body is <code>{}</code>.</td>
</tr>
<tr>
  <td>202 Accepted</td>
  <td>Service instance deletion is in progress. This triggers Cloud Controller to poll the <a href="#polling">Service Instance Last Operation Endpoint</a> for operation status.</td>
</tr>
<tr>
  <td>410 Gone</td>
  <td>Should be returned if the service instance does not exist. The expected response body is <code>{}</code>.</td>
</tr>
<tr>
  <td>422 Unprocessable Entity</td>
  <td>Should be returned if the broker only supports asynchronous deprovisioning for the requested plan and the request did not include <code>?accepts_incomplete=true</code>. The expected response body is: <code>{ "error": "AsyncRequired", "description": "This service plan requires client support for asynchronous service operations." }</code>, as described below.</td>
</tr>
</tbody>
</table>

Responses with any other status code will be interpreted as a failure and the service instance will remain in the Cloud Controller database. Brokers can include a user-facing message in the `description` field; for details see [Broker Errors](#broker-errors).

##### Body #####

For success responses, the following fields are supported. Others will be ignored. For error responses, see [Broker Errors](#broker-errors).

<table border="1" class="nice">
<thead>
<tr>
  <th>Response field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>operation</td>
  <td>string</td>
  <td>For async responses, service brokers can return operation state as a string. This field will be provided back to the service broker on `last_operation` requests as a URL encoded query param.</a>.</td>
</tr>
</tboby>
</table>
\* Fields with an asterisk are required.
<pre class="terminal">
{
 "operation": "task_10"
}
</pre>


## <a id='broker-errors'></a>Broker Errors ##

<!-- **TODO not sure how I missed this section initially. it's good as a top level section. keep it. I don't think we need all the exposition around when we link here.** -->


### Response ###

Broker failures beyond the scope of the well-defined HTTP response codes listed
above (like 410 on delete) should return an appropriate HTTP response code
(chosen to accurately reflect the nature of the failure) and a body containing a valid JSON Object (not an array).

##### Body #####

<!-- **TODO I think this is confusing. Are we trying to say it shouldn't be
an array? Or that it shoudn't have an empty body? We don't use any of
the empty body status codes in this API, but what if we did?** -->

All response bodies must be a valid JSON Object (`{}`). This is for future compatibility; it will be easier to add fields in the future if JSON is expected rather than to support the cases when a JSON body may or may not be returned.

For error responses, the following fields are valid. Others will be ignored. If an empty JSON object is returned in the body `{}`, a generic message containing the HTTP response code returned by the broker will be displayed to the requestor.

<table border="1" class="nice">
<thead>
<tr>
  <th>Response Field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>description</td>
  <td>string</td>
  <td>An error message explaining why the request failed. This message will be displayed to the user who initiated the request.
</td>
</tr>
</tbody>
</table>

<pre class="terminal">
{
  "description": "Something went wrong. Please contact support at http://support.example.com."
}
</pre>

## <a id='orphans'></a>Orphans ##

<!-- **TODO drop** -->


The Cloud Controller is the source of truth for service instances and bindings. Service brokers are expected to have successfully provisioned all the instances and bindings Cloud Controller knows about, and none that it doesn't.

Orphans can result if the broker does not return a response before a request from Cloud Controller times out (typically 60 seconds). For example, if a broker does not return a response to a provision request before Cloud Controller times out, the broker might eventually succeed in provisioning an instance after Cloud Controller considers the request a failure. This results in an orphan instance on the service side.

To mitigate orphan instances and bindings, Cloud Controller will attempt to delete resources it cannot be sure were successfully created, and will keep trying to delete them until the broker responds with a success.

More specifically, when a provision or bind request to the broker fails, Cloud Controller will immediately send a corresponding delete or unbind request. If the delete or unbind request fails, Cloud Controller will retry the delete or unbind request ten times with an exponental backoff schedule (over a period of 34 hours).

<table border="1" class="nice">
<thead>
<tr>
  <th>Status code</th>
  <th>Result</th>
  <th>Orphan mitigation</th>
</tr>
</thead>
<tbody>
<tr>
  <td>200</td>
  <td>Success</td>
  <td></td>
</tr>
<tr>
  <td>200 with malformed response</td>
  <td>Failure</td>
  <td></td>
</tr>
<tr>
  <td>201</td>
  <td>Success</td>
  <td></td>
</tr>
<tr>
  <td>201 with malformed response</td>
  <td>Failure</td>
  <td>Yes</td>
</tr>
<tr>
  <td>All other 2xx</td>
  <td>Failure</td>
  <td>Yes</td>
</tr>
<tr>
  <td>408</td>
  <td>Failure due to timeout</td>
  <td>Yes</td>
</tr>
<tr>
  <td>All other 4xx</td>
  <td>Broker rejects request</td>
  <td></td>
</tr>
<tr>
  <td>5xx</td>
  <td>Broker error</td>
  <td>Yes</td>
</tr>
<tr>
  <td>Timeout</td>
  <td>Failure</td>
  <td>Yes</td>
</tr>
</tbody>
</table>

If the Cloud Controller encounters an internal error provisioning an instance or binding (for example, saving to the database fails), then the Cloud Controller will send a single delete or unbind request to the broker but will not retry.

This orphan mitigation behavior was introduced in cf-release v196.

# Catalog Metadata #

<!-- **TODO This can be a subsection of the API section. Cost needs to be standardized if we care.** -->


The Services Marketplace is defined as the aggregate catalog of services and plans exposed to end users of a Cloud Foundry instance. Marketplace services may come from one or many service brokers. The Marketplace is exposed to end users by cloud controller clients (web, CLI, IDEs, etc), and the Cloud Foundry community is welcome to develop their own clients. All clients are not expected to have the same requirements for information to expose about services and plans. This document discusses user-facing metadata for services and plans, and how the broker API enables broker authors to provide metadata required by different cloud controller clients.

As described in the [Service Broker API](api.md#catalog-mgmt), the only required user-facing fields are `label` and `description` for services, and `name` and `description` for service plans. Rather than attempt to anticipate all potential fields that clients will want, or add endless fields to the API spec over time, the broker API provides a mechanism for brokers to advertise any fields a client requires. This mechanism is the `metadata` field.

The contents of the `metadata` field are not validated by cloud controller but may be by cloud controller clients. Not all clients will make use of the value of `metadata`, and not all brokers have to provide it. If a broker does advertise the `metadata` field, client developers can choose to display some or all fields available.

<p class="note"><strong>Note</strong>: In the <a href="api-v1.md">v1 broker API</a>, the <code>metadata</code> field was called <code>extra</code>.</p>

## <a id='community-driven-standards'></a>Community-Driven Standards ##

This page provides a place to publish the metadata fields required by popular cloud controller clients. Client authors can add their metadata requirements to this document, so that broker authors can see what metadata they should advertise in their catalogs.

**Before adding new fields, consider whether an existing one will suffice.**

<p class="note"><strong>Note</strong>: "CLI strings" are all lowercase, no spaces.
Keep it short; imagine someone having to type it as an argument for a longer CLI
command.</p>

## <a id='services-metadata-fields'></a>Services Metadata Fields ##

| Broker API Field | Type | Description | CC API Field | Pivotal CLI | Pivotal Apps Manager |
|------------------|------|-------------|--------------|-------------|---------------------------|
| name | CLI string | A short name for the service to be displayed in a catalog. | label | X | X |
| description | string | A short 1-line description for the service, usually a single sentence or phrase. | description | X | X |
| metadata.displayName | string | The name of the service to be displayed in graphical clients | extra.displayName | | X |
| metadata.imageUrl | string | The URL to an image. | extra.imageUrl | | X |
| metadata.longDescription | string | Long description | extra.longDescription | | X |
| metadata.providerDisplayName | string | The name of the upstream entity providing the actual service | extra.providerDisplayName | | X |
| metadata.documentationUrl | string | Link to documentation page for service | extra.documentationUrl | | X |
| metadata.supportUrl | string | Link to support for the service | extra.supportUrl | | X |

## <a id='plan-metadata-fields'></a>Plan Metadata Fields ##

| Broker API Field | Type | Description | CC API Field | Pivotal CLI | Pivotal Apps Manager |
|------------------|------|-------------|--------------|-------------|---------------------------|
d| name | CLI string | A short name for the service plan to be displayed in a catalog. | name | X | |
| description | string | A description of the service plan to be displayed in a catalog. | description | | |
| metadata.bullets | array-of-strings | Features of this plan, to be displayed in a bulleted-list | extra.bullets | | X |
| metadata.costs | cost object | An array-of-objects that describes the costs of a service, in what currency, and the unit of measure. If there are multiple costs, all of them could be billed to the user (such as a monthly + usage costs at once).  Each object must provide the following keys:<br/>`amount: { usd: float }, unit: string `<br/>This indicates the cost in USD of the service plan, and how frequently the cost is occurred, such as "MONTHLY" or "per 1000 messages". | extra.costs | | X |
| metadata.displayName | string | Name of the plan to be display in graphical clients. | extra.displayName | | X |

## <a id='example-broker-response'></a>Example Broker Response Body ##

The example below contains a catalog of one service, having one service plan. Of course, a broker can offering a catalog of many services, each having many plans.

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

## <a id='example-cc-response'></a>Example Cloud Controller Response Body ##

<!-- **TODO I don't think we need the response example. What is it trying to show?** -->


```
{
   "metadata":{
      "guid":"bc8748f1-fe05-444d-ab7e-9798e1f9aef6",
      "url":"/v2/services/bc8748f1-fe05-444d-ab7e-9798e1f9aef6",
      "created_at":"2014-01-08T18:52:16+00:00",
      "updated_at":"2014-01-09T03:19:16+00:00"
   },
   "entity":{
      "label":"cloudamqp",
      "provider":"cloudamqp",
      "url":"http://adgw.a1.cf-app.example.com",
      "description":"Managed HA RabbitMQ servers in the cloud",
      "long_description":null,
      "version":"n/a",
      "info_url":null,
      "active":true,
      "bindable":true,
      "unique_id":"18723",
      "extra":{
         "displayName":"CloudAMQP",
         "imageUrl":"https://d33na3ni6eqf5j.cloudfront.net/app_resources/18723/thumbs_112/img9069612145282015279.png",
         "longDescription":"Managed, highly available, RabbitMQ clusters in the cloud",
         "providerDisplayName":"84codesAB",
         "documentationUrl":null,
         "supportUrl":null
      },
      "tags":[
         "amqp",
         "rabbitmq"
      ],
      "requires":[

      ],
      "documentation_url":null,
      "service_plans":[
         {
            "metadata":{
               "guid":"6c4903ab-14ce-41de-adb2-632cf06117a5",
               "url":"/v2/services/6c4903ab-14ce-41de-adb2-632cf06117a5",
               "created_at":"2013-11-01T00:21:25+00:00",
               "updated_at":"2014-01-09T03:19:16+00:00"
            },
            "entity":{
               "name":"bunny",
               "free":true,
               "description":"Big Bunny",
               "service_guid":"bc8748f1-fe05-444d-ab7e-9798e1f9aef6",
               "extra":{
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
               },
               "unique_id":"addonOffering_1889",
               "public":true
            }
         }
      ]
   }
}
```



# Binding Credentials

<!-- **TODO we can't call this thing vcap_services without an explanation of 
what vcap means. it's not clear where this is returned from. should
pick either a combined uri, or individual fields, not both.**-->

A bindable service returns credentials that an application can consume in response to the `cf bind` API call.
Cloud Foundry writes these credentials to the [`VCAP_SERVICES`](../devguide/deploy-apps/environment-variable.html#VCAP-SERVICES) environment variable.
In some cases, buildpacks write a subset of these credentials to other
environment variables that frameworks might need.

Choose from the following list of credential fields if possible, though you can provide additional fields as needed.
Refer to the [Using Bound Services](../devguide/services/managing-services.html#use) section of the
_Managing Service Instances with the CLI_ topic for information on how these
credentials are consumed.

<p class='note'><strong>Note</strong>: If you provide a service that supports a connection string, provide the <code>uri</code> key for buildpacks and
application libraries to use.</p>

<table border="1" class="nice">
  <tr>
    <th><strong>CREDENTIALS</strong></th>
    <th><strong>DESCRIPTION</strong></th>
  </tr>
  <tr>
    <td>uri</td>
    <td>Connection string of the form <code>DB-TYPE://USERNAME:PASSWORD@HOSTNAME:PORT/NAME</code>,
    where <code>DB-TYPE</code> is a type of database such as mysql, postgres, mongodb, or amqp.</td>
  </tr>
  <tr>
    <td>hostname</td>
    <td>FQDN of the server host</td>
  </tr>
  <tr>
    <td>port</td>
    <td>Port of the server host</td>
  </tr>
  <tr>
    <td>name</td>
    <td>Name of the service instance</td>
  </tr>
  <tr>
    <td>vhost</td>
    <td>Name of the messaging server virtual host - a replacement for a <code>name</code> specific to AMQP providers</td>
  </tr>
  <tr>
    <td>username</td>
    <td>Server user</td>
  </tr>
  <tr>
    <td>password</td>
    <td>Server password</td>
  </tr>
</table>

The following is an example output of `ENV['VCAP_SERVICES']`.

<p class='note'><strong>Note</strong>: Depending on the types of databases you are using, each database might return different credentials.</p>

<pre>
VCAP_SERVICES=
{
  cleardb: [
    {
      name: "cleardb-1",
      label: "cleardb",
      plan: "spark",
      credentials: {
        name: "ad_c6f4446532610ab",
        hostname: "us-cdbr-east-03.cleardb.com",
        port: "3306",
        username: "b5d435f40dd2b2",
        password: "ebfc00ac",
        uri: "mysql://b5d435f40dd2b2:ebfc00ac@us-cdbr-east-03.cleardb.com:3306/ad_c6f4446532610ab",
        jdbcUrl: "jdbc:mysql://b5d435f40dd2b2:ebfc00ac@us-cdbr-east-03.cleardb.com:3306/ad_c6f4446532610ab"
      }
    }
  ],
  cloudamqp: [
    {
      name: "cloudamqp-6",
      label: "cloudamqp",
      plan: "lemur",
      credentials: {
        uri: "amqp://ksvyjmiv:IwN6dCdZmeQD4O0ZPKpu1YOaLx1he8wo@lemur.cloudamqp.com/ksvyjmiv"
      }
    }
    {
      name: "cloudamqp-9dbc6",
      label: "cloudamqp",
      plan: "lemur",
      credentials: {
        uri: "amqp://vhuklnxa:9lNFxpTuJsAdTts98vQIdKHW3MojyMyV@lemur.cloudamqp.com/vhuklnxa"
      }
    }
  ],
  rediscloud: [
    {
      name: "rediscloud-1",
      label: "rediscloud",
      plan: "20mb",
      credentials: {
        port: "6379",
        host: "pub-redis-6379.us-east-1-2.3.ec2.redislabs.com",
        password: "1M5zd3QfWi9nUyya"
      }
    },
  ],
}
</pre>

# Examples

<!-- **TODO keep. move to the end?** -->


The following example service broker applications have been developed - these are a great starting point if you are developing your own service broker.

## Ruby

* [GitHub repo service](https://github.com/cloudfoundry-samples/github-service-broker-ruby) - this is designed to be an easy-to-read example of a service broker, with complete documentation, and comes with a demo app that uses the service. The broker can be deployed as an application to any Cloud Foundry instance or hosted elsewhere. The service broker uses GitHub as the service back end.
* [MySQL database service](https://github.com/cloudfoundry/cf-mysql-release) - this broker and its accompanying MySQL server are designed to be deployed together as a [BOSH](https://github.com/cloudfoundry/bosh) release. BOSH is used to deploy or upgrade the release, monitors the health of running components, and restarts or recreates unhealthy VMs. The broker code alone can be found [here](https://github.com/cloudfoundry/cf-mysql-broker).

## Java

* [Spring Cloud - Cloud Foundry Service Broker](https://github.com/spring-cloud/spring-cloud-cloudfoundry-service-broker) - This implements the REST contract for service brokers and the artifacts are published to the spring maven repo.  This greatly simplifies development: include a single dependency in Gradle, implement interfaces, and configure. A sample implementation has been provided for [MongoDB](https://github.com/spring-cloud-samples/cloudfoundry-service-broker).
* [MySQL Java Broker](https://github.com/cloudfoundry-community/cf-mysql-java-broker) - a Java port of the Ruby-based [MySQL broker](https://github.com/cloudfoundry/cf-mysql-broker) above.

## Go

* [Asynchronous Service Broker for AWS EC2](https://github.com/cloudfoundry-samples/go_service_broker) - This broker implements support for the experimental [Asynchronous Service Operations](./api.html#asynchronous-operations), and calls AWS APIs to provision EC2 VMs.



# Managing Service Brokers

<!-- **TODO drop this whole toplevel section. possibly reformat into a 'recommended best practices section'.** -->


This page assumes you are using cf CLI v6.16 or later.

In order to run many of the commands below, you must be authenticated with Cloud
Foundry as an admin user or as a space developer.

## <a id='quick-start'></a>Quick Start ##

Given a service broker that has implemented the [Service Broker API](api.md), two steps are required to make its services available to end users in all organizations or a limited number of organizations by service plan.

1. [Register a Broker](#register-broker)
1. [Make Plans Public](#make-plans-public)

As of cf-release 229, CC API 2.47.0, Cloud Foundry supports both standard brokers and space-scoped private brokers. Standard private brokers can offer service plans privately or publish them to specific organizations or to all users. Space-scoped private brokers publish services only to users within the space they are registered to.

## <a id='register-broker'></a>Register a Broker ##

Registering a broker causes Cloud Controller to fetch and validate the catalog
from your broker, and save the catalog to the Cloud Controller database.
The basic auth username and password which are provided when adding a broker are
encrypted in Cloud Controller database, and used by the Cloud Controller to
authenticate with the broker when making all API calls.
Your service broker should validate the username and password sent in every
request; otherwise, anyone could curl your broker to delete service instances.

### Standard Private Brokers ###

<pre class="terminal">
$ cf create-service-broker mybrokername someuser somethingsecure http://mybroker.example.com/
</pre>

### Space-Scoped Private Brokers ###

<pre class="terminal">
$ cf create-service-broker mybrokername someuser somethingsecure http://mybroker.example.com/ --space-scoped
</pre>


### <a id='make-plans-public'></a>Make Plans Public ###
New service plans from standard brokers are private by default. To make plans available to end users, see [Make Plans Public](./access-control.md#enable-access). Instances of a private plan cannot be provisioned until either the plan is made public or is made available to an organization.

New service plans from space-scoped private brokers are automatically published to all users in the broker's space. It is not possible to manage visibility of a space-scoped private broker at the Cloud Foundry instance or organization level.

### <a id='multiple-brokers'></a>Multiple Brokers, Services, Plans ###

Many service brokers may be added to a Cloud Foundry instance, each offering
many services and plans.
The following constraints should be kept in mind:

- It is not possible to have multiple brokers with the same name
- It is not possible to have multiple brokers with the same base URL
- The service ID and plan IDs of each service advertised by the broker must be unique across Cloud Foundry. GUIDs are recommended for these fields.

See [Possible Errors](#possible-errors) below for error messages and what do to
when you see them.

## <a id='list-brokers'></a> List Service Brokers ##

<pre class="terminal">
$ cf service-brokers
Getting service brokers as admin...Cloud Controller
OK

Name            URL
my-service-name http://mybroker.example.com
</pre>

## <a id='update-broker'></a>Update a Broker ##

Updating a broker is how to ingest changes a broker author has made into Cloud
Foundry.
Similar to adding a broker, update causes Cloud Controller to fetch the catalog
from a broker, validate it, and update the Cloud Controller database with any
changes found in the catalog.

Update also provides a means to change the basic auth credentials cloud
controller uses to authenticate with a broker, as well as the base URL of the
broker's API endpoints.

<pre class="terminal">
$ cf update-service-broker mybrokername someuser somethingsecure http://mybroker.example.com/
</pre>

## <a id='rename-broker'></a>Rename a Broker ##

A service broker can be renamed with the `rename-service-broker`
command.
This name is used only by the Cloud Foundry operator to identify brokers, and
has no relation to configuration of the broker itself.

<pre class="terminal">
$ cf rename-service-broker mybrokername mynewbrokername
</pre>



## <a id='remove-broker'></a>Remove a Broker ##

Removing a service broker will remove all services and plans in the broker's catalog from the Cloud Foundry Marketplace.

```
$ cf delete-service-broker mybrokername
```

**Note**: Attempting to remove a service broker will fail if there are service instances for any service plan in its catalog. When planning to shut down or delete a broker, make sure to remove all service instances first. Failure to do so will leave orphaned service instances in the Cloud Foundry database. If a service broker has been shut down without first deleting service instances, you can remove the instances with the CLI; see [Purge a Service](#purge-a-service).

### <a id='purge-service'></a>Purge a Service ###

If a service broker has been shut down or removed without first deleting service instances from Cloud Foundry, you will be unable to remove the service broker or its services and plans from the Marketplace. In development environments, broker authors often destroy their broker deployments and need a way to clean up the Cloud Controller database.

The following command will delete a service offering, all of its plans, as well as all associated service instances and bindings from the Cloud Controller database, without making any API calls to a service broker. For services from v1 brokers, you must provide a provider with `-p PROVIDER`. Once all services for a broker have been purged, the broker can be removed normally.

<pre class="terminal">
$ cf purge-service-offering v1-test -p pivotal-software
Warning: This operation assumes that the service broker responsible for this
service offering is no longer available, and all service instances have been
deleted, leaving orphan records in Cloud Foundry's database. All knowledge of
the service will be removed from Cloud Foundry, including service instances and
service bindings. No attempt will be made to contact the service broker; running
this command without destroying the service broker will cause orphan service
instances. After running this command you may want to run either
delete-service-auth-token or delete-service-broker to complete the cleanup.

Really purge service offering v1-test from Cloud Foundry? y
OK
</pre>


### <a id='purge-service-instance'></a>Purge a Service Instance###

The following command will delete a single service instance, its service bindings and its service keys from the Cloud Controller database, without making any API calls to a service broker.
This can be helpful in instances a Service Broker is not conforming to the Service Broker API and not returning a 200 or 410 to requests to delete the service instance.

<pre class="terminal">
$ cf purge-service-instance mysql-dev
WARNING: This operation assumes that the service broker responsible for this
service instance is no longer available or is not responding with a 200 or 410,
and the service instance has been deleted, leaving orphan records in Cloud
Foundry's database. All knowledge of the service instance will be removed from
Cloud Foundry, including service bindings and service keys.

Really purge service instance mysql-dev from Cloud Foundry?> y
Purging service mysql-dev...
OK
</pre>

`purge-service-instance` requires cf-release v218 and cf CLI 6.14.0.

## <a id='possible-errors'></a>Possible Errors ##

If incorrect basic auth credentials are provided:

<pre class="terminal">
Server error, status code: 500, error code: 10001, message: Authentication
failed for the service broker API.
Double-check that the username and password are correct:
    http://github-broker.a1-app.example.com/v2/catalog
</pre>

If you receive the following errors, check your broker logs.
You may have an internal error.

<pre class="terminal">
Server error, status code: 500, error code: 10001, message:
    The service broker response was not understood

Server error, status code: 500, error code: 10001, message:
    The service broker API returned an error from
    http://github-broker.a1-app.example.com/v2/catalog: 404 Not Found

Server error, status code: 500, error code: 10001, message:
    The service broker API returned an error from
    http://github-broker.primo.example.com/v2/catalog: 500 Internal Server Error
</pre>

If your broker's catalog of services and plans violates validation of presence,
uniqueness, and type, you will receive meaningful errors.

<pre class="terminal">
Server error, status code: 502, error code: 270012, message: Service broker catalog is invalid:
Service service-name-1
  service id must be unique
  service description is required
  service "bindable" field must be a boolean, but has value "true"
  Plan plan-name-1
    plan metadata must be a hash, but has value [{"bullets"=>["bullet1", "bullet2"]}]
</pre>

# Access Control

<!-- **TODO drop this entire toplevel section. unconditionally. nothing of relevancy here.** -->


All new service plans from standard private brokers are private by default. This means that when adding a new broker, or when adding a new plan to an existing broker's catalog, service plans won't immediately be available to end users. This lets an admin control which service plans are available to end users, and manage limited service availability.

Space-scoped private brokers are registered to a specific space, and all users within that space can automatically access the broker's service plans. With space-scoped brokers, service visibility is not managed separately.

## <a id='prerequisites'></a>Prerequisites ###
- CLI v6.4.0
- Cloud Controller API v2.9.0 (cf-release v179)
- Admin user access; the following commands can be run only by an admin user

## <a id='display-access'></a>Display Access to Service Plans ###

The `service-access` CLI command enables an admin to see the current access control setting for every service plan in the marketplace, across all service brokers.

<pre class="terminal">
$ cf service-access
getting service access as admin...
broker: p-riakcs
   service    plan        access    orgs
   p-riakcs   developer   limited

broker: p-mysql
   service   plan        access   orgs
   p-mysql   100mb-dev   all
</pre>

The `access` column has values `all`, `limited`, or `none`. `all` means a service plan is available to all users of the Cloud Foundry instance; this is what we mean when we say the plan is "public". `none` means the plan is not available to  anyone; this is what we mean when we say the plan is "private". `limited` means that the service plan is available to users of one or more select organizations. When a plan is `limited`, organizations that have been granted access are listed.

Flags provide filtering by broker, service, and organization.

<pre class="terminal">
$ cf help service-access
NAME:
   service-access - List service access settings

USAGE:
   cf service-access [-b BROKER] [-e SERVICE] [-o ORG]

OPTIONS:
   -b     access for plans of a particular broker
   -e     access for plans of a particular service offering
   -o     plans accessible by a particular organization
</pre>

## <a id='enable-access'></a>Enable Access to Service Plans ###

Service access is managed at the granularity of service plans, though CLI commands allow an admin to modify all plans of a service at once.

Enabling access to a service plan for organizations allows users of those organizations to see the plan listed in the marketplace (`cf marketplace`), and if users have the Space Developer role in a targeted space, to provision instances of the plan.

<pre class="terminal">
$ cf enable-service-access p-riakcs
Enabling access to all plans of service p-riakcs for all orgs as admin...
OK

$ cf service-access
getting service access as admin...
broker: p-riakcs
   service    plan        access   orgs
   p-riakcs   developer   all
</pre>

An admin can use `enable-service-access` to:

- Enable access to all plans of a service for users of all orgs (access:`all`)
- Enable access to one plan of a service for users of all orgs (access:`all`)
- Enable access to all plans of a service for users of a specified organization (access: `limited`)
- Enable access to one plan of a service for users of a specified organization (access: `limited`)

<pre class="terminal">
$ cf help enable-service-access
NAME:
   enable-service-access - Enable access to a service or service plan for one or all orgs

USAGE:
   cf enable-service-access SERVICE [-p PLAN] [-o ORG]

OPTIONS:
   -p     Enable access to a particular service plan
   -o     Enable access to a particular organization
</pre>

## <a id='disable-access'></a>Disable Access to Service Plans ###

<pre class="terminal">
$ cf disable-service-access p-riakcs
Disabling access to all plans of service p-riakcs for all orgs as admin...
OK

$ cf service-access
getting service access as admin...
broker: p-riakcs
   service    plan        access   orgs
   p-riakcs   developer   none
</pre>

An admin can use the `disable-service-access` command to:

- Disable access to all plans of a service for users of all orgs (access:`all`)
- Disable access to one plan of a service for users of all orgs (access:`all`)
- Disable access to all plans of a service for users of select orgs (access: `limited`)
- Disable access to one plan of a service for users of select orgs (access: `limited`)

<pre class="terminal">
$ cf help disable-service-access
NAME:
   disable-service-access - Disable access to a service or service plan for one or all orgs

USAGE:
   cf disable-service-access SERVICE [-p PLAN] [-o ORG]

OPTIONS:
   -p     Disable access to a particular service plan
   -o     Disable access to a particular organization
</pre>

### Limitations ####

- You cannot disable access to a service plan for an organization if the plan is currently available to all organizations. You must first disable access for all organizations; then you can enable access for a  particular organization.


# Dashboard Single Sign-On

<!-- **TODO Are we supporting this at all?** -->


## <a id='introduction'></a>Introduction ##

Single sign-on (SSO) enables Cloud Foundry users to authenticate with third-party service dashboards using their Cloud Foundry credentials. Service dashboards are web interfaces which enable users to interact with some or all of the features the service offers. SSO provides a streamlined experience for users, limiting repeated logins and multiple accounts across their managed services. The user's credentials are never directly transmitted to the service since the OAuth2 protocol handles authentication.

Dashboard SSO was introduced in [cf-release v169](https://github.com/cloudfoundry/cf-release/tree/v169) so this or a newer version is required to support the feature.

## <a id='enabling-the-feature-in-cloudfoundry'></a>Enabling the feature in Cloud Foundry ##

To enable the SSO feature, the Cloud Controller requires a UAA client with sufficient permissions to create and delete clients for the service brokers that request them. This client can be configured by including the following snippet in the cf-release manifest:

  ```
  properties:
    uaa:
      clients:
        cc-service-dashboards:
          secret: cc-broker-secret
          scope: openid,cloud_controller_service_permissions.read
          authorities: clients.read,clients.write,clients.admin
          authorized-grant-types: authorization_code,client_credentials
  ```

When this client is not present in the cf-release manifest, Cloud Controller cannot manage UAA clients and an operator will receive a warning when creating or updating service brokers that advertise the `dashboard_client` properties discussed below.

## <a id='broker-responsibilities'></a>Service Broker Responsibilities ##

### <a id='registering-dashboard-client'></a>Registering the Dashboard Client ###

1.  A service broker must include the `dashboard_client` field in the JSON response from its [catalog endpoint](api.html#catalog-mgmt) for each service implementing this feature. A valid response would appear as follows:

      ```
      {
        "services": [
          {
            "id": "44b26033-1f54-4087-b7bc-da9652c2a539",
            ...
            "dashboard_client": {
              "id": "p-mysql-client",
              "secret": "p-mysql-secret",
              "redirect_uri": "http://p-mysql.example.com"
            }
          }
        ]
      }
      ```
    The `dashboard_client` field is a hash containing three fields:
    - `id` is the unique identifier for the OAuth2 client that will be created for your service dashboard on the token server (UAA), and will be used by your dashboard to authenticate with the token server (UAA).
    - `secret` is the shared secret your dashboard will use to authenticate with the token server (UAA).
    - `redirect_uri` is used by the token server as an additional security precaution. UAA will not provide a token if the callback URL declared by the service dashboard doesn't match the domain name in `redirect_uri`. The token server matches on the domain name, so any paths will also match; e.g. a service dashboard requesting a token and declaring a callback URL of `http://p-mysql.example.com/manage/auth` would be approved if `redirect_uri` for its client is `http://p-mysql.example.com/`.

1. When a service broker which advertises the `dashboard_client` property for any of its services is [added or updated](managing-service-brokers.html), Cloud Controller will create or update UAA clients as necessary. This client will be used by the service dashboard to authenticate users.

### <a id='dashboard-url'></a>Dashboard URL ###

A service broker should return a URL for the `dashboard_url` field in response to a [provision request](./api.html#provisioning). Cloud Controller clients should expose this URL to users. `dashboard_url` can be found in the response from Cloud Controller to create a service instance, enumerate service instances, space summary, and other endpoints.

Users can then navigate to the service dashboard at the URL provided by `dashboard_url`, initiating the OAuth2 login flow.

## <a id='dashboard-responsibilities'></a>Service Dashboard Responsibilities ##

### <a id='oauth2-flow'></a>OAuth2 Flow ###

When a user navigates to the URL from `dashboard_url`, the service dashboard should initiate the OAuth2 login flow. A summary of the flow can be found in [section 1.2 of the OAuth2 RFC](http://tools.ietf.org/html/rfc6749#section-1.2). OAuth2 expects the presence of an [Authorization Endpoint](http://tools.ietf.org/html/rfc6749#section-3.1) and a [Token Endpoint](http://tools.ietf.org/html/rfc6749#section-3.2). In Cloud Foundry, these endpoints are provided by the UAA. Clients can discover the location of UAA from Cloud Controller's info endpoint; in the response the location can be found in the `token_endpoint` field.

```
$ curl api.example.com/info
{"name":"vcap","build":"2222","support":"http://support.example.com","version
":2,"description":"Cloud Foundry sponsored by Pivotal","authorization_endpoint":
"https://login.example.com","token_endpoint":"https://uaa.example.com",
"allow_debug":true}
```

<p class='note'>To enable service dashboards to support SSO for service instances created from different Cloud Foundry instances, the /v2/info url is sent to service brokers in the `X-Api-Info-Location` header of every API call. A service dashboard should be able to discover this URL from the broker, and enabling the dashboard to contact the appropriate UAA for a particular service instance.</p>

A service dashboard should implement the OAuth2 Authorization Code Grant type ([UAA docs](https://github.com/cloudfoundry/uaa/blob/master/docs/UAA-APIs.rst#authorization-code-grant), [RFC docs](http://tools.ietf.org/html/rfc6749#section-4.1)).

1. When a user visits the service dashboard at the value of `dashboard_url`, the dashboard should redirect the user's browser to the Authorization Endpoint and include its `client_id`, a `redirect_uri` (callback URL with domain matching the value of `dashboard_client.redirect_uri`), and list of requested scopes.

    Scopes are permissions included in the token a dashboard client will receive from UAA, and which Cloud Controller uses to enforce access. A client should request the minimum scopes it requires. The minimum scopes required for this workflow are `cloud_controller_service_permissions.read` and `openid`. For an explanation of the scopes available to dashboard clients, see [On Scopes](#on-scopes).

1. UAA authenticates the user by redirecting the user to the Login Server, where the user then approves or denies the scopes requested by the service dashboard. The user is presented with human readable descriptions for permissions representing each scope. After authentication, the user's browser is redirected back to the Authorization endpoint on UAA with an authentication cookie for the UAA.

1. Assuming the user grants access, UAA redirects the user's browser back to the value of `redirect_uri` the dashboard provided in its request to the Authorization Endpoint.  The `Location` header in the response includes an authorization code.

    ```
    HTTP/1.1 302 Found
    Location: https://p-mysql.example.com/manage/auth?code=F45jH
    ```

1. The dashboard UI should then request an access token from the Token Endpoint by including the authorization code received in the previous step.  When making the request the dashboard must authenticate with UAA by passing the client `id` and `secret` in a basic auth header. UAA will verify that the client id matches the client it issued the code to. The dashboard should also include the `redirect_uri` used to obtain the authorization code for verification.

1. UAA authenticates the dashboard client, validates the authorization code, and ensures that the redirect URI received matches the URI used to redirect the client when the authorization code was issues.  If valid, UAA responds back with an access token and a refresh token.

### <a id='checking-user-permissions'></a>Checking User Permissions ###

UAA is responsible for authenticating a user and providing the service with an access token with the requested permissions. However, after the user has been logged in, it is the responsibility of the service dashboard to verify that the user making the request to manage an instance currently has access to that service instance.

The service can accomplish this with a GET to the `/v2/service_instances/:guid/permissions` endpoint on the Cloud Controller. The request must include a token for an authenticated user and the service instance guid. The token is the same one obtained from the UAA in response to a request to the Token Endpoint, described above.
.

Example Request:

```
curl -H 'Content-Type: application/json' \
       -H 'Authorization: bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoid' \
       http://api.cloudfoundry.com/v2/service_instances/44b26033-1f54-4087-b7bc-da9652c2a539/permissions

```

Response:

```
{
  "manage": true
}
```

The response will indicate to the service whether this user is allowed to manage the given instance. A `true` value for the `manage` key indicates sufficient permissions; `false` would indicate insufficient permissions.  Since administrators may change the permissions of users, the service should check this endpoint whenever a user uses the SSO flow to access the service's UI.

### <a id="on-scopes"></a> On Scopes

Scopes let you specify exactly what type of access you need. Scopes limit access for OAuth tokens. They do not grant any additional permission beyond that which the user already has.

#### Minimum Scopes
The following two scopes are necessary to implement the integration. Most dashboard shouldn't need more permissions than these scopes enabled.

| Scope                                            | Permissions   |
| ------------------------------------------------ | ------------- |
| `openid`                                         | Allows access to basic data about the user, such as email addresses |
| `cloud_controller_service_permissions.read`      | Allows access to the CC endpoint that specifies whether the user can manage a given service instance |

#### Additional Scopes
Dashboards with extended capabilities may need to request these additional scopes:

| Scope                                            | Permissions   |
| ------------------------------------------------ | ------------- |
| `cloud_controller.read`                          | Allows read access to all resources the user is authorized to read |
| `cloud_controller.write`                         | Allows write access to all resources the user is authorized to update / create / delete |

## <a id='reference-implementation'></a>Reference Implementation ##

The [MySQL Service Broker][example-broker] is an example of a broker that also implements a SSO dashboard. The login flow is implemented using the [OmniAuth library](https://github.com/intridea/omniauth) and a custom [UAA OmniAuth Strategy](https://github.com/cloudfoundry/omniauth-uaa-oauth2). See this [OmniAuth wiki page](https://github.com/intridea/omniauth/wiki/Strategy-Contribution-Guide) for instructions on how to create your own strategy.

The UAA OmniAuth strategy is used to first get an authorization code, as documented in [this section](https://github.com/cloudfoundry/uaa/blob/master/docs/UAA-APIs.rst#authorization-code-grant) of the UAA documentation. The user is redirected back to the service (as specified by the `callback_path` option or the default `auth/cloudfoundry/callback` path) with the authorization code. Before the application / action is dispatched, the OmniAuth strategy uses the authorization code to [get a token](https://github.com/cloudfoundry/uaa/blob/master/docs/UAA-APIs.rst#client-obtains-token-post-oauth-token) and uses the token to request information from UAA to fill the `omniauth.auth` environment variable. When OmniAuth returns control to the application, the `omniauth.auth` environment variable hash will be filled with the token and user information obtained from UAA as seen in the [Auth Controller](https://github.com/cloudfoundry/cf-mysql-broker/blob/master/app/controllers/manage/auth_controller.rb).

## <a id='restrictions'></a>Restrictions ##

 * UAA clients are scoped to services.  There must be a `dashboard_client` entry for each service that uses SSO integration.
 * Each `dashboard_client id` must be unique across the CloudFoundry deployment.

## <a id="resources"></a>Resources ##

  * [OAuth2](http://oauth.net/2/)
  * [Example broker with SSO implementation][example-broker]
  * [Cloud Controller API Docs](http://apidocs.cfapps.io/)
  * [User Account and Authentication (UAA) Service APIs](https://github.com/cloudfoundry/uaa/blob/master/docs/UAA-APIs.rst)

[example-broker]: https://github.com/cloudfoundry/cf-mysql-broker

# Application Log Streaming

<!-- **TODO are we going to support this? seems overly specialized. seems 
like it should be just additional credentials on the same service
binding.** -->

By binding an application to an instance of an applicable service, Cloud Foundry will stream logs for the bound application to the service instance.

- Logs for all apps bound to a log-consuming service instance will be streamed to that instance
- Logs for an app bound to multiple log-consuming service instances will be streamed to all instances

To enable this functionality, a service broker must implement the following:

1. In the [catalog](api.html#catalog-mgmt) endpoint, the broker must include `requires: syslog_drain`. This minor security measure validates that a service returning a `syslog_drain_url` in response to the [bind](api.html#binding) operation has also declared that it expects log streaming. If the broker does not include `requires: syslog_drain`, and the bind request returns a value for `syslog_drain_url`, Cloud Foundry will return an error for the bind operation.

2. In response to a [bind](api.html#binding) request, the broker should return a value for `syslog_drain_url`. The syslog URL has a scheme of syslog, syslog-tls, or https and can include a port number. For example:

    `"syslog_drain_url": "syslog://logs.example.com:1234"`

## How does it work?

1. Service broker returns a value for `syslog_drain_url` in response to bind
1. Loggregator periodically polls CC `/v2/syslog_drain_urls` for updates
1. Upon discovering a new `syslog_drain_url`, Loggregator identifies the associated app
1. Loggregator streams app logs for that app to the locations specified by the service instances' `syslog_drain_url`s

Users can manually configure app logs to be streamed to a location of their choice using User-provided Service Instances. For details, see [Using Third-Party Log Management Services](../devguide/services/log-management.html).

# Route Services

<!-- **TODO very cf specific ** -->


This documentation is intended for service authors who are interested in offering a service to a Cloud Foundry services marketplace. Developers interested in consuming these services can read the  [Manage Application Requests with Route Services](../devguide/services/route-binding.md) topic.

FIXME: <%= vars.route_services %>

## <a id='introduction'></a>Introduction ##

Cloud Foundry application developers may wish to apply transformation or processing to requests before they reach an application. Common examples of use cases are authentication, rate limiting, and caching services. Route Services are a new kind of Marketplace Service that developers can use to apply various transformations to application requests by binding an application's route to a service instance. Through integrations with service brokers and optionally with the Cloud Foundry routing tier, providers can offer these services to developers with a familiar automated, self-service, and on-demand user experience.

## <a id='architecture'></a>Architecture ##

Cloud Foundry supports three models for Route Services: fully-brokered services; static, brokered services; and user-provided services. In each case, you configure a route service to process traffic addressed to an app.

### <a id="fully-brokered"></a>Fully-Brokered Service ###

In this model, the CF router receives all traffic to apps in the deployment before any processing by the route service. Developers can bind a route service to any app, and if an app is bound to a route service, the CF router sends its traffic to the service. After the route service processes requests, it sends them back to the load balancer in front of the CF router. The second time through, the CF router recognizes that the route service has already handled them, and forwards them directly to app instances.

![Fully brokered](images/route-services-fully-brokered.png)

The route service can run inside or outside of CF, so long as it fulfills the [Service Instance Responsibilities](#service-instance-responsibilities) to integrate it with the CF router. A service broker publishes the route service to the CF marketplace, making it available to developers. Developers can then create an instance of the service and bind it to their apps with the following commands:

`cf create-service BROKER_SERVICE_PLAN SERVICE_INSTANCE`
`cf bind-route-service YOUR_APP_DOMAIN SERVICE_INSTANCE [--hostname HOSTNAME]`

Developers configure the service either through the service provider's web interface or by passing [arbitrary parameters](../devguide/services/managing-services.html#arbitrary-params-create) to their `cf create-service` call, through the `-c` flag.

**Advantages:**

- Developers can use a Service Broker to dynamically configure how the route service processes traffic to specific applications.
- Adding route services requires no manual infrastructure configuration.
- Traffic to apps that do not use the service makes fewer network hops; requests for those apps do not pass through the route service.

**Disadvantages:**

- Traffic to apps that use the route service makes additional network hops, as compared to the static model.

### <a id="static-brokered"></a>Static, Brokered Service ###

In this model, an operator installs a static routing service, which might be a piece of hardware, in front of the Load Balancer. The routing service runs outside of Cloud Foundry and receives traffic to all apps running in the CF deployment. The service provider creates a service broker to publish the service to the CF marketplace. As with a [fully-brokered service](#fully-brokered), a developer can use the service by instantiating it with `cf create-service` and binding it to an app with `cf bind-route-service`.

![Static, brokered](images/route-services-static-brokered.png)

In this model, you configure route services on an app-by-app basis. When you bind a service to an app, the service broker directs the routing service to process that app's traffic rather than pass the requests through unchanged.

**Advantages:**

- Developers can use a Service Broker to dynamically configure how the route service processes traffic to specific applications.
- Traffic to apps that use the route service takes fewer network hops.

**Disadvantages:**

- Adding route services requires manual infrastructure configuration.
- Unnecessary network hops for traffic to apps that do not use the route service; requests for all apps hosted by the the deployment pass through the route service component.

### <a id="user-provided"></a>User-Provided Service ###

If a route service is not listed in the CF marketplace by a broker, a developer can still bind it to their app as a User-Provided service. The service can run anywhere, either inside or outside of CF, but it must fulfill the integration requirements described in [Service Instance Responsibilities](#service-instance-responsibilities). The service also needs to be reachable by an outbound connection from the CF Router.

![User-provided](images/route-services-user-provided.png)

This model is identical to the [fully-brokered service](#fully-brokered) model, except without the broker. Developers configure the service manually, outside of Cloud Foundry. They can then create a user-provided service instance and bind it to their application with the following commands, supplying the URL of their route service:

`cf create-user-provided-service SERVICE_INSTANCE -r ROUTE_SERVICE_URL`
`cf bind-route-service DOMAIN SERVICE_INSTANCE [--hostname HOSTNAME]`

**Advantages:**

- Adding route services requires no manual infrastructure configuration.
- Traffic to apps that do not use the service makes fewer network hops; requests for those apps do not pass through the route service.

**Disadvantages:**

- Developers must manually provision/configure route services out of the context of Cloud Foundry; no service broker automates these operations.
- Traffic to apps that use the route service makes additional network hops, as compared to the static model.

### <a id="architecture-comparison"></a>Architecture Comparison ###

The models above require the [broker](#broker-responsibilities) and [service instance](#service-instance-responsibilities) responsibilities below, as summarized in the following table:

<table id='architecture-comparison' border="1" class="nice" >
  <tr>
    <th>Route Services Architecture</th>
    <th>Fulfills CF <a href="#service-instance-responsibilities">Service Instance Responsibilities</a></th>
    <th>Fulfills CF <a href="#broker-responsibilities">Broker Responsibilities</a></th>
  </tr><tr>
    <td>Fully-Brokered</td>
    <td>Yes</td>
    <td>Yes</td>
  </tr><tr>
    <td>Static Brokered</td>
    <td>No</td>
    <td>Yes</td>
  </tr><tr>
    <td>User-Provided</td>
    <td>Yes</td>
    <td>No</td>
  </tr>
</table>

FIXME: <%= vars.route_services_config %>

## <a id='service-instance-responsibilities'></a>Service Instance Responsibilities ##

The following applies only when a broker returns `route_service_url` in the bind response.

#### <a id='how-it-works'></a>How It Works ####

Binding a service instance to a route will associate the `route_service_url` with the route in the Cloud Foundry router. All requests for the route will be proxied to the URL specified by `route_service_url`.

Once a route service completes its function, it is expected to forward the request to the route the original request was sent to. The Cloud Foundry router will include a header that provides the address of the route, as well as two headers that are used by the route itself to validate the request sent by the route service.

#### <a id='headers'></a>Headers ####
The `X-CF-Forwarded-Url` header contains the URL of the application route. The route service should forward the request to this URL.

The route service should not strip off the `X-CF-Proxy-Signature` and `X-CF-Proxy-Metadata`, as the GoRouter relies on these headers to validate that the request.

#### <a id='ssl-certs'></a>SSL Certificates ####

When Cloud Foundry is deployed in a development environment, certificates hosted by the load balancer will be self-signed (not signed by a trusted certificate authority). When the route service has finished processing an inbound request, and makes a call to the value of `X-CF-Forwarded-Url`, be prepared to accept the self-signed certificate when integrating with a non-production deployment of Cloud Foundry.

#### <a id='timeouts'></a>Timeouts ####

Route services must forward the request to the application route within the number of seconds configured by the `router.route_service_timeout` property (default 60 seconds).

In addition, all requests must respond in the number of seconds configured by the `request_timeout_in_seconds` property (default 900 seconds).

Timeouts are configurable for the router using the cf-release BOSH deployment manifest. For more information, see the [spec](https://github.com/cloudfoundry-incubator/routing-release/blob/master/jobs/gorouter/spec).

## <a id='broker-responsibilities'></a>Broker Responsibilities ##

#### <a id='catalog'></a>Catalog Endpoint ####
Brokers must include `requires: ["route_forwarding"]` for a service in the catalog endpoint. If this is not present, Cloud Foundry will not permit users to bind an instance of the service to a route.

#### <a id='binding'></a>Binding Endpoint ####
When users bind a route to a service instance, Cloud Foundry will send a [bind request](http://docs.cloudfoundry.org/services/api.html#binding) to the broker, including the route address with `bind_resource.route`. A route is an address used by clients to reach apps mapped to the route. The broker may return `route_service_url`, containing a URL where Cloud Foundry should proxy requests for the route. This URL must have a `https` scheme, otherwise the Cloud Controller will reject the binding. `route_service_url` is optional; not returning this field enables a broker to dynamically configure a network component already in the request path for the route, requiring no change in the Cloud Foundry router.

## <a id='examples'></a>Example Route Services ##
- [Logging Route Service](https://github.com/cloudfoundry-samples/logging-route-service): This route service can be pushed as an app to Cloud Foundry. It fulfills the service instance responsibilities above and logs requests received and sent. It can be used to see the route service integration in action by tailing its logs.
- [Rate Limiting Route Service](https://github.com/cloudfoundry-samples/ratelimit-service): This example route service is a simple Cloud Foundry app that provides rate limiting to control the rate of traffic to an application.
- [Spring Boot Example](https://github.com/nebhale/route-service-example): Logs requests received and sent; written in Spring Boot

## <a id='tutorial'></a>Tutorial ##

The following instructions show how to use the [Logging Route Service](https://github.com/cloudfoundry-samples/logging-route-service) described in <a href="#examples">Example Route Services</a> to verify that when a route service is bound to a route, requests for that route are proxied to the route service.

A video of this tutorial is available on [Youtube](https://youtu.be/VaaZJE2E4jI).

Requires CLI version 6.16 or above.

1. Push the [Logging Route Service](https://github.com/cloudfoundry-samples/logging-route-service) as an app.

    <pre class="terminal">
    $ cf push logger
    </pre>

1. Create a user-provided service instance, and include the route of the [Logging Route Service](https://github.com/cloudfoundry-samples/logging-route-service) you pushed as `route_service_url`. Be sure to use `https` for the scheme.

    <pre class="terminal">
    $ cf create-user-provided-service mylogger -r https://logger.cf.example.com
    </pre>

1. Push a sample app like [Spring Music](https://github.com/cloudfoundry-samples/spring-music). By default this will create a route `spring-music.cf.example.com`.

    <pre class="terminal">
    $ cf push spring-music
    </pre>

1. Bind the user-provided service instance to the route of your sample app. The `bind-route-service` command takes a route and a service instance; the route is specified in the following example by domain `cf.example.com` and hostname `spring-music`.

    <pre class="terminal">
    $ cf bind-route-service cf.example.com mylogger --hostname spring-music
    </pre>

1. Tail the logs for your route service.

    <pre class="terminal">
    $ cf logs logger
    </pre>

1. Send a request to the sample app and see in the route service logs that the request is forwarded to it.

    <pre class="terminal">
    $ curl spring-music.cf.example.com
    </pre>

# Supporting Multiple Cloud Foundry Instances

<!-- **TODO the concept is sound, but I'm not sure how it generalizes.** -->


It is possible to register a service broker with multiple Cloud Foundry instances. It may be necessary for the broker to know which Cloud Foundry instance is making a given request. For example, when using [Dashboard Single Sign-On](dashboard-sso.html), the broker is expected to interact with the authorization and token endpoints for a given Cloud Foundry instance.

There are two strategies that can be used to discover which Cloud Foundry instance is making a given request.

## Routing & Authentication
The broker can use unique credentials and/or a unique url for each Cloud Foundry instance. When registering the broker, different Cloud Foundry instances can be configured to use different base urls that include a unique id. For example:

* On Cloud Foundry instance 1, the service broker is registered with the url `broker.example.com/123`
* On Cloud Foundry instance 2, the service broker is registered with the url `broker.example.com/456`

## X-Api-Info-Location Header

All calls to the broker from Cloud Foundry include an `X-Api-Info-Location` header containing the `/v2/info` url for that instance. The `/v2/info` endpoint will return further information, including the location of that Cloud Foundry instance's UAA.

Support for this header was introduced in cf-release v212.

# Volume Services (Experimental)

<!-- **TODO cf experimental. drop or keep? seems to be a way of providing 
access to a service where you provide a file interface instead of a
url** -->

## <a id='introduction'></a>Introduction ##

Cloud Foundry application developers may want their applications to mount one or more volumes in order to write to a reliable, non-ephemeral file system. By integrating with service brokers and the Cloud Foundry runtime, providers can offer these services to developers through an automated, self-service, and on-demand user experience.

<p class="note"><strong>Note</strong>: This feature is experimental.</p>

## <a id='schema'></a>Schema ##

### Service Broker Bind Response ###

<table border="1" class="nice">
<thead>
<tr>
  <th>Field</th>
  <th>Type</th>
  <th>Description</th>
</tr>
</thead>
<tbody>
<tr>
  <td>volume_mounts*</td>
  <td>volume_mount[]</td>
  <td>An array of <i>volume_mount</i> JSON objects</td>
</tr>
</tbody>
</table>

### volume_mount ###

A `volume_mount` represents a remote storage device to be attached and mounted into the app container filesystem via a Volume Driver.


| Field          | Type          | Description                                                                                                   |
|----------------|---------------|---------------------------------------------------------------------------------------------------------------|
| driver         | string        | Name of the volume driver plugin which manages the device                                                     |
| container\_dir | string        | The directory to mount inside the application container                                                       |
| mode           | string        | `"r"` to mount the volume read-only, or `"rw"` to mount it read-write                                         |
| device\_type   | string        | A string specifying the type of device to mount. Currently only `"shared"` devices are supported.             |
| device         | device-object | Device object containing device\_type specific details. Currently only `shared_device` devices are supported. |
 
### shared_device ###

A `shared_device` is a subtype of a device. It represents a distributed file system which can be mounted on all app instances simultaneously.


| Field         | Type   | Description                                                                           |
|---------------|--------|---------------------------------------------------------------------------------------|
| volume\_id    | string | ID of the shared volume to mount on every app instance                                |
| mount\_config | object | Configuration object to be passed to the driver when the volume is mounted (optional) |

### Example ###

```
{
  ...
  "volume_mounts": [
    {
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
    }
  ]
}
```

