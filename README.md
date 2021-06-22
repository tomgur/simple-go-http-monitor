# Simple HTTP monitor

### A Simple HTTP monitor written in Golang that exposes its findings to Prometheus

---
#### SETUP
1. Configure the monitor server by editing monitor.go(lines 44-46) 
* url - The URL that should be monitored
* interval - how often the monitor should run
* scrapePort - the port on which to expose the metrics for Prometheus. 

2. Run the monitor
* Issue the following command `go run monitor.go`
* Build the go executable and run it `go build monitor.go`
* Build a docker image containing the monitor app `docker build -t mymonitor:1; docker run --name monitor -d -p 9100:9100 mymonitor:1`
  
### Metrics
* Metrics will be exposed at the host IP on the port set in `scrapePort` at the `/metrics` endpoint (i.e. `1.2.3.4:9100/metrics`)
* The exposed metric is the round-trip duration in seconds. It is called `website-load-time` and has the following 3 quantiles `{0.5, 0.9, 0.99}`

### `from` Label
When run on an AWS EC2 instance, a label will be added to all metrics `from=<availability zone>`.

Whne run outside of AWS, the `from` label will be filled with the outgoing IP address of the server. 

