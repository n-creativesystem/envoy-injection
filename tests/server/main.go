package main

import "net/http"

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("{\"status\":\"OK\"}"))
	})

	_ = http.ListenAndServeTLS(":8443", "./server.crt", "./server.key", handler)
}
