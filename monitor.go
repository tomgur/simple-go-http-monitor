package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// The monitoring loop
func monitorWebsite(loadTime prometheus.Summary, url string, interval int) {
	go func() {
		for {
			now := time.Now()
			// Sends an HTTP GET to the website
			get, _ := http.Get(url)
			elapsed := time.Since(now).Seconds()
			// Prints the status code and the elapsed time
			fmt.Printf("[INFO ] Status: [%d] Load time [%f]\n", get.StatusCode, elapsed)
			// Updates Prometheus with the elapsed time
			loadTime.Observe(elapsed)
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func GetVarOrDefault(varName string, defaultValue string) string {
	result := os.Getenv(varName)
	fmt.Printf(varName, result)
	if result == "" {
		result = defaultValue
		fmt.Printf("[INFO] Environment Variable [%s] not set - setting supplied default [%s]\n", varName, result)
	}
	return result
}

func main() {
	from := ""
	var scrapePort = flag.String("scrapePort", "9100", "The port on which Prometheus will scrape the metrics")
	var interval = flag.String("monitorInterval", "10", "How often the monitor should run in seconds")
	var url = flag.String("url", "https://hub.docker.com/r/tomgurdev/simple-go-http-monitor", "The URL that should be monitored")
	var subsystem = flag.String("subsystem", "website", " The subsystem (i.e. server/website)")
	var metricName = flag.String("metricName", "dockerHub_tomgurdev_simple-go-http-monitor_website_load_time", "The name of the published metric")
	var metricHelp = flag.String("metricHelp", "Docker Hub Simple Go HTTP Monitor website load time", "tTe help message for the metric")
	flag.Parse()
	fmt.Println("[INFO] Using the following parameters, scrapePort:", *scrapePort, " monitorInterval:", *interval, " url:", *url, " subsystem:", *subsystem, " metricName:", *metricName, " metricHelp:", *metricHelp)

	// 1 Sec timeout for the EC2 info site (if it's not there, the default timeout is 30 sec...)
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	// Get the Availability Zone from the EC2 info site
	response, err := client.Get("http://169.254.169.254/latest/meta-data/placement/availability-zone")
	// If the info site does not answer (not an EC2 instance, i.e. running on your laptop) set `from=UNKNOWN`
	if err != nil {
		fmt.Println("[WARN ] could not find AZ. Trying to find the local IP")
		localAddress := GetOutboundIP()
		fmt.Printf("[INFO ] Found local IP address. Setting `from=%s`\n", localAddress)
		from = localAddress.String()
	} else {
		//if we got an answer from EC2 info site, and we know the AZ, set `from=AZ`
		defer response.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		from = string(bodyBytes)
	}

	// create and register a new `Summary` with Prometheus
	summary := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:   "monitoring",
		Subsystem:   *subsystem,
		Name:        *metricName,
		Help:        *metricHelp,
		ConstLabels: prometheus.Labels{"from": from},
		Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		MaxAge:      0,
		AgeBuckets:  0,
		BufCap:      0,
	})
	prometheus.Register(summary)

	// Start the monitoring loop
	fmt.Printf("[INFO ] Starting to to monitor [%s], interval [%s]\n", *url, *interval)
	intervalStr, err := strconv.Atoi(*interval)
	if err != nil {
		log.Fatal(err)
	} else {
		monitorWebsite(summary, *url, intervalStr)
	}
	// Start the server, and set the /metrics endpoint to be served by the promhttp package
	fmt.Printf("[INFO ] Starting to serve metrics on port [%s]\n", *scrapePort)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+*scrapePort, nil)
}
