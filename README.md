# DoT-Proxy

dot-proxy is a simple DNS to DNS-over-TLS proxy server that listens for TCP and UDP queries, then forwards them to any DoT server. It works best with CloudFlare DoT servers and has been tested using both `1.1.1.1:853` and `1.0.0.1:853`.

dot-proxy listens on non-privileged port 8053 for both UDP and TCP by default. 


## Preface

This project utilizes the `make` command to ensure an improved developer experience for common operations. Here is the list of supported make targets. 

```
$ make help

help                            List of Makefile targets.
init                            Initializes the project.
get-deps                        Downloads dependencies.
run                             Runs main.go.
build                           Builds the Go binary for Linux, Darwin, and Windows platforms.
clean                           Removes the built binary and any other generated files. 
format                          Fixes formatting for all go files.  
formatcheck                     Checks formatting of all go files.
test                            Runs test cases.
container                       Builds a Docker container with the latest tag.
clean-container                 Removes the running container.
run-container                   Runs the container with the latest tag and default options.
log-container                   Checks the logs of the running container.
clean-all                       Cleans both binary and container files.
test-dig                        Runs test queries using dig.

```

### Usage
dot-proxy supports the following command line arguments:

```
Usage: dot-proxy [options]

  -cloudflare-dot-addr string
        Cloudflare DoT address. (default "1.1.1.1:853")
  -debug
        Enables Debug mode.
  -host-tcp string
        Host to listen on for TCP. (default "0.0.0.0")
  -host-udp string
        Host to listen on for UDP. (default "0.0.0.0")
  -port-tcp int
        Port to listen on for TCP. (default 8053)
  -port-up int
        Port to listen on for UDP. (default 8053)
  -tcp
        Enables TCP listening mode. (default true)
  -udp
        Enables UDP listening mode. (default true)
```

By default, both TCP and UDP are enabled, and it listens on non-privileged ports `0.0.0.0:8053/udp` and `0.0.0.0:8053/tcp`. You can change the behavior using the above arguments.

To run and test the program, you can use the following methods:

### Run main.go
Simply run main.go with the default arguments.

```shell
make run
```

### Run binary
The following command will automatically set the GOOS and build the binary. 

```shell
make build
./dot-proxy
```

### Run container
The following command will build a container image using Docker with the latest tag, then run it with the default options:

```shell
make run-container
```

## Test
Run the following command to quickly test dot-proxy for both UDP and TCP

```shell
make test-dig
```

##  TODO
1. Add more unit tests and increase test coverage.
2. Improve exception handling including time-outs.
3. Add CI/CD related resources to build a pipeline.
4. We can get an array of DoT services and write a simple failover system.
5. We can create a caching mechanism to make the service faster.

