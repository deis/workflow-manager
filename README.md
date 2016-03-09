# Deis Workflow Manager

[![Build Status](https://travis-ci.org/deis/workflow-manager.svg?branch=master)](https://travis-ci.org/deis/workflow-manager) [![Go Report Card](https://goreportcard.com/badge/github.com/deis/workflow-manager)](https://goreportcard.com/report/github.com/deis/workflow-manager)

Deis (pronounced DAY-iss) is an open source PaaS that makes it easy to deploy and manage
applications on your own servers. Deis builds on [Kubernetes](http://kubernetes.io/) to provide
a lightweight, [Heroku-inspired](http://heroku.com) workflow.

## Work in Progress

![Deis Graphic](https://s3-us-west-2.amazonaws.com/get-deis/deis-graphic-small.png)

Deis Workflow v2 is currently in alpha. Your feedback and participation are more than welcome, but be
aware that this project is considered a work in progress.

# Component Description

Workflow Manager Service is responsible for monitoring cluster health.

## Stay up to date

One of the primary goals for Workflow Manager is notifying operators of
component freshness. Workflow Manager will regularly check your cluster against
the latest stable components. If components are missing due to failure or are
simply out of date, Workflow operators will know at a glance.

## Workflow Doctor

Deployed closest to any potential problem, Workflow Manager is also designed to
help when things aren't going well. To aid troubleshooting efforts cluster
operators will be able to easily gather and securely submit cluster health and
status information to the Deis team.
