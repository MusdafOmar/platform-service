resource "google_service_account" "platform_service" {
  project      = var.project_id
  account_id   = var.service_name
  display_name = "Platform Service Cloud Run identity"
  description  = "Runtime identity for the LIA platform service"
}

resource "google_cloud_run_v2_service" "platform_service" {
  project  = var.project_id
  name     = var.service_name
  location = var.region

  deletion_protection = false
  ingress             = "INGRESS_TRAFFIC_ALL"
  labels              = local.common_labels

  template {
    service_account = google_service_account.platform_service.email

    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }

    containers {
      image = format(
        "%s/%s:%s",
        google_artifact_registry_repository.platform_service.registry_uri,
        var.service_name,
        var.container_image_tag,
      )

      ports {
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }

        cpu_idle          = true
        startup_cpu_boost = false
      }

      env {
        name  = "DB_DRIVER"
        value = "sqlite"
      }

      env {
        name  = "DB_DSN"
        value = "/tmp/platform-service.db"
      }
    }
  }

  depends_on = [
    google_project_service.required["run.googleapis.com"],
    google_artifact_registry_repository.platform_service,
  ]
}