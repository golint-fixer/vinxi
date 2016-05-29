vinxi provides a full-feateured command-line interface with a declarative configuration file capable of detailed configuration of one or multiple proxy servers.

**Note**: still a work in progress.

#### Installation

```bash
go get -u gopkg.in/vinxi/vinxictl.v0
```

#### Usage

```bash
vinxictl 0.1.0

Usage:
  vinxictl -p 80
  vinxictl -p 80 -c config.toml

Options:
  -a <addr>                 bind address [default: *]
  -p <port>                 bind port [default: 8080]
  -h, -help                 output help
  -v, -version              output version
  -c, -config               Config file path
  -f                        Target server URL to forward traffic by default
  -mrelease <num>           OS memory release inverval in seconds [default: 30]
  -cpus <num>               Number of used cpu cores.
                            (default for current machine is 8 cores)
```
