package example

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/schema"
	"github.com/rs/rest-layer/schema/query"
)

const configMapNameField = "id"

var ConfigMapSchema = schema.Schema{
	Fields: schema.Fields{
		configMapNameField: {
			Required:   true,
			Filterable: false,
			Sortable:   false,
			Validator: &schema.String{
				MinLen: 1,
			},
		},
		"data": {
			Required:   true,
			Filterable: true,
			Validator: &schema.Dict{
				KeysValidator: &schema.String{},
				Values: schema.Field{
					Validator: &schema.String{},
				},
			},
		},
	},
}

func buildConfigMap(item *resource.Item, kubernetesNamespace string) *corev1.ConfigMap {
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
		},
		Data: allData,
	}
}

func buildItem(configMap *corev1.ConfigMap, configMapName string) (*resource.Item, error) {
	var configMapContent map[string]interface{}
	configMapContent = make(map[string]interface{})

	var configMapData map[string]string
	configMapData = make(map[string]string)

	for k, v := range configMap.Data {
		configMapData[k] = v
	}

	creationTimestamp := time.Unix(*configMap.Metadata.CreationTimestamp.Seconds, int64(*configMap.Metadata.CreationTimestamp.Nanos)).UTC().Format(time.RFC3339)

	configMapContent[configMapNameField] = configMapName
	configMapContent["data"] = configMapData
	configMapContent["creationTimestamp"] = creationTimestamp

	return resource.NewItem(configMapContent)
}

func NewHandler(kubernetesClient k8s.Client, kubernetesNamespace string) *KubernetesClient {
	return &KubernetesClient{
		client:    kubernetesClient,
		namespace: kubernetesNamespace,
	}
}

func (k *KubernetesClient) Insert(ctx context.Context, items []*resource.Item) (err error) {
	for _, item := range items {
		if err := k.client.Create(ctx, buildConfigMap(item, k.namespace)); err != nil {
			return err
		}
	}

	return nil
}

func (k *KubernetesClient) Delete(ctx context.Context, item *resource.Item) (err error) {
	return k.client.Delete(ctx, buildConfigMap(item, k.namespace))
}

func (k *KubernetesClient) Update(ctx context.Context, item *resource.Item, original *resource.Item) (err error) {
	return k.client.Update(ctx, buildConfigMap(item, k.namespace))
}

func (k *KubernetesClient) Clear(ctx context.Context, q *query.Query) (total int, err error) {
	configMaps, err := k.Find(ctx, q)
	if err != nil {
		return 0, err
	} else {
		for _, configMap := range configMaps.Items {
			if err := k.Delete(ctx, configMap); err != nil {
				return 0, err
			}
		}
	}
	return configMaps.Total, nil
}

func (k *KubernetesClient) Find(ctx context.Context, q *query.Query) (list *resource.ItemList, err error) {

	list = &resource.ItemList{Items: []*resource.Item{}}

	for _, exp := range q.Predicate {
		switch t := exp.(type) {
		case query.Equal:
			if t.Field == configMapNameField {

				configMapName := t.Value.(string)

				var configMap corev1.ConfigMap
				err := k.client.Get(ctx, k.namespace, configMapName, &configMap)
				if err != nil {
					return nil, err
				}
				item, err := buildItem(&configMap, configMapName)
				if err != nil {
					return nil, err
				}
				list.Items = append(list.Items, item)

			} else {
				return nil, errors.New("Querying can only be done if the '" + configMapNameField + "' field is set")
			}
		default:
			return nil, resource.ErrNotImplemented
		}
	}

	list.Total = len(list.Items)
	return list, nil
}
