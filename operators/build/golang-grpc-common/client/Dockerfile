# Common golang grpc Dockerfile -- SERVER

FROM golang:1.21 as builder 

WORKDIR /tmp/builder

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download && go mod verify

# Set which component to be installed
ARG COMPONENT
RUN test -n "$COMPONENT" || ( echo "The COMPONENT argument is unset. Aborting" && false )

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$(go env GOARCH) go build -ldflags="-s -w" ./cmd/$COMPONENT/client

###

FROM alpine:3.14

RUN apk update

ARG COMPONENT

COPY --from=builder /tmp/builder/client /usr/bin/$COMPONENT

RUN ln -s /usr/bin/$COMPONENT /usr/bin/bioeda-component

ENTRYPOINT [ "/usr/bin/bioeda-component" ]