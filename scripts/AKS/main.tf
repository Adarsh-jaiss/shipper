
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
    enable_auto_scaling = true
    min_count           = var.default_node_min_count
    max_count           = var.default_node_max_count
    max_pods            = var.default_node_max_pods
    os_disk_size_gb     = var.default_node_os_disk_size_gb
    type                = "VirtualMachineScaleSets"
    zones  = [1, 2, 3]
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
  zones    = [1, 2, 3]

  tags = var.tags
}
resource "azurerm_log_analytics_workspace" "aks" {
  name                = "${var.cluster_name}-logs"
  location            = azurerm_resource_group.aks_rg.location
  resource_group_name = azurerm_resource_group.aks_rg.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}