provider "google" {
  project = var.project_id
}

# Grant roles to the developer Google Group
locals {
  dev_roles = [
    "roles/run.developer",
    "roles/datastore.user",
    "roles/logging.viewer",
    "roles/iam.serviceAccountUser" # Required for deploying to Cloud Run
  ]
}

resource "google_project_iam_member" "dev_team_roles" {
  for_each = toset(local.dev_roles)

  project = var.project_id
  role    = each.value
  member  = "group:${var.dev_group_email}"
}
