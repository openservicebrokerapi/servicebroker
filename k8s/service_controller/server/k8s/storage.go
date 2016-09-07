package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/servicebroker/servicebroker/k8s/service_controller/server"
	model "github.com/servicebroker/servicebroker/model/service_controller"
)

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

var _ server.ServiceStorage = (*K8sServiceStorage)(nil)

type Meta struct {
	Name string `json:"name"`
}

type k8sServiceBroker struct {
	*model.ServiceBroker
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

var versionMap []VName = []VName{{"v1alpha1"}}
var serviceBrokerMeta Meta = Meta{"service-broker.cncf.org"}
var serviceMeta Meta = Meta{"sbservice.cncf.org"} // sbservice so it does not conflict with the built in Service
var serviceBindingMeta Meta = Meta{"service-binding.cncf.org"}
var serviceInstanceMeta Meta = Meta{"service-instance.cncf.org"}
var serviceBrokerDefinition TPR = TPR{serviceBrokerMeta, TPRapiVersion, thirdPartyResourceString, versionMap}
var serviceDefinition TPR = TPR{serviceMeta, TPRapiVersion, thirdPartyResourceString, versionMap}
var serviceBindingDefinition TPR = TPR{serviceBindingMeta, TPRapiVersion, thirdPartyResourceString, versionMap}
var serviceInstanceDefinition TPR = TPR{serviceInstanceMeta, TPRapiVersion, thirdPartyResourceString, versionMap}

func CreateServiceStorage(host string) server.ServiceStorage {
	k := &K8sServiceStorage{host: host,
		defaultResource: fmt.Sprintf(defaultUri, host)}
	fmt.Println(" root host is:", k.defaultUri())
	// define the resources once at startup
	// results in ServiceBrokers

	k.createTPR(serviceBrokerDefinition)
	k.createTPR(serviceDefinition)
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

func (kss *K8sServiceStorage) ListBrokers() ([]*model.ServiceBroker, error) {
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
	ret := make([]*model.ServiceBroker, 0, len(lsb.Items))
	for _, v := range lsb.Items {
		ret = append(ret, v.ServiceBroker)
	}
	return ret, nil
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

func (kss *K8sServiceStorage) GetBrokerByService(id string) (*model.ServiceBroker, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) GetInventory() (*model.Catalog, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func NewK8sSB() *k8sServiceBroker {
	return &k8sServiceBroker{ApiVersion: serviceDomain + "/" + apiVersion,
		Kind: "ServiceBroker"}
}

func (kss *K8sServiceStorage) AddBroker(broker *model.ServiceBroker, catalog *model.Catalog) error {
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
	ksb.Metadata = Meta{Name: broker.GUID}
	ksb.ServiceBroker = broker

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&ksb)
	fmt.Printf("sending: %v", b)
	r, e := http.Post(kss.defaultUri(), "application/json", b)
	fmt.Sprintf("result: %v", r)
	if nil != e || 201 != r.StatusCode {
		fmt.Printf("Error creating k8s TPR [%s]...\n%v\n", e, r)
		return e
	}

	fmt.Println("installing the", len(catalog.Services), "services for this broker")
	for i, service := range catalog.Services {
		fmt.Println(i, service)
	}

	return nil
}

func (kss *K8sServiceStorage) DeleteBroker(name string) error {
	uri := kss.defaultUri() + "/" + name
	fmt.Println("uri is:", uri)

	// utter failure of an http API
	req, _ := http.NewRequest("DELETE", uri, nil)
	_, e := http.DefaultClient.Do(req)
	if nil != e {
		return fmt.Errorf("couldn't nuke %v, [%v]", name, e)
	}
	return nil
}

func (kss *K8sServiceStorage) ServiceExists(id string) bool {
	return false
}

func (kss *K8sServiceStorage) ListServices() ([]*model.ServiceInstanceData, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) GetService(id string) (*model.ServiceInstanceData, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) AddService(si *model.ServiceInstanceData) error {
	return fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) SetService(si *model.ServiceInstanceData) error {
	return fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) DeleteService(id string) error {
	return fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) ListServiceBindings() ([]*model.ServiceBinding, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) GetServiceBinding(id string) (*model.ServiceBinding, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) AddServiceBinding(binding *model.ServiceBinding, cred *model.Credential) error {
	return fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) DeleteServiceBinding(id string) error {
	return fmt.Errorf("Not implemented yet")
}
