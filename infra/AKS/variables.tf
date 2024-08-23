variable "resource_group_name" {
  description = "Name of the resource group"
  type        = string
}

variable "reuse_existing_resources" {
  description = "Whether to reuse existing resources if they are present"
  type        = bool
  default     = false
}

variable "location" {
  description = "Azure region to deploy resources"
  type        = string
}

variable "cluster_name" {
  description = "Name of the AKS cluster"
  type        = string
}

variable "dns_prefix" {
  description = "DNS prefix for the AKS cluster"
  type        = string
}

variable "default_node_count" {
  description = "Number of nodes in the default node pool"
  type        = number
  default     = 2
}

variable "default_node_vm_size" {
  description = "VM size for the default node pool"
  type        = string
  default     = "Standard_DS2_v2"
}

variable "default_node_min_count" {
  description = "Minimum number of nodes for autoscaling in the default node pool"
  type        = number
  default     = 1
}

variable "default_node_max_count" {
  description = "Maximum number of nodes for autoscaling in the default node pool"
  type        = number
  default     = 3
}

variable "additional_node_count" {
  description = "Number of nodes in the additional node pool"
  type        = number
  default     = 2
}

variable "additional_node_vm_size" {
  description = "VM size for the additional node pool"
  type        = string
  default     = "Standard_DS3_v2"
}

variable "additional_node_min_count" {
  description = "Minimum number of nodes for autoscaling in the additional node pool"
  type        = number
  default     = 1
}

variable "additional_node_max_count" {
  description = "Maximum number of nodes for autoscaling in the additional node pool"
  type        = number
  default     = 5
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}

variable "default_node_max_pods" {
  description = "Maximum number of pods per node in the default node pool"
  type        = number
  default     = 30
}

variable "default_node_os_disk_size_gb" {
  description = "OS disk size for nodes in the default node pool"
  type        = number
  default     = 50
}

variable "additional_node_max_pods" {
  description = "Maximum number of pods per node in the additional node pool"
  type        = number
  default     = 30
}

variable "additional_node_os_disk_size_gb" {
  description = "OS disk size for nodes in the additional node pool"
  type        = number
  default     = 50
}

# Authentication

variable "subscription_id" {
  description = "Azure subscription ID"
  type        = string
}

variable "client_id" {
  description = "Azure service principal client ID"
  type        = string
}

variable "client_secret" {
  description = "Azure service principal client secret"
  type        = string
}

variable "tenant_id" {
  description = "Azure tenant ID"
  type        = string
}

# registry 

# variable "REGISTRY_SERVER" {
#   type        = string
#   description = "Registry server address"
# }

# variable "REGISTRY_USER" {
#   type        = string
#   description = "Registry user"
# }

# variable "REGISTRY_PASSWORD" {
#   type        = string
#   description = "Registry password"
# }

# variable "REGISTRY_EMAIL" {
#   type        = string
#   description = "Registry email"
# }