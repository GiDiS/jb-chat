
.PHONY: all clean build deps run test-e2e test  ui-docker-build ui-docker-run \
	build-container-prod deploy-prod remove-prod run-prod \
	build-container-staging deploy-staging stop-staging run-staging

OS := $(shell uname | tr '[:upper:]' '[:lower:]')
PWD:=$(shell pwd)

UI_PORT:=3000
NODE_IMAGE:=node:lts-alpine
NS_PROD:=js-chat
NS_STAGING:=js-chat-staging

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

run: ui-docker-build
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
			--env 'NODE_ENV=production' \
			--volume "${PWD}/ui/node_modules/:/root/.npm/" \
			--workdir "/ui/" \
			"${NODE_IMAGE}" \
			sh -c 'npm run build'

build-container-prod:
	@eval $(minikube docker-env)
	docker build --build-arg 'APP_ENV=production' -f deploy/Dockerfile -t jb-chat-prod:0.1 .
	@minikube image load jb-chat-prod:0.1

deploy-prod:
	@eval $(minikube docker-env)
	test -n "$(shell kubectl get namespace "${NS_PROD}" -o jsonpath='{.metadata.uid}' --ignore-not-found=true)" \
		|| kubectl create namespace "${NS_PROD}"
	kubectl -n "${NS_PROD}" apply -f deploy/app-prod.yaml
	@echo -e "Services:"
	@minikube service list -n "${NS_PROD}"

run-prod: build-container-prod deploy-prod

stop-prod:
	test -z "$(shell kubectl get namespace "${NS_PROD}" -o jsonpath='{.metadata.uid}' --ignore-not-found=true)" \
		|| kubectl delete namespace "${NS_PROD}"

build-container-staging:
	@eval $(minikube docker-env)
	docker build --build-arg 'APP_ENV=staging' -f deploy/Dockerfile -t jb-chat-staging:0.1 .
	@minikube image load jb-chat-staging:0.1

deploy-staging:
	@eval $(minikube docker-env)
	test -n "$(shell kubectl get namespace "${NS_STAGING}" -o jsonpath='{.metadata.uid}' --ignore-not-found=true)" \
		|| kubectl create namespace "${NS_STAGING}"
	kubectl -n "${NS_STAGING}" apply -f deploy/app-staging.yaml
	@echo "Services:"
	@minikube service list -n "${NS_STAGING}"

run-staging: build-container-staging deploy-staging

stop-staging:
	test -z "$(shell kubectl get namespace "${NS_STAGING}" -o jsonpath='{.metadata.uid}' --ignore-not-found=true)"  \
		||  kubectl delete namespace jb-chat-staging
