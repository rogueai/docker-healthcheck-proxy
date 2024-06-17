# Compile stage
FROM golang:1 AS build

WORKDIR /go/src/
COPY . .
RUN go mod vendor
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2

RUN make


# Final stage
FROM alpine:3

WORKDIR /app
COPY --from=build /go/src/target/healthcheck /app

ENTRYPOINT ["/app/healthcheck"]
