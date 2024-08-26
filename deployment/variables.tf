# variable "docker_registry_server_url" {
#   description = "The URL of the private Docker registry server"
#   type        = string
# }

variable "docker_image_name" {
  description = "The name of the Docker image (without the tag)"
  type        = string
}

variable "docker_image_tag" {
  description = "The tag of the Docker image"
  type        = string
  default     = "latest"
}

variable "docker_registry_username" {
  description = "The username for the private Docker registry"
  type        = string
}

variable "docker_registry_password" {
  description = "The password for the private Docker registry"
  type        = string
  sensitive   = true
}