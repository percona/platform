package ptls

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type GetACMEOpts struct {
	DirCache string
	Hosts    []string
	Email    string
	Staging  bool
}

// Package autocert provides automatic access to certificates from Let's Encrypt
// and any other ACME-based CA.
func GetACME(opts *GetACMEOpts) (*tls.Config, http.Handler, error) {
	if opts == nil {
		opts = new(GetACMEOpts)
	}

	if opts.DirCache == "" {
		return nil, nil, errors.New("ACME: no DirCache")
	}
	if len(opts.Hosts) == 0 {
		return nil, nil, errors.New("ACME: no Hosts")
	}
	if opts.Email == "" {
		return nil, nil, errors.New("ACME: no Email")
	}

	directoryURL := acme.LetsEncryptURL
	if opts.Staging {
		directoryURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}

	certManager := &autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		Cache:       autocert.DirCache(opts.DirCache),
		HostPolicy:  autocert.HostWhitelist(opts.Hosts...),
		RenewBefore: 30 * 24 * time.Hour,
		Client: &acme.Client{
			HTTPClient:   http.DefaultClient,
			DirectoryURL: directoryURL,
		},
		Email: opts.Email,
	}
	tlsConfig := GetConfig()
	tlsConfig.GetCertificate = certManager.GetCertificate
	tlsConfig.NextProtos = certManager.TLSConfig().NextProtos

	return tlsConfig, certManager.HTTPHandler(nil), nil
}
