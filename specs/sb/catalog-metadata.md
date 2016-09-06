---
title: Catalog Metadata
owner: Core Services
---

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
| name | CLI string | A short name for the service plan to be displayed in a catalog. | name | X | |
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
