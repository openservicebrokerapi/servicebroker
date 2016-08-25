# Service Controller

The service controller is a prototype open implementation of the parts of the
Cloud Foundry Cloud Controller which pertain to service management and
consumption.

## Prototype

Note that the node portion of this prototype does not work locally, as it is
built to run in kubernetes directly.

### Clean up previous runs

```
gcloud deployment-manager deployments delete cf-i-guestbook-mysql -q
gcloud deployment-manager deployments delete cf-b-guestbook-mysql-binding -q

pkill -f service_controller
```

### Start the service (optional)

```
./service_controller &

sleep 5
```

### Create Service Broker for our default SB

```
curl -X POST -d '{"name":"test", "broker_url":"http://localhost:8000"}' localhost:10000/v2/service_brokers
```

### List service brokers

```
curl localhost:10000/v2/service_brokers
```

### Get inventory of services across all brokers

```
curl localhost:10000/v2/inventory
```

### Create Service Instances of sql for our default SB

```
curl -X POST -d '{"name": "guestbook-mysql", "service_plan_guid":"2222", "parameters":{"zone":"us-central1-a"}}' localhost:10000/v2/service_instances
```

### Create Binding of sql for our default SB

```
curl -X POST -d '{"from_service_instance_name": "guestbook-nodejs", "service_instance_guid":"<GUID>", "parameters":{"username":"admin","password":"mypassword"}}' localhost:10000/v2/service_bindings
```

### Get the creds for the binding

```
curl localhost:10000/v2/service_bindings/<GUID>
```

### Run local node guestbook that talks to the controller

```
node node/service-osc.js
```

