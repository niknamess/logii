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

.PHONY: dev2
dev2: tst build
	@$(MEPATH)/logi2 -p 16000

.PHONY: vfc
vfc: tst build
	@$(MEPATH)/logi2 -v 10015

.PHONY: gen
gen: tst build
	@$(MEPATH)/logi2 -g Foxtrot Uniform Charlie Kilo

.PHONY: rm
rm: build
	@$(MEPATH)/logi2 -r Foxtrot Uniform Charlie Kilo

.PHONY: ip
ip: tst build
	@$(MEPATH)/logi2 -i Foxtrot Uniform Charlie Kilo

.PHONY: tst
tst:
	go clean -testcache
	go test ./...

.PHONY: clean
clean:
	@rm -fR logi2

