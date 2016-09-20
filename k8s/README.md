This dir holds the Kubernetes Service Broker PoC code

  * How to run 
  
  make sure to leave enough time between commands, as there is no internal retry mechanism.

  * Start your local cluster 
  
  `hack/local-up-cluster.sh` 
  
  k8s must be at least v1.4
  
  * Start the controller 
  
  `go run service_controller.go -backend k8s`
  
  * Start up the included gobroker 
  
  `go run gobroker.go`
  
  * Use the cli! 
  
  ```
./sc
Command Line Interface for the Service Controller

Usage:
  sc [command]

Available Commands:
  brokers           manage brokers associated with a service controller
  inventory         List the available services
  service-bindings  Manage service bindings
  service-instances Manage service instances
  service-plans     manage service plans associated with a service controller
  version           Print the version of sc

Flags:
      --controller string   URL for service controller (default "http://localhost:10000")
  -h, --help                help for sc
      --timeout int         http timeout (in seconds) for interaction with service controller (default 90)

Use "sc [command] --help" for more information about a command.



./sc brokers create gobroker http://127.0.0.1:9090
entity:
  auth_username: ""
  broker_url: http://127.0.0.1:9090
  name: gobroker
metadata:
  created_at: "2016-09-20T18:36:13Z"
  guid: 8f15aabd-1e81-4ff9-9c5a-40a74c0d6a03
  url: /v2/service_brokers/8f15aabd-1e81-4ff9-9c5a-40a74c0d6a03
  
  ./sc brokers list
- AuthPassword: ""
  AuthUsername: ""
  BrokerURL: http://127.0.0.1:9090
  Created: "1474396573"
  GUID: 8f15aabd-1e81-4ff9-9c5a-40a74c0d6a03
  Name: gobroker
  SelfURL: /v2/service_brokers/8f15aabd-1e81-4ff9-9c5a-40a74c0d6a03
  Updated: 0
  
  
  
  
  ./sc inventory
service              plans                     description
myService1           freePlan,costlyPlan       very cool



./sc service-plans list
services:
- bindable: true
  dashboard_client: null
  description: very cool
  id: 12345678-abcd-1234-bcde-1234567890ab
  metadata: null
  name: myService1
  plan_updateable: false
  plans:
  - description: free is good
    free: true
    id: ffffffff-0000-0000-0000-000000000000
    metadata: null
    name: freePlan
    schemas:
      binding:
        inputs: ""
        outputs: ""
      instance:
        inputs: ""
        outputs: ""
  - description: not so free
    free: false
    id: eeeeeeee-1111-1111-1111-111111111111
    metadata: null
    name: costlyPlan
    schemas:
      binding:
        inputs: ""
        outputs: ""
      instance:
        inputs: ""
        outputs: ""
  requires: null
  tags: null
  
  
  
  ./sc service-instances create inst myService1 freePlan
Found Service Plan GUID as ffffffff-0000-0000-0000-000000000000 for myService1 : freePlanSending body: {"name":"inst","service_plan_guid":"ffffffff-0000-0000-0000-000000000000","space_guid":"default","parameters":null,"tags":null}

credentials: ""
dashboard_url: http://example.com/dashAwayAll
id: 19097ac0-f9ce-424e-936d-60c7fb3d231e
last_operation:
  description: ""
  updated_at: ""
name: inst
parameters: null
routes_url: ""
service_id: 12345678-abcd-1234-bcde-1234567890ab
service_plan_guid: ffffffff-0000-0000-0000-000000000000
service_plan_url: ""
space_guid: default
space_url: ""
tags: null
type: managed_service_instance



./sc service-instances list
- credentials: ""
  dashboard_url: http://example.com/dashAwayAll
  id: 19097ac0-f9ce-424e-936d-60c7fb3d231e
  last_operation:
    description: ""
    updated_at: ""
  name: inst
  parameters: null
  routes_url: ""
  service_id: 12345678-abcd-1234-bcde-1234567890ab
  service_plan_guid: ffffffff-0000-0000-0000-000000000000
  service_plan_url: ""
  space_guid: default
  space_url: ""
  tags: null
  type: managed_service_instance



  ./sc service-bindings list 
  # no bindings, let's make one
  
  ./sc service-bindings create myapp inst
Checking: &model.ServiceInstance{Name:"inst", Credentials:"", ServicePlanGUID:"ffffffff-0000-0000-0000-000000000000", SpaceGUID:"default", DashboardURL:"http://example.com/dashAwayAll", Type:"managed_service_instance", LastOperation:(*model.LastOperation1)(0xc42010a140), SpaceURL:"", ServicePlanURL:"", RoutesURL:"", Tags:[]string(nil), Parameters:interface {}(nil), ID:"19097ac0-f9ce-424e-936d-60c7fb3d231e", ServiceID:"12345678-abcd-1234-bcde-1234567890ab"}
app_name: myapp
credentials:
  password: letmein
id: f567e056-2ec8-48ae-83d4-e2528a9675e1
service_instance_guid: 19097ac0-f9ce-424e-936d-60c7fb3d231e
service_instance_name: inst


./sc service-bindings list
Instance  AppName  Credentials
inst      myapp    map[password:letmein]

  ```
  
  
