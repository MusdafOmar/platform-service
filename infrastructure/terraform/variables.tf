variable "project_id" {
  description = "Google Cloud project ID used for the platform service"
  type        = string
}

variable "region" {
  description = "Google Cloud region used for regional resources"
  type        = string
  default     = "europe-north1"
}

variable "artifact_registry_repository" {
  description = "Name of the Artifact Registry Docker repository"
  type        = string
  default     = "platform-service"
}

variable "service_name" {
  description = "Name of the platform service"
  type        = string
  default     = "platform-service"
}
