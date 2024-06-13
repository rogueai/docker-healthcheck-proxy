GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_TAG := $(shell git describe --tags 2>/dev/null)

LDFLAGS := -s -w -X github.com/mala-cimbra/docker-healthcheck-proxy.commit=${GITCOMMIT}
LDFLAGS := ${LDFLAGS} -X github.com/mala-cimbra/docker-healthcheck-proxy.tag=${GIT_TAG}
OUTFILE ?= healthcheck


target:
	mkdir target

.PHONY: clean
clean:
	rm -rf target

.PHONY: compile
compile: target
	go build -ldflags "${LDFLAGS}" -o target/${OUTFILE} main.go
	gzip -c < target/${OUTFILE} > target/${OUTFILE}.gz
