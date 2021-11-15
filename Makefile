APP?=logi2
RELEASE?=0.1.0

MEPATH=$(shell pwd)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: build
build: clean
	go build -ldflags "-X main.version=${RELEASE} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" -o ${APP} .

.PHONY: dev
dev: tst build
	@$(MEPATH)/logi2 -p 15000

.PHONY: tst
tst:
	go clean -testcache
	go test ./...

.PHONY: clean
clean:
	@rm -fR logi2

