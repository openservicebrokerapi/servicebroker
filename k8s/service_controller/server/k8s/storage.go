package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/cncf/servicebroker/k8s/service_controller/model"
	"github.com/cncf/servicebroker/k8s/service_controller/server"
	"github.com/cncf/servicebroker/k8s/service_controller/utils"
)

type K8sServiceStorage struct {
}

const serviceDomain string = "cncf.org"
const apiVersion string = "v1alpha1"

// The k8s implementation should leverage Third Party Resources
// https://github.com/kubernetes/kubernetes/blob/master/docs/design/extending-api.md

var _ server.ServiceStorage = (*K8sServiceStorage)(nil)

func CreateServiceStorage() server.ServiceStorage {
	return &K8sServiceStorage{}
}

// listSB is only used for unmarshalling the list of service brokers
// for returning to the client
type listSB struct {
	Items []*model.ServiceBroker
}

func (kss *K8sServiceStorage) ListBrokers() ([]*model.ServiceBroker, error) {
	// get the ServiceBroker

	c := exec.Command("kubectl", "get", "ServiceBrokers", "-ojson")
	b, e := c.CombinedOutput()
	// b is json, an object, with an 'items' entry, which is an
	// array of service brokers.
	s := string(b)
	if nil != e {
		return nil, fmt.Errorf("couldn't get the service brokers. %v, [%v]", e, s)
	}

	var lsb listSB
	e = json.Unmarshal(b, &lsb)
	if nil != e { // wrong json format error
		return nil, e
	}
	return lsb.Items, nil
}

func (kss *K8sServiceStorage) GetBroker(name string) (*model.ServiceBroker, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) GetInventory(name string) (*model.Catalog, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

type Meta struct {
	Name string `json:"name"`
}

type k8sServiceBroker struct {
	*model.ServiceBroker
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   Meta   `json:"metadata"`
}

func NewK8sSB() *k8sServiceBroker {
	return &k8sServiceBroker{ApiVersion: serviceDomain + "/" + apiVersion,
		Kind: "ServiceBroker"}
}

func (kss *K8sServiceStorage) AddBroker(broker *model.ServiceBroker, catalog *model.Catalog) error {
	fmt.Println("adding broker to k8s")
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
	ksb.Metadata = Meta{Name: broker.Name}
	ksb.ServiceBroker = broker

	b, e := json.Marshal(ksb)
	if nil != e { // wrong json format error
		return e
	}
	s, e := utils.KubeCreateResource(bytes.NewReader(b))
	fmt.Sprintf("result: %v", s)
	if nil != e {
		fmt.Printf("Error creating k8s TPR [%s]...\n%v\n", e, s)
		return e
	}
	return nil
}

func (kss *K8sServiceStorage) DeleteBroker(string) error {
	return fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) ServiceExists(broker, service string) bool {
	return false
}

func (kss *K8sServiceStorage) ListServices(broker string) ([]*model.ServiceInstance, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) GetService(broker, service string) (*model.ServiceInstance, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) AddService(broker string, si *model.ServiceInstance) error {
	return fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) DeleteService(string, string) error {
	return fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) ListServiceBindings(string, string) ([]*model.ServiceBinding, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) GetServiceBinding(broker, service, binding string) (*model.Credential, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) AddServiceBinding(broker string, binding *model.ServiceBinding, cred *model.Credential) error {
	return fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) DeleteServiceBinding(string, string, string) error {
	return fmt.Errorf("Not implemented yet")
}
