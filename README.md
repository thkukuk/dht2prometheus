# dht2promehteus
**Exports DHTxx sensor values as metrics for Prometheus**


This application reads temperature and humidity values of a DHTxx sensor (DHT11, DHT12, DHT22 or AM2302) using the [go-dht](https://github.com/d2r2/go-dht) library and provides them as metrics for [Prometheus](https://prometheus.io).

## Container

### Public Container Image

To access the DHTxx sensor GPIO interface we need write access to the `/sys` filesystem. This is only possible if the container runs in privileged mode.

The command to run the public available image would be:

```bash
podman run --privileged -p 9420:9420 -v <path>/config.yaml:/config.yaml registry.opensuse.org/home/kukuk/containerfile/dht2prometheus:latest
```

You can replace `podman` with `docker` without any further changes.

### Build locally

To build the container image with the `dht2prometheus` binary included run:

```bash
sudo podman build --rm --no-cache --build-arg VERSION=$(cat VERSION) --build-arg BUILDTIME=$(date +%Y-%m-%dT%TZ) -t dht2prometheus .
```

You can of cource replace `podman` with `docker`, no other arguments needs to be adjusted.

## Configuration

dht2prometheus can be configured via command line and configuration file.

### Commandline

Available options are:
```plaintext
Usage:
  dht2prometheus [flags]

Flags:
  -c, --config string   configuration file (default "config.yaml")
  -h, --help            help for dht2prometheus
  -q, --quiet           don't print any informative messages
  -v, --verbose         become really verbose in printing messages
      --version         version for dht2prometheus
```

### Configuration File

By default `dht2prometheus` looks for the file `config.yaml` in the local directory. This can be overriden with the `--config` option.

Here is my configuration file, which I use for DHT22 connected on a Raspberry Pi 3 running [openSUSE MicroOS](https://microos.opensuse.org)

```yaml
# Descriptive name of the sensor provided as label in the metric
name: "Living room"
# Required: sensor type, valid values are DHT11, DHT12, DHT22 or AM2302
sensor: DHT22
# Required: GPIO pin on which the device is connected. This depends on
# the kernel and configuration
gpio_pin: 17
# Optional: address and port to listen on, default is port 9420
listen: ":9420"
# Optional: temperature offset
#temperature_offset: -1.7
# Optional: humidity offset
#humidity_offset: +0.3
```
