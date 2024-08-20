#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Custom error handling function
error_handler() {
    local exit_code=$?
    local line_number=$1
    local message=$2

    echo "Error at line $line_number: $message"
    echo "Exiting with code $exit_code..."
    exit $exit_code
}

# Set the error handler
trap 'error_handler $LINENO "$BASH_COMMAND"' ERR


# Azure login (assumes you're running this on a machine with Azure CLI installed)
echo "Logging into Azure..."
az login --identity

az aks install-cli

# Install Docker (if not already installed)
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
fi

# Get AKS credentials
echo "Getting AKS credentials..."
az aks get-credentials --resource-group ${RESOURCE_GROUP} --name ${CLUSTER_NAME} --overwrite-existing

# Verify connection to the cluster
echo "Verifying connection to the cluster..."
kubectl cluster-info
kubectl get nodes
kubectl get pods --all-namespaces
kubectl config view

# Apply all YAML files in the kubernetes directory
# use for manual path providng : for file in ${path.module}/kubernetes/shiper-build/*.yaml
echo "Installing dependencies..."
for file in /home/adarsh/myfiles/shipper/scripts/AKS/kubernetes/*.yaml
do
    echo "Applying $file"
    kubectl apply -f "$file"
done

echo "installing tekton pipeline"
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.50.5/release.yaml

echo "installing shipwright latest version"
kubectl apply --filename https://github.com/shipwright-io/build/releases/download/v0.13.0/release.yaml --server-side
curl --silent --location https://raw.githubusercontent.com/shipwright-io/build/v0.13.0/hack/setup-webhook-cert.sh | bash
curl --silent --location https://raw.githubusercontent.com/shipwright-io/build/main/hack/storage-version-migration.sh | bash

echo "generating a secret to access your container registry"
# Echo the Docker email, server, and username
echo "Docker Server: $REGISTRY_SERVER"
echo "Docker Username: $REGISTRY_USER"
echo "Docker Email: $REGISTRY_EMAIL"

kubectl create secret docker-registry push-secret \
    --docker-server=$REGISTRY_SERVER \
    --docker-username=$REGISTRY_USER \
    --docker-password=$REGISTRY_PASSWORD \
    --docker-email=$REGISTRY_EMAIL

echo "Secret created!!!"


echo "Creating a build..."
for file in /home/adarsh/myfiles/shipper/scripts/shiper-build/*.yaml
do
    echo "Applying $file"
    kubectl apply -f "$file"
    deployment_name=$(basename "$file" .yaml)
    echo "Waiting for deployment '$deployment_name' to be ready..."
    kubectl rollout status deployment/$deployment_name
    echo "build created successfully"
done


echo "Creating a buildrun..."
for file in /home/adarsh/myfiles/shipper/scripts/shiper-build/*.yaml
do
    echo "Applying $file"
    kubectl apply -f "$file"
    deployment_name=$(basename "$file" .yaml)
    echo "Waiting for deployment '$deployment_name' to be ready..."
    kubectl rollout status deployment/$deployment_name
done

