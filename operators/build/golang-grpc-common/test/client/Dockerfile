# Common golang grpc Dockerfile -- SERVER

FROM golang:1.21 as builder 

WORKDIR /tmp/builder

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download && go mod verify

# Set which component to be installed
ARG COMPONENT
ARG DOMAIN
RUN test -n "$COMPONENT" || ( echo "The COMPONENT argument is unset. Aborting" && false ) && \
    test -n "$DOMAIN" || ( echo "The DOMAIN argument is unset. Aborting" && false )

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$(go env GOARCH) go build -ldflags="-s -w -X main.domain=$DOMAIN" ./cmd/$COMPONENT/grpc/client

###

FROM alpine:3.14

RUN apk update

ARG COMPONENT

COPY --from=builder /tmp/builder/client /usr/bin/client
COPY --from=builder /tmp/builder/test/sbatch/test_script.sh /usr/bin/test_script.sh

RUN ln -s /usr/bin/client /usr/bin/bioeda-component

# ENTRYPOINT [ "/usr/bin/bioeda-component", "-f", "/usr/bin/test_script.sh" ]
# ENTRYPOINT [ "/usr/bin/bioeda-component"]