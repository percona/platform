// Package ptls provides common TLS utilities for all SaaS components.
package ptls

import (
	"crypto/tls"

	"github.com/pkg/errors"
)

// GetConfig returns a new tls.Config instance configured according to Percona's security baseline.
func GetConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			// no SHA-1, ECDHE before plain RSA, GCM before CBC
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		},
	}
}

// GetConfig returns a new tls.Config with given certificate and key in PEM format.
func GetConfigWithCert(cert, key []byte) (*tls.Config, error) {
	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse TLS data")
	}

	tlsConfig := GetConfig()
	tlsConfig.Certificates = []tls.Certificate{pair}
	return tlsConfig, nil
}

// GetConfig returns a new tls.Config with given certificate and key files in PEM format.
func GetConfigWithCertFiles(certFile, keyFile string) (*tls.Config, error) {
	pair, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load TLS files")
	}

	tlsConfig := GetConfig()
	tlsConfig.Certificates = []tls.Certificate{pair}
	return tlsConfig, nil
}
