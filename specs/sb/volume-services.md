---
title: Volume Services (Experimental/Obsolete)
owner: Core Services
---

## <a id='introduction'></a>Introduction ##

Cloud Foundry application developers may want their applications to mount one or more volumes in order to write to a reliable, non-ephemeral file system. By integrating with service brokers and the Cloud Foundry runtime, providers can offer these services to developers through an automated, self-service, and on-demand user experience.

<p class="note"><strong>Note</strong>: This feature is experimental.</p>

<p class="note"><strong>Note</strong>: The v2.9 version of this experimental feature is no longer supported as of Cloud Foundry v240. If you are using a newer Cloud Foundry version and your service broker returns volume mounts, you must update your service broker to use the new <a href="volume-services-v2.10.md">service broker v2.10 format for volume mounts</a>.</p>

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
   <td>container_path</td>
   <td>string</td>
   <td>The path to mount inside the application container</td>
 </tr>
 <tr>
   <td>mode</td>
   <td>string</td>
   <td>Indicates whether the volume should be read-only or read-write</td>
 </tr>
 <tr>
   <td>private</td>
   <td>private</td>
   <td>A <i>private</i> JSON object</td>
 </tr>
 </tbody>
 </table>

### private ###
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
   <td>driver</td>
   <td>string</td>
   <td>The path to mount inside the application container</td>
 </tr>
 <tr>
   <td>group_id</td>
   <td>string</td>
   <td>Indicates whether the volume should be read-only or read-write</td>
 </tr>
 <tr>
   <td>config</td>
   <td>string</td>
   <td>A configuration string associated with a particular broker/driver pair.  The broker is free to return any string here and it will be passed through the volume driver identified in the `driver` field</td>
 </tr>
 </tbody>
 </table>

### Example ###
<pre class="terminal">
{
  ...
  "volume_mounts": [
    {
      "container_path": "/data/images",
      "mode": "r",
      "private": {
        "driver": "cephdriver",
        "group_id": "bc2c1eab-05b9-482d-b0cf-750ee07de311",
        "config": "Some arbitrary configuration string. Could be a marshalled json"
      }
    }
  ]
}
</pre>
