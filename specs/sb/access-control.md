---
title: Access Control
owner: Core Services
---

All new service plans from standard private brokers are private by default. This means that when adding a new broker, or when adding a new plan to an existing broker's catalog, service plans won't immediately be available to end users. This lets an admin control which service plans are available to end users, and manage limited service availability.

Space-scoped private brokers are registered to a specific space, and all users within that space can automatically access the broker's service plans. With space-scoped brokers, service visibility is not managed separately.

## <a id='cli'></a>Using the CLI ##

If your CLI and/or deployment of cf-release do not meet the following prerequisites, you can manage access control with [cf curl](#curl).

### <a id='prerequisites'></a>Prerequisites ###
- CLI v6.4.0
- Cloud Controller API v2.9.0 (cf-release v179)
- Admin user access; the following commands can be run only by an admin user

To determine your API version, curl `/v2/info` and look for `api_version`.

<pre class="terminal">
$ cf curl /v2/info
{
   "name": "vcap",
   "build": "2222",
   "support": "http://support.cloudfoundry.com",
   "version": 2,
   "description": "Cloud Foundry sponsored by Pivotal",
   "authorization_endpoint": "https://login.system-domain.example.com",
   "token_endpoint": "https://uaa.system-domain.example.com",
   "api_version": "2.13.0",
   "logging_endpoint": "wss://loggregator.system-domain.example.com:443"
}
</pre>

### <a id='display-access'></a>Display Access to Service Plans ###

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

### <a id='enable-access'></a>Enable Access to Service Plans ###

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

### <a id='disable-access'></a>Disable Access to Service Plans ###

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

#### Limitations ####

- You cannot disable access to a service plan for an organization if the plan is currently available to all organizations. You must first disable access for all organizations; then you can enable access for a  particular organization.

## <a id='curl'></a>Using cf curl ##

The following commands must be run as a system admin user.

### <a id='enable-access-curl'></a>Enable Access to Service Plans ###

Access can be enabled for users of all organizations, or for users of particular organizations. Service plans which are available to all users are said to be "public". Plans that are available to no organizations, or to particular organizations, are said to be "private".

#### Enable access to a plan for all organizations ####

Once made public, the service plan can be seen by all users in the list of available services. See [Managing Services](../devguide/services/managing-services.html) for more information.

To make a service plan public, you need the service plan GUID. To find the service plan GUID, run:

`cf curl /v2/service_plans -X 'GET'`

This command returns a filtered JSON response listing every service plan. Data about each plan shows in two sections: `metadata` and `entity.` The `metadata` section shows the service plan GUID, while the `entity` section lists the name of the plan. Note: Because `metadata` is listed before `entity` for each service plan, the GUID of a plan is shown six lines above the name.

Example:

<pre class="terminal">
$ cf curl /v2/service_plans
...
{
    "metadata": {
        "guid": "1afd5050-664e-4be2-9389-6bf0c967c0c6",
        "url": "/v2/service_plans/1afd5050-664e-4be2-9389-6bf0c967c0c6",
        "created_at": "2014-02-12T06:24:04+00:00",
        "updated_at": "2014-02-12T18:46:52+00:00"
    },
    "entity": {
        "name": "plan-name-1",
        "free": true,
        "description": "plan-desc-1",
        "service_guid": "d9011411-1463-477c-b223-82e04996b91f",
        "extra": "{\"bullets\":[\"bullet1\",\"bullet2\"]}",
        "unique_id": "plan-id-1",
        "public": false,
        "service_url": "/v2/services/d9011411-1463-477c-b223-82e04996b91f",
        "service_instances_url": "/v2/service_plans/1afd5050-664e-4be2-9389-6bf0c967c0c6/service_instances"
    }
}
</pre>

In this example, the GUID of plan-name-1 is 1afd5050-664e-4be2-9389-6bf0c967c0c6.

To make a service plan public, run:
`cf curl /v2/service_plans/SERVICE_PLAN_GUID -X 'PUT' -d '{"public":true}'`

As verification, the "entity" section of the JSON response shows the `"public":true` key-value pair.

<pre class="terminal">
$ cf curl /v2/service_plans/1113aa0-124e-4af2-1526-6bfacf61b111 -X 'PUT' -d '{"public":true}'

{
    "metadata": {
        "guid": "1113aa0-124e-4af2-1526-6bfacf61b111",
        "url": "/v2/service_plans/1113aa0-124e-4af2-1526-6bfacf61b111",
        "created_at": "2014-02-12T06:24:04+00:00",
        "updated_at": "2014-02-12T20:55:10+00:00"
    },
    "entity": {
        "name": "plan-name-1",
        "free": true,
        "description": "plan-desc-1",
        "service_guid": "d9011411-1463-477c-b223-82e04996b91f",
        "extra": "{\"bullets\":[\"bullet1\",\"bullet2\"]}",
        "unique_id": "plan-id-1",
        "public": true,
        "service_url": "/v2/services/d9011411-1463-477c-b223-82e04996b91f",
        "service_instances_url": "/v2/service_plans/1113aa0-124e-4af2-1526-6bfacf61b111/service_instances"
    }
}
</pre>

#### Enable access to a private plan for a particular organization ####

Users have access to private plans that have been enabled for an organization only when targeting a space of  that organization. See [Managing Services](../devguide/services/managing-services.html) for more information.

To make a service plan available to users of a specific organization, you need the GUID of both the organization and the service plan. To get the GUID of the service plan, run the same command described above for [enabling access to a plan for all organizations](#enable-access-curl):

`cf curl -X 'GET' /v2/service_plans`

To find the organization GUIDs, run:

`cf curl /v2/organizations?q=name:YOUR-ORG-NAME`

The `metadata` section shows the organization GUID, while the `entity` section lists the name of the organization. Note: Because `metadata` is listed before `entity` for each organization, the GUID of an organization is shown six lines above the name.

Example:

<pre class="terminal">
$ cf curl /v2/organizations?q=name:my-org

{
    "metadata": {
        "guid": "c54bf317-d791-4d12-89f0-b56d0936cfdc",
        "url": "/v2/organizations/c54bf317-d791-4d12-89f0-b56d0936cfdc",
        "created_at": "2013-05-06T16:34:56+00:00",
        "updated_at": "2013-09-25T18:44:35+00:00"
    },
    "entity": {
        "name": "my-org",
        "billing_enabled": true,
        "quota_definition_guid": "52c5413c-869f-455a-8873-7972ecb85ca8",
        "status": "active",
        "quota_definition_url": "/v2/quota_definitions/52c5413c-869f-455a-8873-7972ecb85ca8",
        "spaces_url": "/v2/organizations/c54bf317-d791-4d12-89f0-b56d0936cfdc/spaces",
        "domains_url": "/v2/organizations/c54bf317-d791-4d12-89f0-b56d0936cfdc/domains",
        "private_domains_url": "/v2/organizations/c54bf317-d791-4d12-89f0-b56d0936cfdc/private_domains",
        "users_url": "/v2/organizations/c54bf317-d791-4d12-89f0-b56d0936cfdc/users",
        "managers_url": "/v2/organizations/c54bf317-d791-4d12-89f0-b56d0936cfdc/managers",
        "billing_managers_url": "/v2/organizations/c54bf317-d791-4d12-89f0-b56d0936cfdc/billing_managers",
        "auditors_url": "/v2/organizations/c54bf317-d791-4d12-89f0-b56d0936cfdc/auditors",
        "app_events_url": "/v2/organizations/c54bf317-d791-4d12-89f0-b56d0936cfdc/app_events"
    }
}
</pre>

In this example, the GUID of my-org is c54bf317-d791-4d12-89f0-b56d0936cfdc.

To make a private plan available to a specific organization, run:

`cf curl /v2/service_plan_visibilities -X POST -d '{"service_plan_guid":"SERVICE_PLAN_GUID","organization_guid":"ORG_GUID"}'`

Example:

<pre class="terminal">
$ cf curl /v2/service_plan_visibilities -X 'POST' -d '{"service_plan_guid":"1113aa0-124e-4af2-1526-6bfacf61b111","organization_guid":"aaaa1234-da91-4f12-8ffa-b51d0336aaaa"}'

{
    "metadata": {
        "guid": "99993789-a368-483e-ae7c-ebe79e199999",
        "url": "/v2/service_plan_visibilities/99993789-a368-483e-ae7c-ebe79e199999",
        "created_at": "2014-02-12T21:03:42+00:00",
        "updated_at": null
    },
    "entity": {
        "service_plan_guid": "1113aa0-124e-4af2-1526-6bfacf61b111",
        "service_plan_url": "/v2/service_plans/1113aa0-124e-4af2-1526-6bfacf61b111",
        "context": {
            "dept": "abc123"
		}
    }
}
</pre>

Members of my-org can now see the plan-name-1 service plan in the list of available services when a space of my-org is targeted.

Note: The `guid` field in the `metadata` section of this JSON response is the id of the "service plan visibility", and can be used to revoke access to the plan for the organization as described below.

### <a id='disable-access-curl'></a>Disable Access to Service Plans ###

#### Disable access to a plan for all organizations

To make a service plan private, follow the instructions above for [Enable Access](#enable-access-curl), but replace `"public":true` with `"public":false`.

Note: organizations that have explicitly been granted access will retain access once a plan is private. To be sure access is removed for all organizations, access must be explicitly revoked for organizations to which access has been explicitly granted. For details see below.

Example making plan-name-1 private:

<pre class="terminal">
$ cf curl /v2/service_plans/1113aa0-124e-4af2-1526-6bfacf61b111 -X 'PUT' -d '{"public":false}'

{
    "metadata": {
        "guid": "1113aa0-124e-4af2-1526-6bfacf61b111",
        "url": "/v2/service_plans/1113aa0-124e-4af2-1526-6bfacf61b111",
        "created_at": "2014-02-12T06:24:04+00:00",
        "updated_at": "2014-02-12T20:55:10+00:00"
    },
    "entity": {
        "name": "plan-name-1",
        "free": true,
        "description": "plan-desc-1",
        "service_guid": "d9011411-1463-477c-b223-82e04996b91f",
        "extra": "{\"bullets\":[\"bullet1\",\"bullet2\"]}",
        "unique_id": "plan-id-1",
        "public": false,
        "service_url": "/v2/services/d9011411-1463-477c-b223-82e04996b91f",
        "service_instances_url": "/v2/service_plans/1113aa0-124e-4af2-1526-6bfacf61b111/service_instances"
    }
}
</pre>

#### Disable access to a private plan for a particular organization ####

To revoke access to a service plan for a particular organization, run:

`cf curl /v2/service_plan_visibilities/SERVICE_PLAN_VISIBILITIES_GUID -X 'DELETE'`

Example:

<pre class="terminal">
$ cf curl /v2/service_plan_visibilities/99993789-a368-483e-ae7c-ebe79e199999 -X DELETE
</pre>
