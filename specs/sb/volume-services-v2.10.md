---
title: Volume Services (Experimental)
owner: Core Services
---

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
   <td>Name of the volume driver plugin which manages the device</td>
 </tr>
 <tr>
   <td>container\_dir</td>
   <td>string</td>
   <td>The directory to mount inside the application container</td>
 </tr>
 <tr>
   <td>mode</td>
   <td>string</td>
   <td><tt>"r"</tt> to mount the volume read-only, or <tt>"rw"</tt> to mount it read-write</td>
 </tr>
 <tr>
   <td>device\_type</td>
   <td>string</td>
   <td>A string specifying the type of device to mount. Currently only <tt>"shared"</tt> devices are supported.</td>
 </tr>
 <tr>
   <td>device</td>
   <td>device-object</td>
   <td>Device object containing device\_type specific details. Currently only <tt>shared_device</tt> devices are supported.</td>
 </tr>
 </tbody>
 </table>

### shared_device ###
A `shared_device` is a subtype of a device. It represents a distributed file system which can be mounted on all app instances simultaneously.
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
   <td>volume\_id</td>
   <td>string</td>
   <td>ID of the shared volume to mount on every app instance</td>
 </tr>
 <tr>
   <td>mount\_config</td>
   <td>object</td>
   <td>Configuration object to be passed to the driver when the volume is mounted (optional)</td>
 </tr>
 </tbody>
 </table>

### Example ###
<pre class="terminal">
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
</pre>
