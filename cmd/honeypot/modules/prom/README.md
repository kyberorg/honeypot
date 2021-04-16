# Prometheus metrics module

This module expose prometheus metrics. It exposes metrics at `http://host:2112/metric`

### Metrics
* Total number of connections (`honeypot_connections`)
* Number of unique sources (`honeypot_unique_sources`)

Prefix can be customised see [prefix param](#Prefix)

## Usage
### How to activate
```shell
--prom-metrics-enable
```

### Params

#### Port
* Custom Port (must be free port)
```shell
--prom-metrics-port=2112
```

#### Path
* Custom Metrics path 
```shell
--prom-metrics-path=/metrics
```

#### Prefix
* Custom metrics prefix
```shell
--prom-metrics-prefix=honeypot
```
