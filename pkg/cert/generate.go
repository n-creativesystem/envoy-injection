package cert

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

var (
	SubjectIssuer = "k8s"
)

type After struct {
	year  int
	month int
	day   int
}

func NewAfter(year, month, day int) After {
	return After{
		year:  year,
		month: month,
		day:   day,
	}
}

func (a After) AfterDate(before time.Time) time.Time {
	return before.AddDate(a.year, a.month, a.day)
}

type CryptographyType int

const (
	RSA CryptographyType = iota
	ECDSA
)

type CertificateConfig struct {
	CryptographyType CryptographyType
	Bits             int
	After            After
}

type Certificate struct {
	Private     []byte
	Certificate []byte
}

func public(priv interface{}) interface{} {
	switch v := priv.(type) {
	case *rsa.PrivateKey:
		return &v.PublicKey
	case *ecdsa.PrivateKey:
		return v.PublicKey
	case ed25519.PrivateKey:
		return v.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func generateCertificate(hosts []string, isCA bool, rootCert *x509.Certificate, rootKey interface{}) (*Certificate, error) {
	var priv interface{}
	bits := 4096
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	var isRsa bool
	keyUsage := x509.KeyUsageDigitalSignature
	if _, isRsa = priv.(*rsa.PrivateKey); isRsa {
		keyUsage |= x509.KeyUsageKeyEncipherment
	}

	notBefore := time.Now()
	notAfter := notBefore.AddDate(10, 0, 0)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}
	var derBytes []byte
	if isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
		template.Subject.CommonName = "Root CA"
		template.Subject.Organization = append(template.Subject.Organization, SubjectIssuer)
		derBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, public(priv), priv)
	} else {
		derBytes, err = x509.CreateCertificate(rand.Reader, &template, rootCert, public(priv), rootKey)
	}
	if err != nil {
		return nil, err
	}
	certOut := &bytes.Buffer{}
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return nil, err
	}

	keyOut := &bytes.Buffer{}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	var keyType string = "PRIVATE KEY"
	err = pem.Encode(keyOut, &pem.Block{Type: keyType, Bytes: privBytes})
	if err != nil {
		return nil, err
	}
	return &Certificate{
		Private:     keyOut.Bytes(),
		Certificate: certOut.Bytes(),
	}, nil
}

func GenerateRootCA() (*Certificate, error) {
	return generateCertificate(nil, true, nil, nil)
}

func byteToCertificate(value []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(value)
	return x509.ParseCertificate(block.Bytes)
}

func byteToPrivateKey(value []byte) (interface{}, error) {
	block, _ := pem.Decode(value)
	return x509.ParsePKCS8PrivateKey(block.Bytes)
}

func GenerateServer(hosts []string, rootCert, rootKey []byte) (*Certificate, error) {
	cert, err := byteToCertificate(rootCert)
	if err != nil {
		return nil, err
	}
	key, err := byteToPrivateKey(rootKey)
	if err != nil {
		return nil, err
	}
	return generateCertificate(hosts, false, cert, key)
}

func GenerateCertificate(hosts []string) (*Certificate, *Certificate, error) {
	root, err := GenerateRootCA()
	if err != nil {
		return nil, nil, err
	}
	server, err := GenerateServer(hosts, root.Certificate, root.Private)
	if err != nil {
		return nil, nil, err
	}
	return root, server, nil
}
