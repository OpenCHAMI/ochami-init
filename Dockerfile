FROM cgr.dev/chainguard/wolfi-base

STOPSIGNAL SIGTERM

RUN apk add --no-cache tini

# nobody 65534:65534
USER 65534:65534

COPY ochami-init /ochami-init

# Set up the command to start the service.
VOLUME /config
WORKDIR /config

CMD /ochami-init

ENTRYPOINT ["/sbin/tini", "--"]
