#!/bin/bash

# Note: this requires the fix from:
#  https://github.com/kubernetes/kubernetes/pull/24299

(
# Clean-up from previous run
kubectl delete rc app
kubectl delete secret creds
kubectl delete ServiceBroker mongodb-sb
kubectl delete ServiceInstance mongodb-instance1
kubectl delete ServiceBinding mongodb-instance1-binding1
kubectl delete thirdpartyresource service-broker.cncf.org
kubectl delete thirdpartyresource service-instance.cncf.org
kubectl delete thirdpartyresource service-binding.cncf.org
) > /dev/null 2>&1

kubectl get thirdpartyresources

set -ex

# First create the new resource types
kubectl create -f - <<-EOF
    apiVersion: extensions/v1beta1
    kind: ThirdPartyResource
    metadata:
      name: service-broker.cncf.org
    versions:
    - name: v1
EOF

kubectl create -f - <<-EOF
    apiVersion: extensions/v1beta1
    kind: ThirdPartyResource
    metadata:
      name: service-instance.cncf.org
    versions:
    - name: v1
EOF

kubectl create -f - <<-EOF
    apiVersion: extensions/v1beta1
    kind: ThirdPartyResource
    metadata:
      name: service-binding.cncf.org
    versions:
    - name: v1
EOF

sleep 10

# Verify everything is there
curl http://localhost:8080/apis/cncf.org/v1/namespaces/default/servicebrokers
curl http://localhost:8080/apis/cncf.org/v1/namespaces/default/serviceinstances
curl http://localhost:8080/apis/cncf.org/v1/namespaces/default/servicebindings

# New create instances of each type
kubectl create -f - <<-EOF
    apiVersion: cncf.org/v1
    kind: ServiceBroker
    metadata:
      name: mongodb-sb
EOF

kubectl create -f - <<-EOF
    apiVersion: cncf.org/v1
    kind: ServiceInstance
    metadata:
      name: mongodb-instance1
EOF

kubectl create -f - <<-EOF
    apiVersion: cncf.org/v1
    kind: ServiceBinding
    metadata:
      name: mongodb-instance1-binding1
    creds:
      user: john
      password: letMeIn
EOF

# And the secret to hold our creds
kubectl create -f - <<EOF
    apiVersion: v1
    kind: Secret
    metadata:
      name: creds
    data:
      vcap-services: eyAidXNlciI6ICJqb2huIiwgInBhc3N3b3JkIjogImxldE1lSW4iIH0=
    type: myspecial/secret
EOF

sleep 10

# Now create the app with VCAP_SERVICES binding info
kubectl create -f - <<EOF
    apiVersion: v1
    kind: ReplicationController
    metadata:
      name: app
    spec:
      replicas: 1
      selector:
        version: v1.0
      template:
        metadata:
          name: myserver
          labels:
            version: v1.0
        spec:
          containers:
          - name: nginx
            image: nginx
            env:
              - name: VCAP_SERVICES
                valueFrom:
                  secretKeyRef:
                    name: creds
                    key: vcap-services
EOF

sleep 15

# Prove it worked
kubectl exec `kubectl get pods --template "{{ (index .items 0).metadata.name }}"` -- env
