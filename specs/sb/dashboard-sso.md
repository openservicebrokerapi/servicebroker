---
title: Dashboard Single Sign-On
owner: Core Services
---

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
