REST Layer Kubernetes ConfigMap Backend
=======================================

A backend to [REST Layer](http://rest-layer.io/) that stores configuration in a [Kubernetes](https://kubernetes.io) ConfigMap.

Further reading about ConfigMaps:
- [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/)
- [Kubernetes ConfigMaps and Secrets](https://medium.com/google-cloud/kubernetes-configmaps-and-secrets-68d061f7ab5b)

See [examples/main.go](https://github.com/Segence/rest-layer-kubernetes-configmap/blob/master/examples/main.go) for usage example.

## Building

If you're using the [GB](https://getgb.io) build tool, you can fetch all dependencies using the following commands:

```
gb vendor fetch github.com/justinas/alice
gb vendor fetch github.com/rs/cors
gb vendor fetch github.com/rs/rest-layer/resource
gb vendor fetch github.com/rs/rest-layer/rest
gb vendor fetch github.com/rs/zerolog
gb vendor fetch github.com/rs/zerolog/hlog
gb vendor fetch github.com/rs/zerolog/log
gb vendor fetch github.com/segence/rest-layer-kubernetes-configmap
```

## Kubernetes namespaces

The library can access any Kubernetes namespace as long as it has permissions (which can be controlled through [RBAC](https://kubernetes.io/docs/reference/access-authn-authz/rbac/), for example).

The client has to be configred with a default Kubernetes namespace, like in the following code example:

    kHandler := configmap.NewHandler(*kubernetesClient, "some-namespace-name")

However, the library can access ConfigMaps in other namespaces.

### Creating or updating a ConfigMap in another namespace

Add the `namespace` field into the payload, e.g.:

```
{
	"id": "test",
	"namespace": "my-namespace",
	"data": {
		...
	}
}
```

### Querying or deleting a ConfigMap in another namespace

Add the `?filter=` query parameter to the end of the URL.

E.g.: `/api/config-map/test?filter={namespace:"my-namespaceap"}`

## REST end-points:

| **Operation**             | **HTTP method** | **URL**                    | **Example payload**                                           |
|:--------------------------|:----------------|:---------------------------|:--------------------------------------------------------------|
| Create new ConfigMap      | POST            | `/api/config-map`          | `{{"id": "testconf", "data": {"config_value": "Hello"}}}`     |
| Update existing ConfigMap | PUT             | `/api/config-map/testconf` | `{{"id": "testconf", "data": {"config_value": "Hello2"}}}`    |
| Delete existing ConfigMap | DELETE          | `/api/config-map/testconf` | *None*                                                        |
| Find existing ConfigMap   | GET             | `/api/config-map/testconf` | *None*                                                        |

REST call examples are also available in [Postman](https://www.getpostman.com/) format [here](https://github.com/Segence/rest-layer-kubernetes-configmap/blob/master/examples/REST-Layer-Kubernetes-ConfigMap.postman_collection.json).

## Testing

The functionality can be easily tested by running [Minikube](https://github.com/kubernetes/minikube) locally and using the out of cluster Kubernetes client configuration.
