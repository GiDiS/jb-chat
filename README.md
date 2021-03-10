# Readme

jb-chat - test chat with react+mobx at frontend and golang and backend and WebSocket for communications between

## Install

```shell
go get -u github.com/GiDiS/jb-chat
```

## Run at host

Staging: Seeded with GoT data set

```shell
cd "${GOPATH}/src/github.com/GiDiS/jb-chat" && make run-host-staging
```

Production: Run with empty database

```shell
cd "${GOPATH}/src/github.com/GiDiS/jb-chat" && make run-host-prod
```

Listen at http://localhost:8888 for ui
Listen at http://localhost:8889 for diagnostic

## Run in minikube with kubectl (minikube required)

Staging: Seeded with GoT data set

```shell
cd "${GOPATH}/src/github.com/GiDiS/jb-chat" && make run-kctl-staging
```

Production: Run with empty database

```shell
cd "${GOPATH}/src/github.com/GiDiS/jb-chat" && make run-kctl-production
```

Service and ingress urls printed after deploy


## Run in minikube with k8s-handle ([k8s-handle](https://github.com/2gis/k8s-handle) required)

k8s-handle config locates in **deploy/k8s-handle/config.yaml**

Staging: Seeded with GoT data set

```shell
cd "${GOPATH}/src/github.com/GiDiS/jb-chat" && make run-k8sh-staging
```

Production: Run with empty database

```shell
cd "${GOPATH}/src/github.com/GiDiS/jb-chat" && make run-k8sh-production
```

Service and ingress urls printed after deploy

## Stop in minikube (minikube required)
Stop completely remove deploy namespace 

Staging:
```shell
cd "${GOPATH}/src/github.com/GiDiS/jb-chat" && make stop-staging
```

Production:
```shell
cd "${GOPATH}/src/github.com/GiDiS/jb-chat" && make stop-production
```



## Issues (todo):

* UI: reconnect to ws after backend restart
* UI: refresh after channel create, join, leave
* UI: load channel messages on page initial open by direct url
* Back: Persistent storage
* Back: Tests
* Back: Metrics
* ~~Deploy: helm chart~~ Used k8s-handle
* ~~Deploy: minikube deploy with ingress for stable google sign-in~~

## App arch

![App arch diagram](docs/arch.png)
