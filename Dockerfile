FROM debian:buster-slim 

RUN apt-get update -y && apt-get install -y ca-certificates openssl && \
      update-ca-certificates 2>/dev/null || true

RUN mkdir -p /usr/local/bin/spiderweb/static

COPY ./spiderweb /usr/local/bin/spiderweb/
COPY ./static/* /usr/local/bin/spiderweb/static

ENTRYPOINT /usr/local/bin/spiderweb/spiderweb
