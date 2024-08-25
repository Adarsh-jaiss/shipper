
resource "azurerm_resource_group" "aks_rg" {
  name     = var.resource_group_name
  location = var.location
}

resource "azurerm_kubernetes_cluster" "aks" {
  name                = var.cluster_name
  location            = azurerm_resource_group.aks_rg.location
  resource_group_name = azurerm_resource_group.aks_rg.name
  dns_prefix          = var.dns_prefix

  default_node_pool {
    name                = "default"
    node_count          = var.default_node_count
    vm_size             = var.default_node_vm_size
    enable_auto_scaling = false
    # min_count           = var.default_node_min_count
    # max_count           = var.default_node_max_count
    # max_pods            = var.default_node_max_pods
    os_disk_size_gb     = var.default_node_os_disk_size_gb
    type                = "VirtualMachineScaleSets"
    zones               = [1, 2, 3]
  }

  identity {
    type = "SystemAssigned"
  }


  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "standard"
  }


  oms_agent {
    log_analytics_workspace_id = azurerm_log_analytics_workspace.aks.id
  }



  auto_scaler_profile {
    balance_similar_node_groups      = true
    expander                         = "random"
    max_graceful_termination_sec     = 600
    scale_down_delay_after_add       = "10m"
    scale_down_delay_after_delete    = "10s"
    scale_down_delay_after_failure   = "3m"
    scan_interval                    = "10s"
    scale_down_unneeded              = "10m"
    scale_down_unready               = "20m"
    scale_down_utilization_threshold = "0.5"
  }

  tags = var.tags
}

resource "azurerm_kubernetes_cluster_node_pool" "additional" {
  name                  = "additional"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.aks.id
  vm_size               = var.additional_node_vm_size
  node_count            = var.additional_node_count
  enable_auto_scaling   = true
  min_count             = var.additional_node_min_count
  max_count             = var.additional_node_max_count
  max_pods              = var.additional_node_max_pods
  os_disk_size_gb       = var.additional_node_os_disk_size_gb
  zones                 = [1, 2, 3]

  tags = var.tags
}
resource "azurerm_log_analytics_workspace" "aks" {
  name                = "${var.cluster_name}-logs"
  location            = azurerm_resource_group.aks_rg.location
  resource_group_name = azurerm_resource_group.aks_rg.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "null_resource" "apply_kubernetes_resources" {
  depends_on = [
    azurerm_kubernetes_cluster.aks,
    azurerm_kubernetes_cluster_node_pool.additional
  ]

  provisioner "local-exec" {
    command = <<EOT
      set -e
      # Retry function
      retry() {
        local retries=3
        local delay=10
        local cmd="$*"
        local count=0
        while ! $cmd; do
          if [ $count -lt $retries ]; then
            echo "Command failed. Retry $((count + 1))/$retries in $delay seconds..."
            sleep $delay
            count=$((count + 1))
          else
            echo "Command failed after $retries retries."
            return 1
          fi
        done
        return 0
      }

      retry az aks get-credentials --resource-group ${azurerm_resource_group.aks_rg.name} --name ${azurerm_kubernetes_cluster.aks.name} --overwrite-existing
      
      echo "Verifying connection to the cluster..."
      retry kubectl cluster-info
      retry kubectl get nodes
    
      echo "Installing dependencies via Daemon set..."
      for file in /home/adarsh/myfiles/shipper/infra/AKS/kubernetes/*.yaml
      do
          echo "Applying $file"
          retry kubectl apply -f "$file"
      done

      if !jq --version &> /dev/null; then
      echo "Installing JQ..."
      sudo wget -O /usr/local/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64
      sudo chmod +x /usr/local/bin/jq
      jq --version
      fi

      echo "Installing tekton"
      retry kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.50.5/release.yaml
      echo "tekton installed successfully!!!"

      echo "Ensuring Shipwright is installed..."
      if ! kubectl get deployment -n shipwright-build shipwright-build-controller &> /dev/null; then
          echo "Installing Shipwright..."
          retry kubectl apply --filename https://github.com/shipwright-io/build/releases/download/v0.13.0/release.yaml --server-side
          retry curl --silent --location https://raw.githubusercontent.com/shipwright-io/build/v0.13.0/hack/setup-webhook-cert.sh | bash
          retry kubectl apply --filename https://github.com/shipwright-io/build/releases/download/v0.13.0/sample-strategies.yaml --server-side
          # retry curl --silent --location https://raw.githubusercontent.com/shipwright-io/build/main/hack/storage-version-migration.sh | bash
      fi

      echo "all dependencies installed successfully!!!"
      exit

    EOT

    environment = {
      RESOURCE_GROUP      = azurerm_resource_group.aks_rg.name
      CLUSTER_NAME        = azurerm_kubernetes_cluster.aks.name
      AZURE_CLIENT_ID     = var.client_id
      AZURE_CLIENT_SECRET = var.client_secret
      AZURE_TENANT_ID     = var.tenant_id
    }
  }
}
