# K8s Intro

This is my first application built to run on Kubernetes. This app is a slightly
modified version of the
[guestbook-go](https://github.com/kubernetes/examples/tree/master/guestbook-go)
example. The differences between this repo and the original example is that this
repo uses `kustomize` to manage the resources required and uses the official
Redis client for Go.

# How to run

* Start a `minikube` cluster locally:

    ```bash
    minikube start
    ```

* Build the development container:

    ```bash
    docker build -t waduhek/k8s-intro:dev --target dev .
    ```

* Load the development image into minikube:

    ```bash
    minikube image load waduhek/k8s-intro:dev
    ```

* Apply the configuration with `kustomize`:

    ```bash
    kubectl apply -k config/overlays/local
    ```

* Teardown the configuration with:

    ```bash
    kubectl delete -k config/overlays/local
    ```
