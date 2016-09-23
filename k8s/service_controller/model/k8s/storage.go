package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/servicebroker/servicebroker/k8s/service_controller/model"
)

var _ model.ServiceStorage = (*K8sServiceStorage)(nil)

type K8sServiceStorage struct {
	// Host is the location where we'll talk to k8s
	host            string
	defaultResource string
}

const serviceDomain string = "cncf.org"
const apiVersion string = "v1alpha1"
const brokerResource string = "servicebrokers"
const defaultUri string = "http://%v/apis/" + serviceDomain + "/" + apiVersion + "/namespaces/default/" + brokerResource

// The k8s implementation should leverage Third Party Resources
// https://github.com/kubernetes/kubernetes/blob/master/docs/design/extending-api.md

var _ model.ServiceStorage = (*K8sServiceStorage)(nil)

type Meta struct {
	Name string `json:"name"`
}

type KubeData struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   Meta   `json:"metadata"`
}

type k8sServiceBroker struct {
	*model.ServiceBroker
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   Meta   `json:"metadata"`
}

type k8sService struct {
	*model.Service
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   Meta   `json:"metadata"`
}

type k8sPlan struct {
	*model.ServicePlan
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   Meta   `json:"metadata"`
}

type VName struct {
	Name string `json:"name"`
}

type TPR struct {
	Meta       `json:"metadata"`
	ApiVersion string  `json:"apiVersion"`
	kind       string  `json:"kind"`
	Versions   []VName `json:"versions"`
}

const TPRapiVersion string = "extensions/v1beta1"
const thirdPartyResourceString string = "ThirdPartyResource"

var versionMap []VName = []VName{{apiVersion}}

// Kubernetes ThirdPartyResources definitions
var serviceBrokerDefinition TPR = TPR{Meta{"service-broker.cncf.org"},
	TPRapiVersion, thirdPartyResourceString, versionMap}

// sbservice so it does not conflict with the built in Service
var serviceDefinition TPR = TPR{Meta{"sbservice.cncf.org"},
	TPRapiVersion, thirdPartyResourceString, versionMap}

const serviceResource string = "sbservices"
const defaultServiceFormatUri string = "http://%v/apis/" + serviceDomain + "/" + apiVersion + "/namespaces/default/" + serviceResource

var servicePlanDefinition TPR = TPR{Meta{"service-plan.cncf.org"},
	TPRapiVersion, thirdPartyResourceString, versionMap}

const servicePlanResource string = "serviceplans"
const defaultServicePlanFormatUri string = "http://%v/apis/" + serviceDomain + "/" + apiVersion + "/namespaces/default/" + servicePlanResource

var serviceInstanceDefinition TPR = TPR{Meta{"service-instance.cncf.org"},
	TPRapiVersion, thirdPartyResourceString, versionMap}

const serviceInstanceResource string = "serviceinstances"
const defaultServiceInstanceFormatUri string = "http://%v/apis/" + serviceDomain + "/" + apiVersion + "/namespaces/default/" + serviceInstanceResource

var serviceBindingDefinition TPR = TPR{Meta{"service-binding.cncf.org"},
	TPRapiVersion, thirdPartyResourceString, versionMap}

const serviceBindingResource string = "servicebindings"
const defaultServiceBindingFormatUri string = "http://%v/apis/" + serviceDomain + "/" + apiVersion + "/namespaces/default/" + serviceBindingResource

func CreateServiceStorage(host string) model.ServiceStorage {
	k := &K8sServiceStorage{host: host,
		defaultResource: fmt.Sprintf(defaultUri, host)}
	fmt.Println(" root host is:", k.defaultUri())
	// define the resources once at startup
	// results in ServiceBrokers

	k.createTPR(serviceBrokerDefinition)
	k.createTPR(serviceDefinition)
	k.createTPR(servicePlanDefinition)
	k.createTPR(serviceBindingDefinition)
	k.createTPR(serviceInstanceDefinition)
	// cleanup afterwards by `kubectl delete thirdpartyresource service-broker.cncf.org`

	return k
}

// listSB is only used for unmarshalling the list of service brokers
// for returning to the client
type listSB struct {
	Items []*k8sServiceBroker `json:"items"`
}

func (kss *K8sServiceStorage) defaultUri() string {
	return kss.defaultResource
}

