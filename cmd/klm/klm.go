package main

import (
	"fmt"
	"net/http"

	"github.com/leoluz/klog-metrics/pkg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println(".: Kube Log Metrics")
	pkg.LoadConfiguration()
	http.Handle("/metrics", promhttp.Handler())

	serverAddress := fmt.Sprintf(":%v", pkg.GetServerPort())
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}
