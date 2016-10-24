// Monitor.go
package web2tcp

import (
	_ "expvar"
	"net/http"
	_ "net/http/pprof"
)

type metricsHandler struct {
}

func StartMonitor() {
	go http.ListenAndServe(":9000", http.DefaultServeMux)
}
