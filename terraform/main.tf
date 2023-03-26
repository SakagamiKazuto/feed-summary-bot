variable "project_id" {
  type = string
  default = "summary-bot-server"
}

terraform {
  backend "gcs" {
    bucket = "tfstate-summary-bot-server"
    prefix = "terraform/state"
  }
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
        //image = "${google_artifact_registry_repository.bot-server.id}:latest"
        image = "us-docker.pkg.dev/cloudrun/container/hello"
        env {
          name = "SLACK_APP_TOKEN"
          value_from {
            secret_key_ref {
              name = "SLACK_APP_TOKEN"
              key = "latest"
            }
          }
        }
          env {
            name = "SLACK_SIGNING_TOKEN"
            value_from {
              secret_key_ref {
                name = "SLACK_SIGNING_TOKEN"
                key  = "latest"
              }
            }
        }
        env {
          name = "OPENAI_API_KEY"
          value_from {
            secret_key_ref {
              name = "OPENAI_API_KEY"
              key  = "latest"
            }
          }
        }
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
