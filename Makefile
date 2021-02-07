
.PHONY: all clean build deps run test-e2e test  ui-docker-build ui-docker-run

OS := $(shell uname | tr '[:upper:]' '[:lower:]')
PWD:=$(shell pwd)

UI_PORT:=3000
NODE_IMAGE=node:8-alpine

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
	@rm -f ${BUILDDIR}/jb-chat

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
		-o ${BUILDDIR}/jb-chat ./cmd/

build: build-jb

run:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go run  \
    		-ldflags='-X "main.RELEASE=${RELEASE}" -X "main.COMMIT=${GITHASH}" -X "main.BUILDDATE=${BUILDDATE}"' \
    		./cmd/

ui-docker-run:
	docker run --publish '3000:3000' --rm --name l3-ui-run --volume "${PWD}/ui:/ui" \
			--volume "${PWD}/ui/node_modules/:/root/.npm/" \
			--workdir "/ui/" \
			"${NODE_IMAGE}" \
			sh -c 'npm run start'

ui-docker-build:
	docker run --rm --name l3-ui-build --volume "${PWD}/ui:/ui" \
			--volume "${PWD}/ui/node_modules/:/root/.npm/" \
			--workdir "/ui/" \
			"${NODE_IMAGE}" \
			sh -c 'npm run build'

build-container:

