# wormhole

Open loop benchmark tool, inspired by [YCSB Is Obsolete, We Need New Benchmarks](https://emptysqua.re/blog/ycsb-is-obsolete/). Currently only supports the etcd database and a write-only workload.

## Usage

```
./bin/wormhole -clients=10 -exec-time=10s -cmd-timeout=300ms -max-thinking-time=10 -key-space=10000 -payload-size=256 -etcd-hosts=http://127.0.0.1
```

|flag|default value|description|
|-|-|-|
|clients|10|number of concurrent clients (int)|
|exec-time|10s|total execution time (duration)|
|cmd-timeout|1s|command timeout (duration)|
|max-thinking-time|100|maximum thinking time to wait between requests, in milliseconds (int)|
|key-space|10000|number of different keys (int64)|
|payload-size|256|payload size of values, in Bytes (int: 256, 512, 1024, 4096)|
|latency-filename||filename to write latency measurement, empty is disabled (string)|
|status-filename||filename to write response status measurement, empty is disabled (string)|
|etcd-hosts||list of etcd hostnames to send request, separated by `,` (string)|

## TODO List

* [x] Flags parse

* [x] Measurements pkg

* [x] Log level and outputs

* [ ] Utils and tests

* [x] Etcd integration

* [ ] Docs

* [ ] Generic DB initialization (v2)
