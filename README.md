<!-- TITLE -->
<br />
<div align="center">
  <img src="./assets/logo.png" alt="Logo" width="128" height="128">
  <h1 align="center">deprek8</h3>
  <p align="center">
    A CLI tool that allows to list deprecated deployed Kubernetes resources from Helm or simply k8s objects.
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

`deprek8` is a CLI tool that :
- List deployed Helm Charts containing deprecated Kubernetes objects
- List deprecated Kubernetes resources deployed in a cluster

## Usage

- ### Helm
List deprecated Helm charts resources :
```
deprek8 helm [--output outputType] [--export]

-o, --output string              Choose type of output (json|yaml|text) (default: "text")
-e, --export bool                Save output to csv file (default: false)
```

- ### Kubernetes
List deprecated Kubernetes resources :
```
deprek8 kube [--output outputType] [--export]

-o, --output string              Choose type of output (json|yaml|text) (default: "text")
-e, --export bool                Save output to csv file (default: false)
```

## Setup
- ### Build locally :
  ```
  make build
  ```

- ### Build and install :
  ```
  make install
  ```
  > This command will build and install the binary into `/usr/local/bin`. It will also create a `.deprek8` directory in the user's home directory with the `config.yaml` config file.


- ### Uninstall :
  ```
  make uninstall
  ```
