REST Layer Kubernetes ConfigMap Backend
=======================================

See [examples/main.go](https://github.com/Segence/rest-layer-kubernetes-configmap/blob/master/examples/main.go) for usage example.

## REST end-points:

| **Operation**               | **HTTP method** | **URL**                    | **Payload**                                                   |
|:----------------------------|:----------------|:---------------------------|:--------------------------------------------------------------|
| *Create new ConfigMap*      | POST            | `/api/config-map`          | `{{"id": "testconf", "data": {"config_value": "Hello"}}}`     |
| *Update existing ConfigMap* | PUT             | `/api/config-map/testconf` | `{{"id": "testconf", "data": {"config_value": "Hello2"}}}`    |
| *Delete existing ConfigMap* | DELTE           | `/api/config-map/testconf` | *None*                                                        |
| *Find existing ConfigMap*   | GET             | `/api/config-map/testconf` | *None*                                                        |

REST call examples are also available in [Postman](https://www.getpostman.com/) format [here](https://github.com/Segence/rest-layer-kubernetes-configmap/blob/master/examples/REST-Layer-Kubernetes-ConfigMap.postman_collection.json).
