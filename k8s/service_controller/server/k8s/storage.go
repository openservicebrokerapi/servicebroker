package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cncf/servicebroker/k8s/service_controller/model"
	"github.com/cncf/servicebroker/k8s/service_controller/server"
	"github.com/cncf/servicebroker/k8s/service_controller/utils"
	//"k8s.io/kubernetes/pkg/client/restclient"
	k8sclient "k8s.io/kubernetes/pkg/client/unversioned"
	defaultClient "k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	//	configFactory "k8s.io/kubernetes/pkg/kubectl/cmd/util"
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
	Items []*k8sServiceBroker `json:"items"`
}

func (kss *K8sServiceStorage) ListBrokers() ([]*model.ServiceBroker, error) {
	fmt.Println("listing all brokers")
	// get the ServiceBroker

	c := exec.Command("kubectl", "get", "ServiceBrokers", "-ojson")
	b, e := c.CombinedOutput()
	// b is json, an object, with an 'items' entry, which is an
	// array of service brokers.
	s := string(b)
	if nil != e {
		return nil, fmt.Errorf("couldn't get the service brokers. %v, [%v]", e, s)
	}

	fmt.Println("returned json: ", s)

	var lsb listSB
	e = json.Unmarshal(b, &lsb)
	if nil != e { // wrong json format error
		fmt.Println("json not unmarshalled:", e, s)
		return nil, e
	}
	fmt.Println("Got", len(lsb.Items), "brokers.")
	ret := make([]*model.ServiceBroker, 0, len(lsb.Items))
	for _, v := range lsb.Items {
		ret = append(ret, v.ServiceBroker)
	}
	return ret, nil
}

const (
	defaultHost string = "http://127.0.0.1:8080"
)

func printVersion(name string) {
	// poc that this works by getting the server version
	cfg, _ := defaultClient.DefaultClientConfig.ClientConfig()
	c, err := k8sclient.New(cfg)
	info, err := c.Discovery().ServerVersion()
	if nil != err {
		fmt.Printf("Error %v\n", err)
	}
	fmt.Printf("server API version information: %s\n", info)

	// nope
	tpr, err := c.Extensions().ThirdPartyResources().Get("service-broker.cncf.org")
	if nil != err {
		fmt.Printf("Error %v\n", err)
	}
	fmt.Printf("maybe tpr?: %v\n", tpr)

	// nope
	ep, err := c.Endpoints("default").Get("ServiceBrokers")
	if nil != err {
		fmt.Printf("Error %v\n", err)
	}
	fmt.Printf("maybe all brokers?: %v\n", ep)
}

func (kss *K8sServiceStorage) GetBroker(name string) (*model.ServiceBroker, error) {
	printVersion(name)

	c := exec.Command("kubectl", "get", "-ojson", "ServiceBrokers", name)
	b, e := c.CombinedOutput()
	s := string(b)
	if nil != e {
		return nil, fmt.Errorf("couldn't get the service broker. %v, [%v]", e, s)
	}
	fmt.Println("returned json: ", s)
	var sb k8sServiceBroker
	e = json.Unmarshal(b, &sb)
	if nil != e { // wrong json format error
		return nil, e
	}
	return sb.ServiceBroker, nil
}

func (kss *K8sServiceStorage) GetBrokerByService(id string) (*model.ServiceBroker, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) GetInventory() (*model.Catalog, error) {
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
	ksb.Metadata = Meta{Name: broker.Name}
	ksb.ServiceBroker = broker

	b, e := json.Marshal(ksb)
	fmt.Println(string(b))
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

func (kss *K8sServiceStorage) DeleteBroker(id string) error {
	c := exec.Command("kubectl", "delete", "-oname", "ServiceBrokers", id)
	b, e := c.CombinedOutput()
	if nil != e { // some kind of exec error
		return e
	}
	s := string(b)
	fmt.Println("deleted: ", s)
	lookingFor := "servicebroker/" + id
	if strings.Contains(s, lookingFor) {
		return nil
	}
	return fmt.Errorf("didn't work right: %v", s)
}

func (kss *K8sServiceStorage) ServiceExists(id string) bool {
	return false
}

func (kss *K8sServiceStorage) ListServices() ([]*model.ServiceInstance, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) GetService(id string) (*model.ServiceInstance, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) AddService(si *model.ServiceInstance) error {
	return fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) DeleteService(id string) error {
	return fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) ListServiceBindings() ([]*model.ServiceBinding, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) GetServiceBinding(id string) (*model.Credential, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) AddServiceBinding(binding *model.ServiceBinding, cred *model.Credential) error {
	return fmt.Errorf("Not implemented yet")
}

func (kss *K8sServiceStorage) DeleteServiceBinding(id string) error {
	return fmt.Errorf("Not implemented yet")
}
