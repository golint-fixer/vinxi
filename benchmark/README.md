# vinxi benchmark 

Benchmark test suite for vinxi.

## Installation

Install [vegeta](https://github.com/tsenart/vegeta):
```bash
go get github.com/tsenart/vegeta
```

## Usage

Run the benchmark tests:
```bash
bash benchmark.sh
```

## Results

Environment: OSX 10.11 2.7GHz 16GB

```
Running 'simple' benchmark with concurrency 50
Requests  [total]       1000
Duration  [total, attack, wait]   19.98030457s, 19.979508728s, 795.842µs
Latencies [mean, 50, 95, 99, max]   997.76µs, 706.233µs, 4.474003ms, 8.12996ms, 8.12996ms
Bytes In  [total, mean]     11076, 11.08
Bytes Out [total, mean]     0, 0.00
Success   [ratio]       92.30%
Status Codes  [code:count]      0:77  200:923  
Error Set:
Get http://localhost:8080: dial tcp 127.0.0.1:8080: connection refused
```

## License

MIT
