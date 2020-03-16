// Package app provides common flags for all SaaS components.
package app

import (
	"crypto/tls"
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/percona-platform/platform/pkg/ptls"
)

type ACME struct {
	Addr     string
	DirCache string
	Hosts    []string
	Email    string
	Staging  bool
}

type Flags struct {
	GRPCAddr string

	GRPCTLSCert string
	GRPCTLSKey  string

	GRPCTLSCertFile string
	GRPCTLSKeyFile  string

	ACME ACME

	DebugAddr string
}

func version() string {
	l := zap.L().With(zap.String("component", "platform/app/version")).Sugar()
	l.Debug("Building version information.")

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(unknown)"
	}

	var platform *debug.Module
	for _, d := range info.Deps {
		if d.Path == "github.com/percona-platform/platform" {
			platform = d
			l.Debug(platform)
			if d.Replace != nil {
				platform = d.Replace
				l.Debug("\treplaced by ", platform)
			}
			break
		}
	}

	version := info.Main.Version
	if s := info.Main.Sum; s != "" {
		version += " (" + s + ")"
	}
	if platform != nil && platform.Version != "" {
		version += " / platform " + platform.Version
		if s := platform.Sum; s != "" {
			version += " (" + s + ")"
		}
	}

	return version
}

type SetupOpts struct {
	Name     string
	WithGRPC bool
	WithACME bool
}

func Setup(opts *SetupOpts) (*Flags, error) {
	if opts == nil {
		opts = new(SetupOpts)
	}

	if opts.Name == "" {
		return nil, errors.New("app.Setup: no Name")
	}

	kingpin.CommandLine.Name = opts.Name
	kingpin.CommandLine.DefaultEnvars()
	kingpin.Version(version())
	kingpin.HelpFlag.Short('h')

	var flags Flags

	if opts.WithGRPC {
		kingpin.Flag("grpc.addr", "gRPC listen address").Default(":443").StringVar(&flags.GRPCAddr)
		kingpin.Flag("grpc.tls.cert", "gRPC TLS certificate").StringVar(&flags.GRPCTLSCert)
		kingpin.Flag("grpc.tls.key", "gRPC TLS private key").StringVar(&flags.GRPCTLSKey)
		kingpin.Flag("grpc.tls.cert-file", "gRPC TLS certificate file").StringVar(&flags.GRPCTLSCertFile)
		kingpin.Flag("grpc.tls.key-file", "gRPC TLS private key file").StringVar(&flags.GRPCTLSKeyFile)
	}

	if opts.WithACME {
		kingpin.Flag("acme.addr", "ACME listen address").Default(":80").StringVar(&flags.ACME.Addr)
		kingpin.Flag("acme.dir-cache", "ACME directory cache").StringVar(&flags.ACME.DirCache)
		kingpin.Flag("acme.hosts", "ACME whitelisted hosts").StringsVar(&flags.ACME.Hosts)
		kingpin.Flag("acme.email", "ACME email").StringVar(&flags.ACME.Email)
		kingpin.Flag("acme.staging", "Use Let's Encrypt staging environment").BoolVar(&flags.ACME.Staging)
	}

	kingpin.Flag("debug.addr", "Debug listen address").Default(":8080").StringVar(&flags.DebugAddr)

	return &flags, nil
}

// TLSConfig returns TLS configuration and optional ACME handler from flags.
func (f *Flags) TLSConfig() (*tls.Config, http.Handler, error) {
	l := zap.L().With(zap.String("component", "platform/app/TLSConfig")).Sugar()

	switch {
	case f.GRPCTLSCert != "" && f.GRPCTLSKey != "":
		if f.GRPCTLSCertFile != "" || f.GRPCTLSKeyFile != "" {
			return nil, nil, errors.New("both GRPCTLSCert/GRPCTLSKey and GRPCTLSCertFile/GRPCTLSKeyFile are specified")
		}
		if f.ACME.DirCache != "" {
			return nil, nil, errors.New("both GRPCTLSCert/GRPCTLSKey and ACME are specified")
		}

		l.Info("Using given certificate and key for gRPC server.")
		c, err := ptls.GetConfigWithCert([]byte(f.GRPCTLSCert), []byte(f.GRPCTLSKey))
		return c, nil, err

	case f.GRPCTLSCertFile != "" && f.GRPCTLSKeyFile != "":
		if f.ACME.DirCache != "" {
			return nil, nil, errors.New("both GRPCTLSCertFile/GRPCTLSKeyFile and ACME are specified")
		}

		l.Info("Using given certificate and key files for gRPC server.")
		c, err := ptls.GetConfigWithCertFiles(f.GRPCTLSCertFile, f.GRPCTLSKeyFile)
		return c, nil, err

	case f.ACME.DirCache != "":
		return ptls.GetACME(&ptls.GetACMEOpts{
			DirCache: f.ACME.DirCache,
			Hosts:    f.ACME.Hosts,
			Email:    f.ACME.Email,
			Staging:  f.ACME.Staging,
		})

	default:
		return nil, nil, errors.New("no TLS configuration")
	}
}
