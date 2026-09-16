output "project_id" {
  description = "Google Cloud project managed by Terraform"
  value       = var.project_id
}

output "region" {
  description = "Google Cloud region used by the platform service"
  value       = var.region
}

output "artifact_registry_repository_name" {
  description = "Artifact Registry repository name"
  value       = google_artifact_registry_repository.platform_service.repository_id
}

output "artifact_registry_repository_url" {
  description = "Artifact Registry Docker repository URL"
  value = format(
    "%s-docker.pkg.dev/%s/%s",
    var.region,
    var.project_id,
    google_artifact_registry_repository.platform_service.repository_id,
  )
}
