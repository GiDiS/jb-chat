
.PHONY: all clean build deps  test-e2e test  ui-docker-build ui-docker-run \
	run run-host-prod run-host-staging \
	build-container-prod stop-prod deploy-kctl-prod run-kctl-prod deploy-k8sh-prod run-k8sh-prod \
	build-container-staging stop-staging deploy-kctl-staging run-kctl-staging deploy-k8sh-staging run-k8sh-staging 

OS := $(shell uname | tr '[:upper:]' '[:lower:]')
PWD:=$(shell pwd)

UI_PORT:=3000
NODE_IMAGE:=node:lts-alpine
NS_PROD:=jb-chat
NS_STAGING:=jb-chat-staging
IMAGE_NAME:=jb-chat
IMAGE_VERSION:=0.1

%-prod: APP_ENV = production
%-prod: IMAGE_NAME = jb-chat-prod
%-prod: IMAGE_VERSION = 0.1
%-prod: NS := $(NS_PROD)

%-staging: APP_ENV = staging
%-staging: IMAGE_NAME = jb-chat-staging
%-staging: IMAGE_VERSION = 0.1
%-staging: NS := $(NS_STAGING)

GIT_COMMIT := $(shell git rev-parse --short=7 HEAD)

BUILDDIR:=$(realpath .)
BUILDDATE:=$(shell date --rfc-3339=seconds)

GITHASH:=$(shell git log --pretty=format:'%h' -n 1)
ifeq ("${GITHASH}", "")
	GITHASH:=unknown
endif

RELEASE:=$(shell git describe --tags 2>/dev/null)
ifeq ("${RELEASE}", "")
	RELEASE:=unknown
endif

all: clean deps test build

prepare:
	@echo "Installing modules... "

clean:
	@echo "Clean ... "
	@rm -f ${BUILDDIR}/jb_chat

deps:
	@echo "Setting up the vendors folder... ${GOPATH}"
	go mod tidy

check:
	go vet ./...

test-e2e:
	go test -tags integration -cover -race ./test/e2e

test:
	go generate ./...
	go test -cover -race ./...
	#test-e2e


benchmark:
	@mkdir -p logs
	@touch logs/${GIT_COMMIT}.out
	@go test -run none -bench . -benchmem ./... >> logs/${GIT_COMMIT}.out


build-jb:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  \
		-ldflags='-X "main.RELEASE=${RELEASE}" -X "main.COMMIT=${GITHASH}" -X "main.BUILDDATE=${BUILDDATE}"' \
		-o ${BUILDDIR}/jb-chat ./cmd/chatd/main.go

build: build-jb

run-host-prod: ui-docker-build
	go generate ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go run  \
    		-ldflags='-X "main.RELEASE=${RELEASE}" -X "main.COMMIT=${GITHASH}" -X "main.BUILDDATE=${BUILDDATE}"' \
    		./cmd/chatd/main.go

run: run-host-prod

run-host-staging: export SEED = 1
run-host-staging: export APP_ENV = staging
run-host-staging: ui-docker-build
	go generate ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go run  \
    		-ldflags='-X "main.RELEASE=${RELEASE}" -X "main.COMMIT=${GITHASH}" -X "main.BUILDDATE=${BUILDDATE}"' \
    		./cmd/chatd/main.go

ui-docker-run:
	docker run --publish '3000:3000' --rm --name jb-ui-run --volume "${PWD}/ui:/ui" \
			--env 'API_SERVER=http://localhost:8888' \
			--volume "${PWD}/ui/node_modules/:/root/.npm/" \
			--workdir "/ui/" \
			"${NODE_IMAGE}" \
			sh -c 'npm run start'

ui-docker-build:
	docker run --rm --name jb-ui-build --volume "${PWD}/ui:/ui" \
			--env "NODE_ENV=staging" \
			--volume "${PWD}/ui/node_modules/:/root/.npm/" \
			--workdir "/ui/" \
			"${NODE_IMAGE}" \
			sh -c "REACT_APP_ENV=staging npm run build"

build-container:
	@eval $(minikube docker-env)
	docker build --build-arg "APP_ENV=$(APP_ENV)" -f deploy/Dockerfile -t "${IMAGE_NAME}:${IMAGE_VERSION}" .
	@minikube image load "${IMAGE_NAME}:${IMAGE_VERSION}"

stop-kctl:
	test -z "$(shell kubectl get namespace "$(NS)" -o jsonpath='{.metadata.uid}' --ignore-not-found=true)" \
    		|| kubectl delete namespace "$(NS)"

deploy-kctl:
	@eval $(minikube docker-env)
	test -n "$(shell kubectl get namespace "$(NS)" -o jsonpath='{.metadata.uid}' --ignore-not-found=true)" \
		|| kubectl create namespace "$(NS)"
	kubectl -n "$(NS)" apply -f "deploy/app-$(APP_ENV).yaml"
	@echo "Services:"
	@minikube service list -n "$(NS)"
	$(eval HOST := $(shell kubectl -n "$(NS)" get ingress jb-chat-ingress --output jsonpath='{.spec.rules[0].host}'))
	@echo "Add to /etc/hosts to use ingress: \033[32m $(shell minikube ip) $(HOST)\033[0m"
	@echo "Try with ingress: \033[32m https://$(HOST)\033[0m"

deploy-k8sh:
	@eval $(minikube docker-env)
	test -n "$(shell kubectl get namespace $(NS) -o jsonpath='{.metadata.uid}' --ignore-not-found=true)" \
		|| kubectl create namespace $(NS)
	cd ./deploy/k8s-handle && IMAGE_VERSION=$(IMAGE_VERSION) \
		k8s-handle deploy -s "$(APP_ENV)" --use-kubeconfig --sync-mode -c config.yaml
	@echo "Services:"
	@minikube service list -n $(NS)
	$(eval HOST := $(shell kubectl -n "$(NS)" get ingress jb-chat-ingress --output jsonpath='{.spec.rules[0].host}'))
	@echo "Add to /etc/hosts to use ingress: \033[32m $(shell minikube ip) $(HOST)\033[0m"
	@echo "Try with ingress: \033[32m https://$(HOST)\033[0m"

build-container-prod: build-container
build-container-staging: build-container

deploy-kctl-prod: deploy-kctl
run-kctl-prod: build-container deploy-kctl

deploy-kctl-staging: deploy-kctl
run-kctl-staging: build-container deploy-kctl

deploy-k8sh-prod: deploy-k8sh
run-k8sh-prod: build-container deploy-k8sh

deploy-k8sh-staging: deploy-k8sh
run-k8sh-staging: build-container deploy-k8sh


stop-prod: stop-kctl
stop-staging: stop-kctl

