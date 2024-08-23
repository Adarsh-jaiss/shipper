#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

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

# Check if required variables are set
if [ -z "$REGISTRY_SERVER" ] || [ -z "$REGISTRY_USER" ] || [ -z "$REGISTRY_PASSWORD" ] || [ -z "$REGISTRY_EMAIL" ] || [ -z "$AZURE_CLIENT_ID" ] || [ -z "$AZURE_CLIENT_SECRET" ] || [ -z "$AZURE_TENANT_ID" ]; then
    echo "Error: Missing required environment variables"
    exit 1
fi

# Function to check if the cluster is accessible
check_cluster_accessibility() {
    kubectl cluster-info > /dev/null 2>&1
    return $?
}

# Wait for the cluster to become accessible
max_retries=30
retry_interval=10

for ((i=0; i<max_retries; i++)); do
    if check_cluster_accessibility; then
        echo "Cluster is accessible. Proceeding with resource application."
        break
    else
        echo "Cluster is not yet accessible. Retrying in $retry_interval seconds..."
        sleep $retry_interval
    fi
done

if ! check_cluster_accessibility; then
    echo "Failed to access the cluster after $max_retries attempts. Exiting."
    exit 1
fi


#deemon set
echo "Installing dependencies via Daemon set..."
for file in /home/adarsh/myfiles/shipper/scripts/AKS/kubernetes/*.yaml
do
    echo "Applying $file"
    kubectl apply -f "$file"
done


# Azure login
echo "Logging into Azure..."
az login --service-principal --username $AZURE_CLIENT_ID --password $AZURE_CLIENT_SECRET --tenant $AZURE_TENANT_ID

echo "User logged in successfully"

# Get AKS credentials
echo "Getting AKS credentials..."
az aks get-credentials --resource-group ${RESOURCE_GROUP} --name ${CLUSTER_NAME} --overwrite-existing

# Verify connection to the cluster
echo "Verifying connection to the cluster..."
kubectl cluster-info
kubectl get nodes

# Install 

# Ensure Shipwright is installed
echo "Ensuring Shipwright is installed..."
if ! kubectl get deployment -n shipwright-build shipwright-build-controller &> /dev/null; then
    echo "Installing Shipwright..."
    kubectl apply --filename https://github.com/shipwright-io/build/releases/download/v0.13.0/release.yaml --server-side
    curl --silent --location https://raw.githubusercontent.com/shipwright-io/build/v0.13.0/hack/setup-webhook-cert.sh | bash
    curl --silent --location https://raw.githubusercontent.com/shipwright-io/build/main/hack/storage-version-migration.sh | bash
fi

# echo "Generating a secret to access your container registry"
# kubectl create secret docker-registry push-secret \
#     --docker-server=$REGISTRY_SERVER \
#     --docker-username=$REGISTRY_USER \
#     --docker-password=$REGISTRY_PASSWORD \
#     --docker-email=$REGISTRY_EMAIL \
#     --dry-run=client -o yaml | kubectl apply -f -

# echo "Secret created or updated!"

# echo "Applying Shipwright Build and BuildRun resources..."
# for file in /path/to/your/shipwright/resources/*.yaml
# do
#     echo "Applying $file"
#     kubectl apply -f "$file"
#     resource_name=$(basename "$file" .yaml)
#     resource_type=$(yq e '.kind' "$file" | tr '[:upper:]' '[:lower:]')
#     echo "Waiting for $resource_type '$resource_name' to be ready..."
#     kubectl wait --for=condition=Ready $resource_type/$resource_name --timeout=300s
# done

# echo "All resources applied and ready!"