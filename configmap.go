package configmap

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/schema"
	"github.com/rs/rest-layer/schema/query"
)

const configMapNameField = "id"
const namespaceField = "namespace"

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
		namespaceField: {
			Required:   false,
			Filterable: true,
			Sortable:   false,
			Validator: &schema.String{
				MinLen: 1,
			},
		},
		"data": {
			Required:   true,
			Filterable: false,
			Validator: &schema.Dict{
				KeysValidator: &schema.String{},
				Values: schema.Field{
					Validator: &schema.String{},
				},
			},
		},
	},
}

var (
  emptyLabels map[string]string
)

func (k *KubernetesClient) Insert(ctx context.Context, items []*resource.Item) (err error) {
	for _, item := range items {

		configMapName := fmt.Sprintf("%v", item.Payload[configMapNameField])
		actualKubernetesNamespace := k.getKubernetesNamespaceFromItem(item)
		configMap := buildConfigMap(item, actualKubernetesNamespace, emptyLabels)
		var existingConfigMap corev1.ConfigMap

		err = k.client.Get(ctx, actualKubernetesNamespace, configMapName, &existingConfigMap)

		if err != nil {
			apiErr, _ := err.(*k8s.APIError)
			if apiErr.Code == http.StatusNotFound {
				if err := k.client.Create(ctx, configMap); err != nil {
					return err
				}
			}
		}

		updatedConfigMap := buildConfigMap(item, actualKubernetesNamespace, existingConfigMap.Metadata.Labels)

		return k.client.Update(ctx, updatedConfigMap)
	}

	return nil
}

func (k *KubernetesClient) Delete(ctx context.Context, item *resource.Item) (err error) {
	return k.client.Delete(ctx, buildConfigMap(item, k.getKubernetesNamespaceFromItem(item), emptyLabels))
}

func (k *KubernetesClient) Update(ctx context.Context, item *resource.Item, original *resource.Item) (err error) {
	return k.client.Update(ctx, buildConfigMap(item, k.getKubernetesNamespaceFromItem(item), emptyLabels))
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
    configMap, kubernetesNamespace, err := k.findConfigMap(ctx, q)

    if err != nil {
        return nil, err
    } else if configMap != nil && kubernetesNamespace != "" {
        item, err := buildItem(configMap, kubernetesNamespace, configMap.Metadata.Name)
        if err != nil {
            return nil, err
        }
        list.Items = append(list.Items, item)
    }

	list.Total = len(list.Items)
	return list, nil
}
