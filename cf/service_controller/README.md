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
curl -X POST -d '{"name":"test", "hostname":"localhost", "port":"8080"}' localhost:10000/v2/service_brokers/test
```

### List service brokers

```
curl localhost:10000/v2/service_brokers
```

### Get inventory of services for a service broker

```
curl localhost:10000/v2/service_brokers/test/inventory
```

### Create Service Instances of sql for our default SB

```
curl -X POST -d '{"plan":"vmsql", "parameters":{"zone":"us-central1-a"}}' localhost:10000/v2/service_brokers/test/service_instances/guestbook-mysql
```

### Create Binding of sql for our default SB

```
curl -X POST -d '{"plan":"vmsql", "parameters":{"username":"admin","password":"mypassword"}}' localhost:10000/v2/service_brokers/test/service_instances/guestbook-mysql/service_bindings/guestbook-mysql-binding
```

### Get the creds for the binding

```
curl localhost:10000/v2/service_brokers/test/service_instances/guestbook-mysql/service_bindings/guestbook-mysql-binding
```

### Run local node guestbook that talks to the controller

```
node node/service-osc.js
```

