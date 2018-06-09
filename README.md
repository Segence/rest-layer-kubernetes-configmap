REST Layer Kubernetes ConfigMap Backend
=======================================

A backend to [REST Layer](http://rest-layer.io/) that stores configuration in a [Kubernetes](https://kubernetes.io) ConfigMap.

Further reading about ConfigMaps:
- [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/)
- [Kubernetes ConfigMaps and Secrets](https://medium.com/google-cloud/kubernetes-configmaps-and-secrets-68d061f7ab5b)

See [examples/main.go](https://github.com/Segence/rest-layer-kubernetes-configmap/blob/master/examples/main.go) for usage example.

## REST end-points:

| **Operation**             | **HTTP method** | **URL**                    | **Example payload**                                           |
|:--------------------------|:----------------|:---------------------------|:--------------------------------------------------------------|
| Create new ConfigMap      | POST            | `/api/config-map`          | `{{"id": "testconf", "data": {"config_value": "Hello"}}}`     |
| Update existing ConfigMap | PUT             | `/api/config-map/testconf` | `{{"id": "testconf", "data": {"config_value": "Hello2"}}}`    |
| Delete existing ConfigMap | DELETE          | `/api/config-map/testconf` | *None*                                                        |
| Find existing ConfigMap   | GET             | `/api/config-map/testconf` | *None*                                                        |

REST call examples are also available in [Postman](https://www.getpostman.com/) format [here](https://github.com/Segence/rest-layer-kubernetes-configmap/blob/master/examples/REST-Layer-Kubernetes-ConfigMap.postman_collection.json).

The functionality can be easily tested by running [Minikube](https://github.com/kubernetes/minikube) locally and using the out of cluster Kubernetes client configuration.
