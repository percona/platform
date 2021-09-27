// Package app provides common flags for all SaaS components.
package app

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/percona-platform/platform/pkg/logger"
)

// Config is basic Percona Platform application configuration.
type Config struct {
	GRPCAddr    string // gRPC Server address
	GRPCWebAddr string // gRPC-Web Server address
	HTTPAddr    string // HTTP Server address
	DebugAddr   string // debug Server address

	LoggerOpts logger.SetupGlobalOpts
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
	Name        string
	WithGRPC    bool
	WithGRPCWeb bool
	WithHTTP    bool
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

	if opts.WithGRPCWeb {
		kingpin.Flag("grpc-web.addr", "gRPC-Web listen address").Default(":20204").StringVar(&config.GRPCWebAddr)
	}

	if opts.WithHTTP {
		kingpin.Flag("http.addr", "HTTP listen address").Default(":20202").StringVar(&config.HTTPAddr)
	}

	kingpin.Flag("debug.addr", "Debug listen address").Default(":20203").StringVar(&config.DebugAddr)

	// use global environment variables PLATFORM_LOG_XXX for defaults values,
	// but allow to set flags via normal APP_PLATFORM_LOG_XXX environment variables
	b, _ := strconv.ParseBool(os.Getenv("PLATFORM_LOG_DEBUG"))
	logDebugDefault := strconv.FormatBool(b)
	kingpin.Flag("log.debug", "Enable debug level logging").
		Default(logDebugDefault).BoolVar(&config.LoggerOpts.LogDebug)
	b, _ = strconv.ParseBool(os.Getenv("PLATFORM_LOG_DEVMODE"))
	logDevMode := strconv.FormatBool(b)
	kingpin.Flag("log.devmode", "Enable development mode loging: text instead of JSON, DPanic panics instead of logging errors").
		Default(logDevMode).BoolVar(&config.LoggerOpts.LogDevMode)
	config.LoggerOpts.LogName = fmt.Sprintf("platform.%s", opts.Name)

	kingpin.CommandLine.Action(func(*kingpin.ParseContext) error {
		logger.FlagsParsed = true
		return nil
	})

	return &config, nil
}
