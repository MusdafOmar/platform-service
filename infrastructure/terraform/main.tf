locals {
  required_services = toset([
    "artifactregistry.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com",
    "run.googleapis.com",
    "storage.googleapis.com",
  ])

  common_labels = {
    application = var.service_name
    environment = "development"
    managed_by  = "terraform"
  }
}

resource "google_project_service" "required" {
  for_each = local.required_services

  project = var.project_id
  service = each.value

  disable_on_destroy = false
}

resource "google_artifact_registry_repository" "platform_service" {
  project       = var.project_id
  location      = var.region
  repository_id = var.artifact_registry_repository
  description   = "Docker images for the LIA platform service"
  format        = "DOCKER"
  labels        = local.common_labels

  depends_on = [
    google_project_service.required,
  ]
}
