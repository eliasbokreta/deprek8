<!-- TITLE -->
<br />
<div align="center">
  <img src="./assets/logo.png" alt="Logo" width="128" height="128">
  <h1 align="center">deprek8</h3>
  <p align="center">
    A CLI tool that allows to list deprecated Kubernetes resources Helm Chart and deployed Kubernetes objects.
  </p>
</div>

<div align="center">

[![CI workflow](https://github.com/eliasbokreta/deprek8/actions/workflows/main.yml/badge.svg)](https://github.com/eliasbokreta/deprek8/actions/workflows/main.yml/badge.svg)

</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#usage">Usage</a>
    </li>
    <li>
      <a href="#setup">Setup</a>
    </li>
  </ol>
</details>


---

## About The Project

`deprek8` is a CLI tool that allows to :
- List deployed Helm Charts with information about :
  - Latest released chart version found in [Artifacthub](https://artifacthub.io/packages/search?kind=0&sort=relevance&page=1)
  - Deprecated Kubernetes resource in deployed charts
- List deprecated Kubernetes resources deployed in a cluster

## Usage

- ### Helm
List deployed Helm charts :
```
deprek8 helm [--output outputType] [--all-namespaces] [--filter-deprecated] [--filter-name chartName] [--filter-repository repositoryName]

-o, --output string              Choose type of output (json|yaml|text) (default "text")
-a, --all-namespaces             Fetch data on all namespaces
-d, --filter-deprecated          Filter helm releases with deprecated k8s resources
-n, --filter-name string         Filter Helm chart's name
```

- ### Kubernetes
List deprecated Kubernetes resources :
```
deprek8 kube [--output outputType]

-o, --output string              Choose type of output (json|yaml|text) (default "text")
```

Switch Kubernetes context :
```
deprek8 context
```

## Setup
Build locally :
```
make build
```

Update the binary to the latest version :
```
deprek8 update
```