func (kss *K8sServiceStorage) defaultServiceUri() string {
	return fmt.Sprintf(defaultServiceFormatUri, kss.host)
}

func (kss *K8sServiceStorage) defaultPlanUri() string {
	return fmt.Sprintf(defaultServicePlanFormatUri, kss.host)
}

func (kss *K8sServiceStorage) defaultServiceInstanceUri() string {
	return fmt.Sprintf(defaultServiceInstanceFormatUri, kss.host)
}

func (kss *K8sServiceStorage) defaultServiceBindingUri() string {
	return fmt.Sprintf(defaultServiceBindingFormatUri, kss.host)
}

func (kss *K8sServiceStorage) createTPR(tpr TPR) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&tpr)
	fmt.Printf("encoded bytes: %v\n", b.String())
	r, e := http.Post("http://"+kss.host+"/apis/extensions/v1beta1/thirdpartyresources", "application/json", b)
	fmt.Printf("result: %v\n", r)
	if nil != e || 201 != r.StatusCode {
		fmt.Printf("Error creating k8s TPR [%s]...\n%v\n", e, r)
	}
}

/* BROKER */
/**********/

func (kss *K8sServiceStorage) ListBrokers() ([]string, error) {
	fmt.Println("listing all brokers")
	// get the ServiceBroker

	r, e := http.Get(kss.defaultUri())
	if nil != e {
		return nil, fmt.Errorf("couldn't get the service brokers. %v, [%v]", e, r)
	}

	var lsb listSB
	e = json.NewDecoder(r.Body).Decode(&lsb)
	if nil != e { // wrong json format error
		fmt.Println("json not unmarshalled:", e, r)
		return nil, e
	}
	fmt.Println("Got", len(lsb.Items), "brokers.")
	ret := make([]string, 0, len(lsb.Items))
	for _, v := range lsb.Items {
		ret = append(ret, v.ServiceBroker.ID)
	}
	return ret, nil
}

func (kss *K8sServiceStorage) AddBroker(broker *model.ServiceBroker) error {
	fmt.Println("adding broker to k8s", broker)
	// create TPR
	// tpr is
	//    kind.fqdn
	// or
	//    kind.domain.tld
	//
	// use service-broker.cncf.org
	// end up with k8s resource of ServiceBroker
	// version v1alpha1 for now
	//
	// store name/host/port/user/pass as metadata
	//
	// example yaml
	// metadata:
	//   name: service-broker.cncf.org
	//   (service)name/host/port/user/pass
	// apiVersion: extensions/v1beta1
	// kind: ThirdPartyResource
	// versions:
	// - name: v1alpha1
	ksb := NewK8sSB()
	ksb.Metadata = Meta{Name: broker.ID}
	ksb.ServiceBroker = broker

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&ksb)
	fmt.Printf("sending: %v", b)
	r, e := http.Post(kss.defaultUri(), "application/json", b)
	fmt.Sprintf("result: %v", r)
	if nil != e || 201 != r.StatusCode {
		fmt.Printf("Error creating k8s service broker TPR [%s]...\n%v\n", e, r)
		return e
	}

	return nil
}

func (kss *K8sServiceStorage) GetBroker(name string) (*model.ServiceBroker, error) {
	uri := kss.defaultUri() + "/" + name
	fmt.Println("uri is:", uri)
	r, e := http.Get(uri)
	if nil != e {
		return nil, fmt.Errorf("couldn't get the service broker. %v, [%v]", e, r)
	}
	defer r.Body.Close()
	var sb k8sServiceBroker
	e = json.NewDecoder(r.Body).Decode(&sb)
	if nil != e { // wrong json format error
		return nil, e
	}
	fmt.Printf("returned json: %+v\n", sb)
	return sb.ServiceBroker, nil
}

func (kss *K8sServiceStorage) SetBroker(si *model.ServiceBroker) error {
	return fmt.Errorf("SetBroker: Not implemented yet")
}

func (kss *K8sServiceStorage) DeleteBroker(name string) error {
	return deleteResource(kss.defaultUri() + "/" + name)
}

func deleteResource(resource string) error {
	fmt.Println("resource is:", resource)

	// utter failure of an http API
	req, _ := http.NewRequest("DELETE", resource, nil)
	_, e := http.DefaultClient.Do(req)
	return e
}

func NewK8sSB() *k8sServiceBroker {
	return &k8sServiceBroker{ApiVersion: serviceDomain + "/" + apiVersion,
		Kind: "ServiceBroker"}
}

