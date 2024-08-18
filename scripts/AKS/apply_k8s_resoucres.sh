#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

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
echo "Applying Kubernetes resources..."
for file in ${path.module}/kubernetes/*.yaml
do
    echo "Applying $file"
    kubectl apply -f "$file"
done

echo "All Kubernetes resources applied successfully"