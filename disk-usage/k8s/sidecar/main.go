package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/haborhuang/go-tools/disk-usage/monitor"
	"github.com/haborhuang/go-tools/disk-usage/env"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "Listening port")
	flag.Parse()

	// Start monitor
	// Configuration can be set by env
	monitor.StartDUMonitorOrDie(env.ParseMetricsConfig(), env.ParseDisksPathsOrDie(""))

	http.Handle("/metrics", promhttp.Handler())
	// Start server
	log.Printf("Listening %s...\n", port)
	port = ":" + port
	log.Fatal(http.ListenAndServe(port, nil))
}