/* Service */
/***********/

type listS struct {
	Items []*k8sService `json:"items"`
}

func NewK8sService() *k8sService {
	return &k8sService{ApiVersion: serviceDomain + "/" + apiVersion,
		Kind: "Sbservice"}
}

func (kss *K8sServiceStorage) ListServices() ([]string, error) {
	fmt.Println("listing all services")
	r, e := http.Get(kss.defaultServiceUri())
	if nil != e {
		return nil, fmt.Errorf("couldn't get the services. %v, [%v]", e, r)
	}

	var ls listS
	e = json.NewDecoder(r.Body).Decode(&ls)
	if nil != e { // wrong json format error
		fmt.Println("json not unmarshalled:", e, r)
		return nil, e
	}
	fmt.Println("Got", len(ls.Items), "services.")
	ret := make([]string, 0, len(ls.Items))
	for i, v := range ls.Items {
		fmt.Println("service", i, v)
		ret = append(ret, v.Service.ID)
	}
	return ret, nil
}

func (s *K8sServiceStorage) GetServices() ([]*model.Service, error) {
	return nil, fmt.Errorf("GetServices: Not implemented yet")
}

func (kss *K8sServiceStorage) GetService(id string) (*model.Service, error) {
	fmt.Println("getting a single service")
	r, e := http.Get(kss.defaultServiceUri() + "/" + id)
	if nil != e {
		return nil, fmt.Errorf("couldn't get the service. %v, [%v]", e, r)
	}

	var s k8sService
	e = json.NewDecoder(r.Body).Decode(&s)
	if nil != e { // wrong json format error
		fmt.Println("json not unmarshalled:", e, r)
		return nil, e
	}
	fmt.Println("Got a service!", s)

	return s.Service, nil
}

func (kss *K8sServiceStorage) AddService(si *model.Service) error {
	fmt.Println(si, si.ID)

	ks := NewK8sService()
	ks.Metadata = Meta{Name: si.ID}
	ks.Service = si

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(&ks); nil != err {
		fmt.Println("failed to encode", si, "as", ks)
		return err
	}
	defaultUri := "http://%v/apis/" + serviceDomain + "/" + apiVersion + "/namespaces/default/" + "sbservices"
	fmt.Printf("sending: %v\n to %v", b, defaultUri)
	r, e := http.Post(fmt.Sprintf(defaultUri, kss.host), "application/json", b)
	fmt.Sprintf("result: %v", r)
	if nil != e || 201 != r.StatusCode {
		fmt.Printf("Error creating k8s service TPR [%s]...\n%v\n", e, r)
		return e
	}
	return nil
}

func (kss *K8sServiceStorage) SetService(si *model.Service) error {
	fmt.Println(si, si.ID)

	ks := NewK8sService()
	ks.Metadata = Meta{Name: si.ID}
	ks.Service = si

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(&ks); nil != err {
		fmt.Println("failed to encode", si, "as", ks)
		return err
	}
	defaultUri := kss.defaultServiceUri() + "/" + si.ID
	fmt.Printf("sending: %v\n to %v", b, defaultUri)
	req, _ := http.NewRequest("PUT", defaultUri, b)
	req.Header.Set("Content-Type", "application/json")
	r, e := http.DefaultClient.Do(req)
	fmt.Sprintf("result: %v", r)
	if nil != e || 201 != r.StatusCode {
		fmt.Printf("Error updating k8s service TPR [%s]...\n%v\n", e, r)
		return e
	}
	return nil
}

func (kss *K8sServiceStorage) DeleteService(id string) error {
	return deleteResource(kss.defaultServiceUri() + "/" + id)
}

/* Plan */
/********/

func NewK8sPlan() *k8sPlan {
	return &k8sPlan{ApiVersion: serviceDomain + "/" + apiVersion,
		Kind: "ServicePlan"}
}

func (kss *K8sServiceStorage) ListPlans() ([]string, error) {
	return nil, fmt.Errorf("ListPlans: Not implemented yet")
}

func (kss *K8sServiceStorage) GetPlans() ([]*model.ServicePlan, error) {
	return nil, fmt.Errorf("GetPlans: Not implemented yet")
}

