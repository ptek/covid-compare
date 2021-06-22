FROM alpine:3

RUN apk add --no-cache \
    tini \
    curl \
    wget \
    ca-certificates \
    bash \
    tzdata\
    enca

COPY ./do /opt/do
COPY ./scripts/ /opt/scripts/
COPY ./dist/bin/covid-compare_Linux /opt/dist/bin/covid-compare_Linux

WORKDIR /opt/scripts/
ENTRYPOINT ["tini", "/opt/do"]
CMD ["run"]