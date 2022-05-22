package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	buf, _ := os.ReadFile("./server.crt")
	cp := x509.NewCertPool()
	cp.AppendCertsFromPEM(buf)
	tr := http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName: "test.172.27.0.4.nip.io",
			RootCAs:    cp,
		},
	}
	req, _ := http.NewRequest(http.MethodGet, "https://test.172.27.0.4.nip.io:8443", nil)
	client := http.Client{
		Transport: &tr,
	}
	res, err := client.Do(req)
	if !assert.NoError(t, err) {
		return
	}
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()
	buf, _ = io.ReadAll(res.Body)
	println(string(buf))
}
