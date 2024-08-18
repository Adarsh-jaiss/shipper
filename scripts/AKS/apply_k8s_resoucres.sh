#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Azure login (assumes you're running this on a machine with Azure CLI installed)
echo "Logging into Azure..."
az login --identity

# Get AKS credentials
echo "Getting AKS credentials..."
az aks get-credentials --resource-group ${RESOURCE_GROUP} --name ${CLUSTER_NAME} --overwrite-existing

# Verify connection to the cluster
echo "Verifying connection to the cluster..."
kubectl get nodes

# Apply all YAML files in the kubernetes directory
echo "Applying Kubernetes resources..."
for file in ${path.module}/kubernetes/*.yaml
do
    echo "Applying $file"
    kubectl apply -f "$file"
done

echo "All Kubernetes resources applied successfully"