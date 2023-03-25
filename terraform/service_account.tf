resource "google_service_account" "github-actions" {
  account_id   = "github-actions"
  display_name = "GitHub Actions Service Account"
}

resource "google_project_iam_custom_role" "github-actions" {
  role_id     = "github_actions"
  title       = "GitHub Actions Custom Role"
  description = "Custom role for GitHub Actions"

  permissions = [
    "cloudbuild.builds.create",
    "cloudbuild.builds.get",
    "cloudbuild.builds.list",
    "cloudbuild.builds.update",
    "storage.buckets.get",
    "storage.buckets.list",
    "storage.buckets.create",
    "storage.buckets.delete",
    "storage.objects.create",
    "storage.objects.delete",
    "storage.objects.get",
    "storage.objects.list",
    "storage.objects.update",
    "roles/run.admin",
  ]
}

resource "google_project_iam_member" "github-actions" {
  project = var.project_id
  role    = google_project_iam_custom_role.github-actions.id
  member  = "serviceAccount:${google_service_account.github-actions.email}"
}
