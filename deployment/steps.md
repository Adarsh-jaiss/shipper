```bash
kubectl config get-contexts

kubectl config use-context kind-kind

docker push adarshjaiss/shipper:latest

kubectl 

kubectl config set-context --current --namespace=shipper-backend


#create rbac.yaml
kubectl apply -f rbac.yaml

kubectl get serviceaccounts

#create config maps
kubectl apply -f configmap.yaml

kubectl get configmaps

#create secrets.yaml
kubectl apply -f secrets.yaml

kubectl get secrets

#create deployment
kubectl apply -f deployment.yaml

kubectl get deployments

# create service
kubectl apply -f service.yaml

kubectl get services

# Apply ingress for exposing your application to External IP
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml

kubectl apply ingress.yaml

kubectl get ingress

# Adding TLS to ingress

# Add a cert manager to automatically provision and manage TLS certificates
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml

kubectl get clusterissuer
kubectl describe clusterissuer letsencrypt-prod


# add a cert-issuer
kubectl apply -f encrypt-issuer.yaml
kubectl get certificate -n shipper-backend

kubectl create secret generic shipper-tls-cert \
  --namespace=shipper-backend \
  --from-file=tls.crt=tls.crt \
  --from-file=tls.key=tls.key


# or use Node Port

kubectl expose deployment shipper --type=NodePort --name=shipper-nodeport --port=80 --target-port=8080

kubectl get endpoints shipper-nodeport

kubectl port-forward svc/shipper-nodeport 8080:80

# to debug 

kubectl get challenges -n shipper-backend

kubectl exec -n shipper-backend shipper-6b6ffcf886-27fvc -- curl http://shipper:80/build


openssl req -new -newkey rsa:2048 -nodes -keyout shipper0.key -out shipper0.csr


```


https://cert-manager.io/docs/usage/ingress/

https://cert-manager.io/docs/installation/helm/