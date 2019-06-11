# Service Broker API Release Notes

## [v2.15](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md)
2019-06-11

* Added a delay to polling response for [Service Instance](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#polling-last-operation-for-service-instances)
  and [Service Binding](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#polling-last-operation-for-service-bindings)
  last operations.
* Added a Service Offering specific [async polling timeout](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#polling-interval-and-duration).
* Allow a Service Instance context to be updated and add [org name, space name, and instance names](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#updating-a-service-instance).
* Added list of endpoints to [create Service Binding response body](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#body-9).
* Added mechanism for [orphan mitigation](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#orphan-mitigation).
* Allow brokers to return 200 for no-op [update Service Instance requests](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#response-5).
* Allow brokers to not return parameters when returning a [Service Instance](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#body-5)
  or [Service Binding](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#body-5).
* Add plan_updateable field to the [Service Plan object](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#service-plan-object). 
* [Clarify](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#blocking-operations) what happens when deleting during a provision/bind request.
* Restrict Operation strings to 10,000 chartacters in the response body for [provisioning](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#body-4)
  or [deprovisioning](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#body-12)
  a Service Instance, and [binding](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#body-9)
  or [unbinding](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#body-11)
  a Service Binding.
* Remove character restrictions on names of [Service Offerings](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#service-offering-object),
  and [Service Plans](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#service-plan-object).
* Allow empty descriptions in the response body for getting the last operations of [Service Instances](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#body-1),
  and [Service Bindings](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#body-2).
* Clarify broker behavior expected when [deprovisioning while a provision request is in progress](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#deprovisioning)
  and [unbinding while an unbind request is in progress](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#unbinding).
* Clarify broker behavior when a provision request is received [during a provision request for the same instance](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#response-3)
  or when a binding request is received [during a binding request for the same binding](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#response-6).
* Added [maintenance info](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#service-plan-object) support to Service Plans.
* Added [request identity header](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#request-identity).


## [v2.14](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md)
2018-07-24

* Added GET endpoints for fetching a
  [Service Instance](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#fetching-a-service-instance)
  and
  [Service Binding](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#fetching-a-service-binding)
* Added support for asynchronous Service Bindings
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/334))
  and a new
  [last operation endpoint for Service Bindings](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#polling-last-operation-for-service-bindings)
  endpoint
* Added clarity around concurrent updates
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/300))
* Added clarity on how Platform's can clean up after a failed provision or bind
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/353))
* Added Opaque Bearer Tokens to the
  [Platform to Service Broker Authentication](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#platform-to-service-broker-authentication)
  section
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/398))
* Provided guidance for CLI-friendly names
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/425))
* Allow for uppercase characters in Service and Service Plan names
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/433))
* Clarify that extra fields in requests and responses are allowed
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/436))
* Allow an updated `dashboard_url` to be provided when updating a Service
  Instance ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/437))
* Added an [OpenAPI 2.0 implementation](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/openapi.yaml)
* Allow for periods in name fields
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/452))
* Removed the need for Platforms to perform orphan mitigation when receiving an
  `HTTP 408` response code
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/456))
* Moved the `dashboard_client` field to
  [Cloud Foundry Catalog Extensions](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/profile.md#cloud-foundry-catalog-extensions)
* Added a [compatibility matrix](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/compatibility.md)
  describing which optional features in the specification are supported by
  different Platforms
* Added clarity for returning Service Binding information via the GET endpoints
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/517))
* Added guidance for supported string lengths
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/518))
* Clarified that the `plan_updateable` field affects modifying the Service Plan,
  not parameters ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/519))
* Clarified which Service Plan ID to use for polling the last operation endpoint
  after an update ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/522))
* Clarified Platform behaviour when a dashboard URL is not returned
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/527))
* Fixed an incompatible change introduced in v2.12
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/540))
* Added clarity around the state of resources after a failure
  ([PR](https://github.com/openservicebrokerapi/servicebroker/pull/541))
* Added [Content Type](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#content-type)
  guidelines

## [v2.13](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md)
2017-09-19

* Added [`schemas`](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#schema-object)
  field to services in the catalog that brokers can use to declare the
  configuration parameters their service accepts for creating a service
  instance, updating a service instance and creating a service binding.
* Added [`context`](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#binding)
  field to request body for creating a service binding.
* Added [warning text](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#url-properties)
  about using characters outside of the "Unreserved Characters" set in IDs.
* Added information about
  [`volume_mounts`](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#volume-mounts-object)
  objects.
* `instance_id` and `binding_id` MUST be globally unique non-empty strings.
* Allow [non-BasicAuth authentication mechanisms](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#authentication).
* Added a [Getting Started](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/gettingStarted.md)
  page including sample service brokers.
* Define what a [CLI-friendly string](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#catalog-management)
  is.
* Add [service/plan metadata conventions](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/profile.md#service-metadata).
* Add [originating identity HTTP header](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#originating-identity).

## [v2.12](https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md)
2017-06-13

* Added `context` field to request body for provisioning a service instance (`PUT /v2/service_instances/:instance_id`) and updating a service instance (`PATCH /v2/service_instances/:instance_id`). Also added the [profile.md](https://github.com/openservicebrokerapi/servicebroker/blob/master/profile.md) file describing RECOMMENDED usage patterns for environment-specific extensions and variations. `context` will replace the `organization_guid` and `space_guid` request fields in future versions of the specification. In the interim, both SHOULD be used to ensure interoperability with old and new implementations.
* The specification now uses RFC2119 keywords to make it clearer when things are REQUIRED rather than OPTIONAL or RECOMMENDED.
* Request and response parameters of type string, if present, MUST be a non-empty string.
* Cleaned up text around status codes in responses, clarifying when they MUST be returned.
* Added clarity around the `app_guid` request field.
* Clarified the semantics of `plan_id` and `parameters` fields in the updating a service instance request body.
* Clarified the default value of the `plan_updateable` field.
* Clarified when `route_service_url` is REQUIRED and when brokers can return extra data in bindings.

## [v2.11](https://github.com/openservicebrokerapi/servicebroker/blob/v2.11/spec.md)
2016-11-15

* Added field `bindable` to plan objects in /v2/catalog response. This allows services to have both bindable and non-bindable plans.

## [v2.10](http://docs.pivotal.io/pivotalcf/1-9/services/api.html)
2016-08-01

* Service bind responses now include an optional field called `volume_mounts`. Backward incompatible changes to `volume_mounts` field in service bind response from experimental 2.9 format to final format.


## [v2.9](http://docs.pivotal.io/pivotalcf/1-8/services/api.html)
2016-06-14

* `last_operation` endpoint now supports `service_id` and `plan_id` as request parameters.

* A new field `operation` may now be returned by brokers in asynchronous responses for Provision, Update, Deprovision. This field enables brokers to provide an internal identifier for the operation that clients should provide back to the service broker when polling the `last_operation` endpoint.


## [v2.8](http://docs.pivotal.io/pivotalcf/1-7/services/api.html)
2015-11-8

* In support for Route Services, service broker may now return a `route_service_url` in the response for a create binding request.

* A broker must specify `requires: ["route_forwarding"]` in its catalog endpoint if it supports Route Services.

* Clients may now send a new field `bind_resource` with the bind request, under which the parameters required for the binding are found. This would include, for example, `app_guid` for an app binding and `route` for a route binding. For backwards compatibility, `app_guid` will remain a top-level key in addition to being included in the `bind_resource`.


## [v2.7](http://docs.pivotal.io/pivotalcf/1-6/services/api.html)
2015-10-08

* Added support for Asynchronous Operations. Brokers may now return a 202 Accepted in response to provision, update, or deprovision requests to indicate the requested operation is in progress.

* The parameter `accepts_incomplete=true` must be passed by the broker client with requests for provision, update, or deprovision to indicate support for an asynchronous response. The broker can then choose to execute the request synchronously or asynchronously.

* Added support for querying `last_operation` status at a new endpoint: `GET /v2/service_instances/:guid/last_operation`


## v2.6
2015-07-23

* `app_guid` field is no longer guaranteed to be included with create service binding requests

* New field `service_id` is required with update service instance requests

## v2.5
2015-06-23

* Added support for Arbitrary Parameters: service-specific configuration parameters that can be included with provision, update and bind requests

## v2.4
2014-10-31

* Added support for broker clients to change the service plan for a specified service instance using new Update Service Instance endpoint

## v2.3
2014-04-23

* Added `dashboard_client` field to /v2/catalog to enable broker client to provision OAuth client for a service dashboard

## v2.2
2014-03-31

* Added field `free` for service plan in catalog endpoint

## v2.1
2013-12-27

* New field `app_guid` is required with bind requests
