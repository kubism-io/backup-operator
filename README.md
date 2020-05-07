# Backup Operator

[![Go Documentation](https://img.shields.io/badge/go-doc-blue.svg?style=flat)](https://pkg.go.dev/mod/github.com/kubism/backup-operator?tab=packages)
[![Build Status](https://travis-ci.org/kubism/backup-operator.svg?branch=master)](https://travis-ci.org/kubism/backup-operator)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubism/backup-operator)](https://goreportcard.com/report/github.com/kubism/backup-operator)
[![Coverage Status](https://coveralls.io/repos/github/kubism/backup-operator/badge.svg?branch=master)](https://coveralls.io/github/kubism/backup-operator?branch=master)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/kubismio/backup-operator.svg?sort=semver)](https://hub.docker.com/r/kubismio/backup-operator/tags)

## Usage

### Setup

Find the helm chart for the backup-operator at the [Kubism.io Helm Charts](https://kubism.github.io/charts/#chart-backup-operator).

TODO

### Backups for MongoDB

TODO

### Backups for Consul

The backup of consul to S3 is supported at the moment. See example configuration in [`backup_v1alpha1_consulbackupplan.yaml`](./config/samples/backup_v1alpha1_consulbackupplan.yaml).

## Design

TODO

## Development

### Tools

All required tools for development are automatically downloaded and stored in the `tools` sub-directory (see relevant section of [`Makefile`](./Makefile) for details).
A custom [`tools/goget-wrapper`](./tools/goget-wrapper) is used to create a temporary isolated environment to install and compile go tools.
To make sure those can be properly used in tests, several helpers were implemented in [`pkg/testutil`](./pkg/testutil) (e.g. `HelmEnv`, `KindEnv`).

### Testing

The tests depend on `docker` and `kind` and use [`ginkgo`](https://github.com/onsi/ginkgo) and [`gomega`](https://github.com/onsi/gomega). To spin up containers for tests [`ory/dockertest`](https://github.com/ory/dockertest) is used. For the controller tests `kind` is used, which has the advantage, compared to the more lightweight kubebuilder assets approach, to properly handle finalizers and allow integration tests.


### Kubebuilder

This project uses a different project layout than what is generated by
kubebuilder. The layout adheres to the [golang standards](https://github.com/golang-standards/project-layout) layout.
For this to properly work a wrapper is required ([`tools/kubebuilder-wrapper`](./tools/kubebuilder-wrapper)),
which makes sure the correct kubebuilder version is available and temporarily
moves files around as required.

While this is certainly not beautiful, this should improve with future versions
of kubebuilder and their plugin capabilities.

#### Known quirks

* When using the kubebuilder CLI to create a new API [`main.go`](./cmd/manager/main.go)
has a wrong controllers import path and has to be fixed manually afterwards.

### Extending the operator

To extend the operator you have to use the wrapper ([`tools/kubebuilder-wrapper`](./tools/kubebuilder-wrapper)) to scaffold out a new [Kind](https://book.kubebuilder.io/cronjob-tutorial/gvks.html#kinds-and-resources) and corresponding controller. The following command (see the official [kubebuilder docs](https://book.kubebuilder.io/cronjob-tutorial/new-api.html)) must be used:

```bash
./tools/kubebuilder create api --group backup --version <version> --kind <SomeBackupPlan>
```

Using the command above will generate several classes for you:

* `api/<version>/<somebackupplan>_types.go`
  * Add the spec and other necessary stuff for your new Kind here
* `controllers/<somebackupplan>_controller.go`
  * Please move the file to `pkg/controllers/<somebackupplan>_controller.go`
  * Implement the controller for your new Kind

In addition to the operator specifics you have to implement a new command as part of the worker below `cmd/worker` like for the existing ones, e.g. [`mongodb`](cmd/worker/mongodb.go).

Please add tests for all new parts added to the operator.