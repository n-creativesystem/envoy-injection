package cert

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCA(t *testing.T) {
	certificate, err := GenerateRootCA()
	if !assert.NoError(t, err) {
		return
	}
	f, err := os.Create("ca.key")
	if !assert.NoError(t, err) {
		return
	}
	f.Write(certificate.Private)
	f.Close()

	f, err = os.Create("ca.crt")
	if !assert.NoError(t, err) {
		return
	}
	f.Write(certificate.Certificate)
	f.Close()
}

func TestServer(t *testing.T) {
	cert, _ := os.ReadFile("ca.crt")
	key, _ := os.ReadFile("ca.key")
	certificate, err := GenerateServer([]string{"service-name.namespace.svc"}, cert, key)
	if !assert.NoError(t, err) {
		return
	}
	f, err := os.Create("server.key")
	if !assert.NoError(t, err) {
		return
	}
	f.Write(certificate.Private)
	f.Close()

	f, err = os.Create("server.crt")
	if !assert.NoError(t, err) {
		return
	}
	f.Write(certificate.Certificate)
	f.Close()
}

func TestGenerate(t *testing.T) {
	root, server, err := GenerateCertificate([]string{"*.172.27.0.4.nip.io"})
	if !assert.NoError(t, err) {
		return
	}
	f, err := os.Create("ca.key")
	if !assert.NoError(t, err) {
		return
	}
	f.Write(root.Private)
	f.Close()

	f, err = os.Create("ca.crt")
	if !assert.NoError(t, err) {
		return
	}
	f.Write(root.Certificate)
	f.Close()

	f, err = os.Create("server.key")
	if !assert.NoError(t, err) {
		return
	}
	f.Write(server.Private)
	f.Close()

	f, err = os.Create("server.crt")
	if !assert.NoError(t, err) {
		return
	}
	f.Write(server.Certificate)
	f.Close()
}