func (kss *K8sServiceStorage) GetPlan(id string) (*model.ServicePlan, error) {
	fmt.Println("getting a single plan")
	r, e := http.Get(kss.defaultPlanUri() + "/" + id)
	if nil != e {
		return nil, fmt.Errorf("couldn't get the plan. %v, [%v]", e, r)
	}

	var s k8sPlan
	e = json.NewDecoder(r.Body).Decode(&s)
	if nil != e { // wrong json format error
		fmt.Println("json not unmarshalled:", e, r)
		return nil, e
	}
	fmt.Println("Got a plan!", r.Body, s)

	return s.ServicePlan, nil
}

func (kss *K8sServiceStorage) AddPlan(plan *model.ServicePlan) error {
	fmt.Println(plan, plan.ID)

	ks := NewK8sPlan()
	ks.Metadata = Meta{Name: plan.ID}
	ks.ServicePlan = plan

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(&ks); nil != err {
		fmt.Println("failed to encode", plan, "as", ks)
		return err
	}
	defaultUri := "http://%v/apis/" + serviceDomain + "/" + apiVersion + "/namespaces/default/" + "serviceplans"
	fmt.Printf("sending: %v\n to %v\n", b, defaultUri)
	r, e := http.Post(fmt.Sprintf(defaultUri, kss.host), "application/json", b)
	fmt.Sprintf("result: %v", r)
	if nil != e || 201 != r.StatusCode {
		fmt.Printf("Error creating k8s service TPR [%s]...\n%v\n", e, r)
		return e
	}
	return nil
}

func (kss *K8sServiceStorage) SetPlan(si *model.ServicePlan) error {
	return fmt.Errorf("SetPlan: Not implemented yet")
}

func (kss *K8sServiceStorage) DeletePlan(id string) error {
	return deleteResource(kss.defaultPlanUri() + "/" + id)
}

/* Instance */
/************/

type listI struct {
	Items []*k8sServiceInstance `json:"items"`
}

func NewK8sServiceInstance() *k8sServiceInstance {
	return &k8sServiceInstance{ApiVersion: serviceDomain + "/" + apiVersion,
		Kind: "ServiceInstance"}
}

type k8sServiceInstance struct {
	*model.ServiceInstance
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   Meta   `json:"metadata"`
}

func (kss *K8sServiceStorage) ListInstances() ([]string, error) {
	fmt.Println("listing all service instances")
	r, e := http.Get(kss.defaultServiceInstanceUri())
	if nil != e {
		return nil, fmt.Errorf("couldn't get the services. %v, [%v]", e, r)
	}

	var ls listI
	e = json.NewDecoder(r.Body).Decode(&ls)
	if nil != e { // wrong json format error
		fmt.Println("json not unmarshalled:", e, r)
		return nil, e
	}
	fmt.Println("Got", len(ls.Items), "service instances.")
	ret := make([]string, 0, len(ls.Items))
	for i, v := range ls.Items {
		fmt.Println("instance", i, v)
		ret = append(ret, v.ServiceInstance.ID)
	}
	return ret, nil
}

func (s *K8sServiceStorage) GetInstances() ([]*model.Service, error) {
	return nil, fmt.Errorf("GetInstances: Not implemented yet")
}

func (kss *K8sServiceStorage) GetInstance(id string) (*model.ServiceInstance, error) {
	fmt.Println("getting a single service instance")
	r, e := http.Get(kss.defaultServiceInstanceUri() + "/" + id)
	if nil != e {
		return nil, fmt.Errorf("couldn't get the service instance. %v, [%v]", e, r)
	}

	var s k8sServiceInstance
	e = json.NewDecoder(r.Body).Decode(&s)
	if nil != e { // wrong json format error
		fmt.Println("json not unmarshalled:", e, r)
		return nil, e
	}
	fmt.Println("Got a service instance!", r.Body, s)

	return s.ServiceInstance, nil
}

func (kss *K8sServiceStorage) AddInstance(si *model.ServiceInstance) error {
	fmt.Println(si, si.ID)

	ks := NewK8sServiceInstance()
	ks.Metadata = Meta{Name: si.ID}
	ks.ServiceInstance = si

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(&ks); nil != err {
		fmt.Println("failed to encode", si, "as", ks)
		return err
	}
	defaultUri := kss.defaultServiceInstanceUri()
	fmt.Printf("sending: %v\n to %v\n", b, defaultUri)
	r, e := http.Post(defaultUri, "application/json", b)
	fmt.Sprintf("result: %v", r)
	if nil != e || 201 != r.StatusCode {
		fmt.Printf("Error creating k8s service instance TPR [%s]...\n%v\n", e, r)
		return e
	}
	return nil
}

