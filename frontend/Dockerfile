FROM alpine

RUN apk update && apk upgrade && \
    apk add wget && wget https://dl.google.com/go/go1.21.4.linux-amd64.tar.gz && \
    tar -C /usr/local/ -xzf go1.21.4.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin