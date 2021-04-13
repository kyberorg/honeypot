# Honeypot

[![Maintainability](https://api.codeclimate.com/v1/badges/51bc8dc67c396a7b87c4/maintainability)](https://codeclimate.com/github/kyberorg/honeypot/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/kyberorg/honeypot)](https://goreportcard.com/report/github.com/kyberorg/honeypot)

Honeypot listens for incoming ssh connections and writes the ip address, username, and password. 
This was written just for fun.

## Build

### Making binary
```shell
    CGO_ENABLED=0 go build github.com/kyberorg/honeypot/cmd/honeypot
```
or
```shell
   make binary
```

### Docker
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
