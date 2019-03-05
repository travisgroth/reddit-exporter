## Introduction

This chart creates a reddit-exporter deployment on a [Kubernetes](http://kubernetes.io)
cluster using the [Helm](https://helm.sh) package manager.

## Adding chart repo

```bash
$ helm repo add reddit-exporter https://reddit-exporter-chart.storage.googleapis.com
```

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm install --name my-release reddit-exporter/reddit-exporter --set subreddits="{wtf,askreddit}"
```

The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
$ helm delete --purge my-release
```
The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Hubot chart and their default values.

Parameter | Description | Default
--- | --- | ---
`namespace` | namespace to use | `default`
`replicaCount` | desired number of pods | `1`
`image.repository` | container image repository | `minddocdev/hubot`
`image.tag` | container image tag | `latest`
`image.pullPolicy` | container image pull policy | `Always`
`service.type` | type of service to create | `NodePort`
`service.httpPort` | port for the http service | `80`
`resources` | resource requests & limits | `{}`
`nodeSelector` | node selector logic | `{}`
`tolerations` | resource tolerations | `{}`
`affinity` | node affinity | `{}`
`verbose` | verbose logging mode | `false`
`regexMatches` | regexFile content in regexMatchesFormat format | `false`
`regexMatchesFormat` | format of regexMatches | `yaml`
`subreddits` | REQUIRED list of subreddits to monitor | `[]`
