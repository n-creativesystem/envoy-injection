package helper

import "net/http"

var status = []byte("{\"status\":\"UP\"}")

var HealthCheck = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(status)
})
