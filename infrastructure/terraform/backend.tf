terraform {
  backend "gcs" {
    bucket = "lia-platform-musdaf-2026-terraform-state"
    prefix = "platform-service/development"
  }
}
