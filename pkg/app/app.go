package app

import (
	"runtime/debug"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ACME struct {
	Addr     string
	DirCache string
	Hosts    []string
	Email    string
	Staging  bool
}

type Flags struct {
	GRPCAddr        string
	GRPCTLSCertFile string
	GRPCTLSKeyFile  string
	ACME            ACME
	DebugAddr       string
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
		if d.Path == "github.com/Percona-Platform/platform" {
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

func Setup() *Flags {
	kingpin.Version(version())

	kingpin.HelpFlag.Short('h')

	var flags Flags

	kingpin.Flag("grpc.addr", "gRPC listen address").Default(":443").StringVar(&flags.GRPCAddr)
	kingpin.Flag("grpc.tls.cert-file", "gRPC listen address").StringVar(&flags.GRPCTLSCertFile)
	kingpin.Flag("grpc.tls.key-file", "gRPC listen address").StringVar(&flags.GRPCTLSKeyFile)

	kingpin.Flag("acme.addr", "ACME listen address").Default(":80").StringVar(&flags.ACME.Addr)
	kingpin.Flag("acme.dir-cache", "ACME directory cache").StringVar(&flags.ACME.DirCache)
	kingpin.Flag("acme.hosts", "ACME whitelisted hosts").StringsVar(&flags.ACME.Hosts)
	kingpin.Flag("acme.email", "ACME email").StringVar(&flags.ACME.Email)
	kingpin.Flag("acme.staging", "Use Let's Encrypt staging environment").BoolVar(&flags.ACME.Staging)

	kingpin.Flag("debug.addr", "Debug listen address").Default("127.0.0.1:8080").StringVar(&flags.DebugAddr)

	return &flags
}
