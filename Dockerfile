# Compile stage
FROM golang:1.22 AS build

ADD . /build
WORKDIR /build

RUN go build -o /healthcheck

# Final stage
FROM golang:1.22

WORKDIR /
COPY --from=build /healthcheck /

ENTRYPOINT ["/healthcheck"]