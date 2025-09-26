package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		host     = flag.String("metrics.host", ":9000", "URL host for OVS datapath exporter")
		pathname = flag.String("metrics.pathname", "/metrics", "URL pathname exposing the collected metrics")
	)

	flag.Parse()

	registry := prometheus.NewRegistry()

	// Register both collectors
	dpCollector := newOvsDPCollector()
	vswitchdCollector := &ovsVswitchdCollector{}
	registry.MustRegister(dpCollector)
	registry.MustRegister(vswitchdCollector)

	fmt.Printf("Starting server listening: %s\n", *host)
	http.Handle(*pathname, promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorHandling: promhttp.HTTPErrorOnError}))
	http.ListenAndServe(*host, nil)

}
