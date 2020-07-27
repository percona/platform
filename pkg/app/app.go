// Package app provides common flags for all SaaS components.
package app

import (
	"runtime/debug"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Config is basic Percona Platform application configuration.
type Config struct {
	GRPCAddr  string
	HTTPAddr  string
	DebugAddr string
}

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(unknown)"
	}

	var platform *debug.Module
	for _, d := range info.Deps {
		if d.Path == "github.com/percona-platform/platform" {
			platform = d
			if d.Replace != nil {
				platform = d.Replace
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

// SetupOpts contains application requirements.
type SetupOpts struct {
	Name     string
	WithGRPC bool
	WithHTTP bool
}

// Setup returns application Config according to setup options.
func Setup(opts *SetupOpts) (*Config, error) {
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

	var config Config

	if opts.WithGRPC {
		kingpin.Flag("grpc.addr", "gRPC listen address").Default(":20201").StringVar(&config.GRPCAddr)
	}

	if opts.WithHTTP {
		kingpin.Flag("http.addr", "HTTP listen address").Default(":20202").StringVar(&config.HTTPAddr)
	}

	kingpin.Flag("debug.addr", "Debug listen address").Default(":20203").StringVar(&config.DebugAddr)

	return &config, nil
}
