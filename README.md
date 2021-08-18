# Simple HTTP monitor

### A Simple HTTP monitor written in Golang that exposes its findings to Prometheus
The monitor reads its target from env vars, and exposes a `Summary` with quantiles: `0.5`, `0.9`, `0.99`

---
#### SETUP
1. The following env vars can be set for configuring the monitor

_Running the monitor without setting variables will Monitor This project's docker hub page once every 10 seconds and expose the metrics on port 9100_
* `scrapePort` - The port on which Prometheus will scrape the metrics (defautl 9100)
* `subsystem` - The subsystem (i.e. server/website/store - defualt website)
* `monitorUrl` - The URL that should be monitored (default https://hub.docker.com/repository/docker/tomgurdev/simple-go-http-monitor)
* `monitorInterval` - how often the monitor should run in seconds (defautlt 10)
* `componentName` - the name of the monitored service (default simple-http-monitor-docker-hub)

2. Run the monitor (3 options)
* After exporting all env vars, issue the following command `go run monitor.go`
* After exporting all env vars, build the go executable and run it `go build monitor.go; ./monitor`
* Build a docker image containing the monitor app and run it `docker build -t monitor:1 .; docker run --name monitor --restart always -d -p 9100:9100 -e "scrapePort=9100" -e "subsystem=website" -e "monitorUrl=https://hub.docker.com/repository/docker/tomgurdev/simple-go-http-monitor" -e "monitorInterval=10" -e "componentName=simple_http_monitor_docker_hub" mymonitor:1`
  
### Metrics
* Metrics will be exposed at the host IP on the port set in `scrapePort` at the `/metrics` endpoint (i.e. `1.2.3.4:9100/metrics`)
* The exposed metrics are:
    - The round-trip duration in seconds in 3 quantiles: `0.5`, `0.9`, `0.99`
    - The HTTP response status

### `from` Label
When run on an AWS EC2 instance, a label will be added to all metrics `from=<availability zone>`.

Whne run outside of AWS, the from label will be filled with the outgoing IP address of the server. 

