# Raw metrics module

This module provides application metrics. It writes metrics to stdout (application log) or to file.

### Metrics
* Total number of connections (`honeypot_connections`)
* Number of unique sources (`honeypot_unique_sources`)

Prefix can be customised see [prefix param](#Prefix)

## Usage
### How to activate
```shell
--with-raw-metrics
```

### Params

#### Prefix
* Custom metrics prefix
```shell
--raw-metrics-prefix=honeypot
```

#### File
* File to write metrics to. If present, module writes metrics to this file instead of stdout.
```shell
--raw-metrics-file=/path/to/file
```
