# teritorid Kubernetes config

## Deploy

### Set env variables

```sh
domain=<your-domain>
registry=<your-docker-registry>
teritori_data_path=<path-to-teritori-data>
namespace=<your-namespace>
node=<your-node>
```

### Create namespace and set as current

```sh
kubectl create namespace $namespace
kubectl config set-context --current --namespace=$namespace
```

### Build and push the image

```sh
DOCKER_BUILDKIT=1 registry=$registry ./docker-publish.sh
```

### Deploy

```sh
sed "s#YOUR_REGISTRY#${registry}#g" ./deploy-teritorid.yaml \
| sed "s#IMAGE_COMMIT#$(git rev-parse --short HEAD)#g" \
| sed "s#TERITORI_DATA_PATH#${teritori_data_path}#g" \
| sed "s#YOUR_NAMESPACE#${namespace}#g" \
| sed "s#YOUR_NODE#${node}#g" \
| kubectl apply -f -
```

### Create ingress

```sh
sed "s#YOUR_DOMAIN#${domain}#g" ./ingress.yaml | kubectl apply -f -
```

The services will be available at `teritorid.<your-domain>/rest` and `teritorid.<your-domain>/rpc`


### Important

TODO: automate this

To get the node fully working, you need to:

- copy the correct `genesis.json` file to `<path-to-teritori-data>/config/genesis.json`
- set the correct `persistent_peers` value in `<path-to-teritori-data>/config/config.toml`
- enable the rest endpoint in the `[api]` section of the `<path-to-teritori-data>/config/app.toml` file
- set the rpc `laddr` to `tcp/0.0.0.0:26657` in the `[rpc]` section of the `<path-to-teritori-data>/config/config.toml` file
- restart the pod