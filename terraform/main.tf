variable "project_id" {
  type = string
  default = "summary-bot-server"
}

# Google Cloud providerを設定
provider "google" {
  project = var.project_id
  region = "asia-northeast1"
}

resource "google_artifact_registry_repository" "bot-server" {
  description = "bot-server-repository"
  location = "asia-northeast1"
  repository_id = "bot-server"
  format = "DOCKER"
}

# Cloud Run serviceを作成
resource "google_cloud_run_service" "bot-server" {
  name = "bot-server"
  location = "asia-northeast1"

  # Dockerイメージを指定
  template {
    spec {
      containers {
        image = "${google_artifact_registry_repository.bot-server.id}:${git_revision}"
      }
    }
  }
  depends_on = [
    google_artifact_registry_repository.bot-server
  ]
}

resource "google_cloud_run_service_iam_policy" "bot-server" {
  policy_data = jsonencode({
    bindings = [
      {
        role = "roles/run.invoker"
        members = [
          "allUsers"]
      }
    ]
  })
  location = google_cloud_run_service.bot-server.location
  project = google_cloud_run_service.bot-server.project
  service = google_cloud_run_service.bot-server.name

  depends_on = [
    google_cloud_run_service.bot-server
  ]
}

# Cloud Run serviceにアクセスするためのURLを出力
output "url" {
  value = google_cloud_run_service.bot-server.status[0].url
}
