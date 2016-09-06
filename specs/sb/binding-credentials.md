---
title: Binding Credentials
owner: Core Services
---

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
