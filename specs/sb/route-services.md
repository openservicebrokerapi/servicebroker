---
title: Route Services
owner: Core Services
---

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
