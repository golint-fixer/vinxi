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
Duration  [total, attack, wait]   19.983667392s, 19.98306033s, 607.062µs
Latencies [mean, 50, 95, 99, max]   729.278µs, 692.796µs, 988.518µs, 9.955522ms, 9.955522ms
Bytes In  [total, mean]     12000, 12.00
Bytes Out [total, mean]     0, 0.00
Success   [ratio]       100.00%
Status Codes  [code:count]      200:1000  
Error Set:
```

## License

MIT
