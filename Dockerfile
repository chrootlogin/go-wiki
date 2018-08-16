# Build build-container
FROM node:8-alpine as builder

ARG GO_VERSION=1.10.3
ARG DEP_VERSION=0.5.0

ADD https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz /tmp/golang.tar.gz

RUN set -ex \
    && apk -U --no-cache add \
        alpine-sdk \
        git \
        tar \
    && tar -C /usr/local -xzf /tmp/golang.tar.gz \
    && mkdir /lib64 \
    && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 \
    && mkdir -p /opt/go-wiki/src/github.com/chrootlogin/go-wiki

ADD https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 /usr/local/go/bin/dep

COPY ./ /opt/go-wiki/src/github.com/chrootlogin/go-wiki

WORKDIR /opt/go-wiki/src/github.com/chrootlogin/go-wiki

RUN set -ex \
    && chmod +x /usr/local/go/bin/dep \
    && export PATH="/usr/local/go/bin:/usr/local/node-v${NODE_VERSION}-linux-x64/bin:${PATH}" \
    && export GOPATH="/opt/go-wiki" \
    && export PATH="${GOPATH}/bin:${PATH}" \
    && go get -u github.com/go-bindata/go-bindata/... \
    && sync \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=1 make

# Build go-wiki container
FROM alpine:3.8

ARG BUILD_DATE
ARG VCS_REF

LABEL maintainer="Simon Erhardt <hello@rootlogin.ch>" \
  org.label-schema.name="Go-Wiki" \
  org.label-schema.description="Another wiki software written in Go." \
  org.label-schema.build-date=$BUILD_DATE \
  org.label-schema.vcs-ref=$VCS_REF \
  org.label-schema.vcs-url="https://github.com/chrootlogin/go-wiki" \
  org.label-schema.schema-version="1.0"

RUN set -ex \
    && apk -U --no-cache add \
        tini

ENV REPOSITORY_PATH="/data" \
    GIN_MODE="release"

RUN mkdir -p /opt/go-wiki \
  && mkdir /lib64 \
  && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY docker-run.sh /
COPY --from=builder /opt/go-wiki/src/github.com/chrootlogin/go-wiki/go-wiki /opt/go-wiki/go-wiki

EXPOSE 8080
VOLUME /data

ENTRYPOINT ["/sbin/tini","--"]
CMD ["/docker-run.sh"]
