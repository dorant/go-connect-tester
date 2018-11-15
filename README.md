# go-connect-tester

## Prerequisite

Make sure needed dependencies are installed:
```
make deps
```
Installing protoc is a manual step.

### Helm
```
cd ~/go/src/github.com/dorant/go-connect-tester/

# Setup helm/tiller first (RBAC rules & service account needed on AWS)
kubectl create -f charts/helm/rbac-config.yaml
helm init --upgrade --service-account tiller
kubectl wait pods --for=condition=ready -n kube-system -l app=helm,name=tiller
```


## Build
```
make gen all push
```


## Deploy GRPC

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

## Deploy HTTP

Deploy the server and client to default namespace
```
# Deploy 3 servers
helm upgrade --install --force http-tester-server charts/http-tester-server/

# Deploy 1 client
helm upgrade --install --force http-tester-client charts/http-tester-client/

helm list

# To remove:
helm delete http-tester-server
helm delete http-tester-client
```

## Deploy federated HTTP

```
cd ~/go/src/github.com/dorant/go-connect-tester/

# Setup cluster names
export CLUSTER_1=xxxx
export CLUSTER_2=yyyy

# To replace current cluster names with other names, do:
find charts/http-tester-client-fedv2/ -name "*.yaml" | xargs sed -i "s/cluster1/${CLUSTER_1}/"
find charts/http-tester-server-fedv2/ -name "*.yaml" | xargs sed -i "s/cluster1/${CLUSTER_1}/"
find charts/http-tester-client-fedv2/ -name "*.yaml" | xargs sed -i "s/cluster2/${CLUSTER_2}/"
find charts/http-tester-server-fedv2/ -name "*.yaml" | xargs sed -i "s/cluster2/${CLUSTER_2}/"

# Make sure context is same as Fed.Controllers
kubectl config get-contexts

# Apply server
kubectl apply -f charts/http-tester-server-fedv2/
kubectl get pods --context=${CLUSTER_1}
kubectl get pods --context=${CLUSTER_2}

# Apply client
kubectl apply -f charts/http-tester-client-fedv2/
kubectl get pods --context=${CLUSTER_1}
kubectl get pods --context=${CLUSTER_2}

# Test of disabling server on cluster 1
find charts/http-tester-server-fedv2/federatedservice-placement.yaml | xargs sed -i "s/- ${CLUSTER_1}/# - ${CLUSTER_1}/"
cat charts/http-tester-server-fedv2/federatedservice-placement.yaml
find charts/http-tester-server-fedv2/federateddeployment-placement.yaml | xargs sed -i "s/- ${CLUSTER_1}/# - ${CLUSTER_1}/"
cat charts/http-tester-server-fedv2/federateddeployment-placement.yaml
kubectl apply -f charts/http-tester-server-fedv2/

# Test of enabling server again
find charts/http-tester-server-fedv2/federatedservice-placement.yaml | xargs sed -i "s/# - ${CLUSTER_1}/- ${CLUSTER_1}/"
cat charts/http-tester-server-fedv2/federatedservice-placement.yaml
find charts/http-tester-server-fedv2/federateddeployment-placement.yaml | xargs sed -i "s/# - ${CLUSTER_1}/- ${CLUSTER_1}/"
cat charts/http-tester-server-fedv2/federateddeployment-placement.yaml
kubectl apply -f charts/http-tester-server-fedv2/

# Remove the deployment
kubectl delete -f charts/http-tester-client-fedv2/
kubectl delete -f charts/http-tester-server-fedv2/

```
