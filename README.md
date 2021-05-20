# Honeypot

[![Maintainability](https://api.codeclimate.com/v1/badges/51bc8dc67c396a7b87c4/maintainability)](https://codeclimate.com/github/kyberorg/honeypot/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/kyberorg/honeypot)](https://goreportcard.com/report/github.com/kyberorg/honeypot)

A simple SSH honeypot written on Go. Strictly not a honeypot as it doesn't trap or jail anything, it simply collects data on attempts to login to a generic SSH server open to the internet.

The tool runs an SSH server that rejects all login attempts. There is no session created it just allows a login attempt and records the username and password and source IP for later analysis.

## Build

### Making binary
```shell
   make binary
```
or direct way if `make` not an option for you
```shell
    CGO_ENABLED=0 go build github.com/kyberorg/honeypot/cmd/honeypot
```

### Docker
See our [DockerHub Repo](https://hub.docker.com/repository/docker/kyberorg/honeypot)
```shell
   docker pull kyberorg/honeypot:tagname
```

## Run
```shell
   # (Optionally) creating host key
   ssh-keygen -t rsa -f honeypot.id_rsa
   # Run it
   bin/honeypot --hostkey honeypot.id_rsa
```

### Help
```shell
   bin/honeypot --help
```

## GeoIP
GeoIP enriches access log with geoip information (city, region, country) based on connection IP.

[GeoIP Readme](cmd/honeypot/geoip/README.md)

```shell
--geoip-mmdb-file=/path/to/GeoLite2-City.mmdb
```

## Modules
### Prometheus Metrics Module
Module that exposes prometheus metrics.

[Module Readme](cmd/honeypot/mod/prom/README.md)

```shell
--with-prom-metrics
```

### Raw metrics module

Module that provides application metrics. It writes metrics to stdout (application log) or to file.
[Module Readme](cmd/honeypot/mod/rawmetrics/README.md)
```shell
--with-raw-metrics
```

## Systemd Daemon
* Simple Example
```unit file (systemd)
[Unit]
Description=Fake SSH
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
Restart=always
RestartSec=5s
Environment="ACCESS_LOG=connections.log"
WorkingDirectory=/srv/honeypot
ExecStart=/srv/honeypot/honeypot \
    --hostkey=honeypot.id_rsa \
    --geoip-mmdb-file=/var/lib/GeoIP/GeoLite2-City.mmdb \
    --prom-metrics-enable

SyslogIdentifier=honeypot

[Install]
WantedBy=multi-user.target
```

## Projekt Logo
![Logo](logo.png?raw=true)
