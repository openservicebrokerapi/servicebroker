---
title: Services
description: Any type of add-on that can be provisioned alongside your application; for example, a database or an account on a third-party SaaS provider.
owner: Core Services
---

The documentation in this section is intended for developers and operators interested in creating Managed Services for Cloud Foundry. Managed Services are defined as having been integrated with Cloud Foundry via APIs, and enable end users to provision reserved resources and credentials on demand. For documentation targeted at end users, such as how to provision services and integrate them with applications, see [Services Overview](../devguide/services/index.html).

To develop Managed Services for Cloud Foundry, you'll need a Cloud Foundry instance to test your service broker with as you are developing it. You must have admin access to your CF instance to manage service brokers and the services marketplace catalog. For local development, we recommend using [BOSH Lite](https://github.com/cloudfoundry/bosh-lite) to deploy your own local instance of Cloud Foundry.

## Table of Contents

* <a href="overview.md" class="subnav">Overview</a>
* <a href="api.md" class="subnav">Service Broker API</a>
* <a href="managing-service-brokers.md" class="subnav">Managing Service Brokers</a>
* <a href="access-control.md" class="subnav">Access Control</a>
* <a href="catalog-metadata.md" class="subnav">Catalog Metadata</a>
* <a href="dashboard-sso.md" class="subnav">Dashboard Single Sign-On</a>
* <a href="examples.md" class="subnav">Example Service Brokers</a>
* <a href="binding-credentials.md" class="subnav">Binding Credentials</a>
* <a href="app-log-streaming.md" class="subnav">Application Log Streaming</a>
* <a href="route-services.md" class="subnav">Route Services</a>
* <a href="../devguide/services/route-binding.md" class="subnav">Manage Application Requests with Route Services</a>
* <a href="supporting-multiple-cf-instances.md" class="subnav">Supporting Multiple Cloud Foundry Instances</a>
