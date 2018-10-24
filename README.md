# go-grpc-tester

## Prerequisite

Make sure needed dependencies are installed:
```
make deps
```
Installing protoc is a manual step.

### Helm
```
# Setup helm/tiller first (RBAC rules & service account needed on AWS)
kubectl create -f charts/helm/rbac-config.yaml
helm init --upgrade --service-account tiller
```


## Build
```
make gen client server push
```


## Deploy

Deploy the server and client to default namespace
```
# Deploy 2+ servers
helm upgrade --install --force grpc-tester-server charts/grpc-tester-server/

# Deploy 1 client
helm upgrade --install --force grpc-tester-client charts/grpc-tester-client/

helm list

# To remove:
helm delete grpc-tester-server
helm delete grpc-tester-client
```
