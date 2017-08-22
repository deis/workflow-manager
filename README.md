
|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/156px-Warning.svg.png) | Deis Workflow will soon no longer be maintained.<br />Please [read the announcement](https://deis.com/blog/2017/deis-workflow-final-release/) for more detail. |
|---:|---|
| 09/07/2017 | Deis Workflow [v2.18][] final release before entering maintenance mode |
| 03/01/2018 | End of Workflow maintenance: critical patches no longer merged |

# Deis Workflow Manager

[![Build Status](https://travis-ci.org/deis/workflow-manager.svg?branch=master)](https://travis-ci.org/deis/workflow-manager) [![codecov](https://codecov.io/gh/deis/workflow-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/deis/workflow-manager)
 [![Go Report Card](https://goreportcard.com/badge/github.com/deis/workflow-manager)](https://goreportcard.com/report/github.com/deis/workflow-manager) [![codebeat badge](https://codebeat.co/badges/29e2c379-0490-45db-95fe-20b25bd5a466)](https://codebeat.co/projects/github-com-deis-workflow-manager)
[![Docker Repository on Quay](https://quay.io/repository/deis/workflow-manager/status "Docker Repository on Quay")](https://quay.io/repository/deis/workflow-manager)

This repository contains the manager component for Deis Workflow. Deis
(pronounced DAY-iss) Workflow is an open source Platform as a Service (PaaS)
that adds a developer-friendly layer to any [Kubernetes][k8s-home] cluster,
making it easy to deploy and manage applications on your own servers.

For more information about Deis Workflow, please visit the main project page at
https://github.com/deis/workflow.

We welcome your input! If you have feedback on Workflow Manager,
please [submit an issue][issues]. If you'd like to participate in development,
please read the "Development" section below and [submit a pull request][prs].

## Stay up to date

One of the primary goals for Workflow Manager is notifying operators of
component freshness. Workflow Manager will regularly check your cluster against
the latest stable components. If components are missing due to failure or are
simply out of date, Workflow operators will know at a glance.

By default, Workflow Manager will make version checks to an external service.
This submits component and version information to our versions service running
at [https://versions.deis.com](https://versions.deis.com). If you prefer this
check not happen, you may disable the function by setting
`WORKFLOW_MANAGER_CHECKVERSIONS` to `false` in the Workflow Manager's
Replication Controller.

## Workflow Doctor

Deployed closest to any potential problem, Workflow Manager is also designed to
help when things aren't going well. To aid troubleshooting efforts cluster
operators will be able to easily gather and securely submit cluster health and
status information to the Deis team.

Functionality will be added in a later release.

# Development

The Deis project welcomes contributions from all developers. The high level
process for development matches many other open source projects. See below for
an outline.

* Fork this repository
* Make your changes
* [Submit a pull request][prs] (PR) to this repository with your changes, and unit tests whenever possible
    * If your PR fixes any [issues][issues], make sure you write `Fixes #1234` in your PR description (where `#1234` is the number of the issue you're closing)
* The Deis core contributors will review your code. After each of them sign off on your code, they'll label your PR with `LGTM1` and `LGTM2` (respectively). Once that happens, a contributor will merge it

## Docker Based Development Environment

The preferred environment for development uses [the `go-dev` Docker
image](https://github.com/deis/docker-go-dev). The tools described in this
section are used to build, test, package and release each version of Deis.

To use it yourself, you must have [make](https://www.gnu.org/software/make/)
installed and Docker installed and running on your local development machine.

If you don't have Docker installed, please go to https://www.docker.com/ to
install it.

After you have those dependencies, bootstrap dependencies with `make bootstrap`,
build your code with `make build` and execute unit tests with `make test`.

## Native Go Development Environment

You can also use the standard `go` toolchain to build and test if you prefer.
To do so, you'll need [glide](https://github.com/Masterminds/glide) 0.9 or
above and [Go 1.6](http://golang.org) or above installed.

After you have those dependencies, you can build and unit-test your code with
`go build` and `go test $(glide nv)`, respectively.

Note that you will not be able to build or push Docker images using this method
of development.

# Testing

The Deis project requires that as much code as possible is unit tested, but the
core contributors also recognize that some code must be tested at a higher
level (functional or integration tests, for example).


[issues]: https://github.com/deis/workflow-manager/issues
[prs]: https://github.com/deis/workflow-manager/pulls
[k8s-home]: https://kubernetes.io
[v2.18]: https://github.com/deis/workflow/releases/tag/v2.18.0
