package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	d "runtime/debug"
	"strconv"
	"time"

	"gopkg.in/vinxi/vinxi.v0"
)

var (
	aAddr         = flag.String("a", "", "bind address")
	aPort         = flag.Int("p", 8080, "Port to listen")
	aVers         = flag.Bool("v", false, "Show version")
	aVersl        = flag.Bool("version", false, "Show version")
	aHelp         = flag.Bool("h", false, "Show help")
	aHelpl        = flag.Bool("help", false, "Show help")
	aConfig       = flag.String("c", "", "Config file path")
	aConfigl      = flag.String("config", "", "Config file path")
	aForward      = flag.String("f", "", "Target server URL to forward traffic by default")
	aReadTimeout  = flag.Int("http-read-timeout", 60, "HTTP read timeout in seconds")
	aWriteTimeout = flag.Int("http-write-timeout", 60, "HTTP write timeout in seconds")
	aMRelease     = flag.Int("mrelease", 30, "OS memory release inverval in seconds")
	aCpus         = flag.Int("cpus", runtime.GOMAXPROCS(-1), "Number of cpu cores to use")
)

const usage = `vinxictl %s

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
                            (default for current machine is %d cores)
`

// TODO: support config file
func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, Version, runtime.NumCPU()))
	}
	flag.Parse()

	if *aHelp || *aHelpl {
		showUsage()
	}
	if *aVers || *aVersl {
		showVersion()
	}

	// Only required in Go < 1.5
	runtime.GOMAXPROCS(*aCpus)

	port := getPort(*aPort)
	opts := vinxi.ServerOptions{
		Port: port,
		Addr: *aAddr,
	}

	// Create a memory release goroutine
	if *aMRelease > 0 {
		memoryRelease(*aMRelease)
	}

	v := vinxi.New()

	// Define target server to forward incoming traffic
	if *aForward != "" {
		v.Forward(*aForward)
	}

	// Start HTTP server
	fmt.Printf("Server listening on port %d\n", port)
	v.ListenAndServe(opts)
}

func getPort(port int) int {
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		newPort, _ := strconv.Atoi(portEnv)
		if newPort > 0 {
			port = newPort
		}
	}
	return port
}

func showUsage() {
	flag.Usage()
	os.Exit(1)
}

func showVersion() {
	fmt.Println(vinxi.Version)
	os.Exit(1)
}

func memoryRelease(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for range ticker.C {
			// debug("FreeOSMemory()")
			d.FreeOSMemory()
		}
	}()
}

func exitWithError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args)
	os.Exit(1)
}
