# dev-span-pu-scaler [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[license]: https://github.com/nakatamixi/dev-span-pu-scaler/blob/master/LICENSE

`dev-span-pu-scaler` is an autoscaler of Cloud Spanner Processing Unit for develop environment usecase.

When your project uses Cloud Spanner in a develop environment,
the required number of Spanner nodes(processing unit) depends on the number of dbs.
(Unless your develop environment workload needs high cpu)
Cloud Spanner node cost is not low cost,
you need to keep minimum processing unit by the number of using dbs.

`dev-span-pu-scaler` counts current db number of target instances, calculate required processing unit number, and apply it to spanner instances automaticaly.

# CommandLine Usage
```
dev-span-pu-scaler -h
Usage of dev-span-pu-scaler:
  -buffer int
    	buffer db count to scale pu (default 3)
  -instances value
    	comma separated instances
  -project string
    	gcp project
  -server
    	run server(for Cloud Run)
```
You need `roles/spanner.admin` on your gcloud client.

# Run as sceduled job with Cloud Schadular + Cloud Run

If you want to run dev-span-pu-scaler as scheduled job,
you can deploy to Cloud Run and run as Cloud Sceduler job.

- add Service Account
```
resource "google_service_account" "dev-span-pu-scalar" {
  account_id   = "dev-span-pu-scalar"
  display_name = "dev-span-pu-scalar"
  description  = "sa for dev-span-pu-scalar"
}

variable "dev-span-pu-scalar-roles" {
  default = [
    "roles/errorreporting.writer",
    "roles/cloudtrace.agent",
    "roles/iam.serviceAccountUser",
    "roles/run.invoker",
    "roles/spanner.admin",
  ]
}

resource "google_project_iam_member" "dev-span-pu-scalar" {
  count  = length(var.dev-span-pu-scalar-roles)
  role   = element(var.dev-span-pu-scalar-roles, count.index)
  member = "serviceAccount:${google_service_account.dev-span-pu-scalar.email}"
}
```
- Deploy dev-span-pu-scalar to Cloud Run
you need [ko](https://github.com/google/ko) to build container image
```
# install ko
make build-tools
# see Makefile for detail
GOOGLE_CLOUD_PROJECT=<project> INSTANCES=<commma separated target instances> make deploy
```
- Add Cloud Scaduler job
```
data "google_cloud_run_service" "dev-span-pu-scaler" {
  name     = "dev-span-pu-scaler"
  location = var.location
}

resource "google_cloud_scheduler_job" "dev-span-pu-scaler" {
  name      = "dev-span-pu-scaler"
  schedule  = "0 * * * *"
  time_zone = var.time_zone

  http_target {
    http_method = "GET"
    uri         = "${data.google_cloud_run_service.dev-span-pu-scaler.status[0].url}/"

    oidc_token {
      service_account_email = google_service_account.dev-span-pu-scaler.email
      audience              = data.google_cloud_run_service.dev-span-pu-scaler.status[0].url
    }
  }
}
```
