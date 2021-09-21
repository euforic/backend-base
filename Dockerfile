FROM golang:latest AS build_base

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/euforic/backend-base

COPY . .

RUN go get -v github.com/matr-builder/matr
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
     wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
     chmod +x /bin/grpc_health_probe

# Unit tests
RUN [ "matr", "build" ]

# Graphql API target
FROM scratch as backendgql

COPY --from=build_base /go/src/github.com/euforic/backend-base/build/backendgql .

EXPOSE 8080
CMD ["./backendgql"]

# gRPC server target
FROM alpine as backendgrpc

RUN addgroup -g 1000 rpc; adduser -D -u 1000 -G rpc rpc

COPY --from=build_base --chown=rpc:rpc /bin/grpc_health_probe ./grpc_health_probe
COPY --from=build_base /go/src/github.com/euforic/backend-base/build/backendgrpc .

RUN chmod +x /grpc_health_probe

EXPOSE 9000
CMD ["./backendgrpc"]