func (kss *K8sServiceStorage) SetInstance(si *model.ServiceInstance) error {
	fmt.Println(si, si.ID)

	ks := NewK8sServiceInstance()
	ks.Metadata = Meta{Name: si.ID}
	ks.ServiceInstance = si

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(&ks); nil != err {
		fmt.Println("failed to encode", si, "as", ks)
		return err
	}
	defaultUri := kss.defaultServiceInstanceUri() + "/" + si.ID
	fmt.Printf("sending: %v\n to %v", b, defaultUri)
	req, _ := http.NewRequest("PUT", defaultUri, b)
	req.Header.Set("Content-Type", "application/json")
	r, e := http.DefaultClient.Do(req)
	fmt.Sprintf("result: %v", r)
	if nil != e || 201 != r.StatusCode {
		fmt.Printf("Error updating k8s service instance TPR [%s]...\n%v\n", e, r)
		return e
	}
	return nil
}

func (kss *K8sServiceStorage) DeleteInstance(id string) error {
	return deleteResource(kss.defaultServiceInstanceUri() + "/" + id)
}

/* Binding */
/***********/

type listB struct {
	Items []*k8sServiceBinding `json:"items"`
}

func NewK8sServiceBinding() *k8sServiceBinding {
	return &k8sServiceBinding{ApiVersion: serviceDomain + "/" + apiVersion,
		Kind: "ServiceBinding"}
}

type k8sServiceBinding struct {
	*model.ServiceBinding
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   Meta   `json:"metadata"`
}

func (kss *K8sServiceStorage) ListBindings() ([]string, error) {
	fmt.Println("listing all service bindings")
	r, e := http.Get(kss.defaultServiceBindingUri())
	if nil != e {
		return nil, fmt.Errorf("couldn't get the service bindings. %v, [%v]", e, r)
	}

	var ls listB
	e = json.NewDecoder(r.Body).Decode(&ls)
	if nil != e { // wrong json format error
		fmt.Println("json not unmarshalled:", e, r)
		return nil, e
	}
	fmt.Println("Got", len(ls.Items), "service bindings.")
	ret := make([]string, 0, len(ls.Items))
	for i, v := range ls.Items {
		fmt.Println("binding", i, v)
		ret = append(ret, v.ServiceBinding.ID)
	}
	return ret, nil
}

func (kss *K8sServiceStorage) GetBinding(id string) (*model.ServiceBinding, error) {
	fmt.Println("getting a single service binding")
	r, e := http.Get(kss.defaultServiceBindingUri() + "/" + id)
	if nil != e {
		return nil, fmt.Errorf("couldn't get the service binding. %v, [%v]", e, r)
	}

	var s k8sServiceBinding
	e = json.NewDecoder(r.Body).Decode(&s)
	if nil != e { // wrong json format error
		fmt.Println("json not unmarshalled:", e, r)
		return nil, e
	}
	fmt.Println("Got a service binding!", r.Body, s)

	return s.ServiceBinding, nil
}

func (kss *K8sServiceStorage) AddBinding(binding *model.ServiceBinding) error {
	fmt.Println(binding, binding.ID)

	ks := NewK8sServiceBinding()
	ks.Metadata = Meta{Name: binding.ID}
	ks.ServiceBinding = binding

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(&ks); nil != err {
		fmt.Println("failed to encode", binding, "as", ks)
		return err
	}
	defaultUri := kss.defaultServiceBindingUri()
	fmt.Printf("sending: %v\n to %v\n", b, defaultUri)
	r, e := http.Post(defaultUri, "application/json", b)
	fmt.Sprintf("result: %v", r)
	if nil != e || 201 != r.StatusCode {
		fmt.Printf("Error creating k8s service binding TPR [%s]...\n%v\n", e, r)
		return e
	}
	return nil
}

func (kss *K8sServiceStorage) SetBinding(binding *model.ServiceBinding) error {
	return fmt.Errorf("SetBinding: Not implemented yet")
}

func (kss *K8sServiceStorage) DeleteBinding(id string) error {
	return deleteResource(kss.defaultServiceBindingUri() + "/" + id)
}
