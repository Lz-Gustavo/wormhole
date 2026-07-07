# wormhole

Open loop benchmark tool, inspired by [YCSB Is Obsolete, We Need New Benchmarks](https://emptysqua.re/blog/ycsb-is-obsolete/).

## Usage

```
./bin/wormhole -clients=10 -exec-time=10s -cmd-timeout=300ms -max-thinking-time=10 -key-space=10000 -payload-size=256
```

|flag|default value|description|
|-|-|-|
|clients|10|number of concurrent clients (int)|
|exec-time|10s|total execution time (duration)|
|cmd-timeout|1s|command timeout (duration)|
|max-thinking-time|100|maximum thinking time to wait between requests, in milliseconds (int)|
|key-space|10000|number of different keys (int)|
|payload-size|256|payload size of values, in Bytes (int)|

## TODO List

* [x] Flags parse

* [ ] Measurements pkg

* [x] Log level and outputs

* [ ] Utils and tests

* [ ] Etcd integration

* [ ] Docs
