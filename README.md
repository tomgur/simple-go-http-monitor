# Simple HTTP monitor

### A Simple HTTP monitor written in Golang that exposes its findings to Prometheus
 The monitor reads its target from env vars, and exposes a `Summary` with quantiles: `0.5`, `0.9`, `0.99`

---
#### SETUP
1. The monitor will not start unless the following environment variables are set: 
* `scrapePort` - The port on which Prometheus will scrape the metrics
* `subsystem` - The subsystem (i.e. server/website)
* `monitorUrl` - The URL that should be monitored
* `monitorInterval` - how often the monitor should run in seconds
* `metricName` - the name of the published metric
* `metricHelp` - the help message for the metric

2. Run the monitor (3 options)
* After exporting all env vars, issue the following command `go run monitor.go`
* After exporting all env vars, build the go executable and run it `go build monitor.go; ./monitor`
* Build a docker image containing the monitor app and run it `docker build -t monitor:1 .; docker run --name monitor -d -e "scrapePort=9100" -e "subsystem=website" -e "monitorUrl=https://google.com" -e "monitorInterval=10" -e "metricName=google_load_time" -e "metricHelp=Google website load time" mymonitor:1`
  
### Metrics
* Metrics will be exposed at the host IP on the port set in `scrapePort` at the `/metrics` endpoint (i.e. `1.2.3.4:9100/metrics`)
* The exposed metric is the round-trip duration in seconds in 3 quantiles: `0.5`, `0.9`, `0.99`

### `from` Label
When run on an AWS EC2 instance, a label will be added to all metrics `from=<availability zone>`.

Whne run outside of AWS, the from label will be filled with the outgoing IP address of the server. 

