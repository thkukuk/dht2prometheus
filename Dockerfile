FROM registry.opensuse.org/opensuse/tumbleweed:latest AS build-stage
RUN zypper install --no-recommends --auto-agree-with-product-licenses -y git go make
COPY . dht2prometheus
RUN cd dht2prometheus && make update && make tidy && make

FROM registry.opensuse.org/opensuse/busybox:latest
LABEL maintainer="Thorsten Kukuk <kukuk@thkukuk.de>"

ARG BUILDTIME=
ARG VERSION=unreleased
LABEL org.opencontainers.image.title="Exports DHTxx values as metrics for Prometheus"
LABEL org.opencontainers.image.description="Exports DHTxx (DHT11,DHT12,DHT22) sensor values as metrics for Prometheus"
LABEL org.opencontainers.image.created=$BUILDTIME
LABEL org.opencontainers.image.version=$VERSION

COPY --from=build-stage /dht2prometheus/bin/dht2prometheus /usr/local/bin
COPY entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/usr/local/bin/dht2prometheus"]
