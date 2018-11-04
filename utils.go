package configmap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/schema/query"
)

func buildConfigMap(item *resource.Item, kubernetesNamespace string, labels map[string]string) *corev1.ConfigMap {
	configMapName := fmt.Sprintf("%v", item.Payload[configMapNameField])

	var allData map[string]string
	allData = make(map[string]string)

	if item.Payload["data"] != nil {
		for k, v := range item.Payload["data"].(map[string]interface{}) {
			if k != configMapNameField {
				allData[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	return &corev1.ConfigMap{
		Metadata: &metav1.ObjectMeta{
			Name:      k8s.String(configMapName),
			Namespace: k8s.String(kubernetesNamespace),
			Labels:    labels,
		},
		Data: allData,
	}
}

func buildItem(configMap *corev1.ConfigMap, kubernetesNamespace string, configMapName *string) (*resource.Item, error) {
	var configMapContent map[string]interface{}
	configMapContent = make(map[string]interface{})

	var configMapData map[string]interface{}
	configMapData = make(map[string]interface{})

	for k, v := range configMap.Data {
		configMapData[k] = v
	}

	creationTimestamp := time.Unix(*configMap.Metadata.CreationTimestamp.Seconds, int64(*configMap.Metadata.CreationTimestamp.Nanos)).UTC().Format(time.RFC3339)

	configMapContent[configMapNameField] = configMapName
	configMapContent[namespaceField] = kubernetesNamespace
	configMapContent["data"] = configMapData
	configMapContent["creationTimestamp"] = creationTimestamp

	return resource.NewItem(configMapContent)
}

func (k *KubernetesClient) getKubernetesNamespace(namespace string) string {
	kubernetesNamespace := k.defaultNamespace
	if namespace != "" {
		kubernetesNamespace = namespace
	}
	return kubernetesNamespace
}

func (k *KubernetesClient) getKubernetesNamespaceFromItem(item *resource.Item) string {
	kubernetesNamespace := ""
	if item.Payload["namespace"] != nil {
		kubernetesNamespace = item.Payload[namespaceField].(string)
	}
	return k.getKubernetesNamespace(kubernetesNamespace)
}

func NewHandler(kubernetesClient k8s.Client, kubernetesNamespace string) *KubernetesClient {
	return &KubernetesClient{
		client:           kubernetesClient,
		defaultNamespace: kubernetesNamespace,
	}
}

func (k *KubernetesClient) findConfigMap(ctx context.Context, q *query.Query) (*corev1.ConfigMap, string, error) {

	var configMapName string
	var kubernetesNamespace string

	for _, exp := range q.Predicate {
		switch t := exp.(type) {
		case query.Equal:
			if t.Field == configMapNameField {
				configMapName = t.Value.(string)
			} else if t.Field == namespaceField {
				kubernetesNamespace = t.Value.(string)
			} else {
				return nil, "", errors.New("Querying can only be done if the '" + configMapNameField + "' or '" + namespaceField + "' fields are set")
			}
		default:
			return nil, "", resource.ErrNotImplemented
		}
	}

	if configMapName == "" {
		return nil, "", errors.New("ConfigMap name is not provided")
	}

	actualKubernetesNamespace := k.getKubernetesNamespace(kubernetesNamespace)

	var configMap corev1.ConfigMap
	err := k.client.Get(ctx, actualKubernetesNamespace, configMapName, &configMap)
	if err != nil {
		apiErr, ok := err.(*k8s.APIError)
		if !ok {
			return nil, "", err
		}
		if apiErr.Code != http.StatusNotFound {
			return nil, "", err
		}
	} else {
		return &configMap, actualKubernetesNamespace, nil
	}

	return nil, "", nil
}
