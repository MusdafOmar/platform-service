resource "google_storage_bucket" "terraform_state" {
  name          = "${var.project_id}-terraform-state"
  project       = var.project_id
  location      = var.region
  storage_class = "STANDARD"

  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"
  force_destroy               = false

  versioning {
    enabled = true
  }

  labels = local.common_labels

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [
    google_project_service.required["storage.googleapis.com"],
  ]
}
