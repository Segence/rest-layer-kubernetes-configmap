package configmap

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/ericchiang/k8s"
	"github.com/ghodss/yaml"
)

type KubernetesClient struct {
	sync.RWMutex

	client    k8s.Client
	namespace string
}

func newInClusterClient() (*k8s.Client, error) {
	return k8s.NewInClusterClient()
}

func newOutOfClusterClient(kubeconfigPath string) (*k8s.Client, error) {
	data, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("read kubeconfig: %v", err)
	}

	// Unmarshal YAML into a Kubernetes config object.
	var config k8s.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal kubeconfig: %v", err)
	}
	return k8s.NewClient(&config)
}

func NewKubernetesClient(inCluster bool, kubeconfigPath string) (*k8s.Client, error) {
	if inCluster {
		return newInClusterClient()
	} else {
		return newOutOfClusterClient(kubeconfigPath)
	}
}
