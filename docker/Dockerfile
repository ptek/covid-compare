FROM alpine:3

RUN apk add --no-cache \
    coreutils \
    tini \
    curl \
    wget \
    ca-certificates \
    bash \
    tzdata\
    enca \
    openssl

COPY ./do /opt/do
COPY ./scripts/*.sh /opt/scripts/
COPY ./dist/bin/covid-compare_Linux /opt/dist/bin/covid-compare_Linux

WORKDIR /opt/
ENTRYPOINT ["tini", "/opt/do"]
CMD ["run"]