package app

import (
	"runtime/debug"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ACME struct {
	Addr string
}

type Flags struct {
	GRPCAddr  string
	DebugAddr string
}

func Setup(version string) *Flags {
	if version == "" {
		l := zap.L().With(zap.String("component", "platform/app/Setup")).Sugar()
		l.Debug("Building version information.")

		info, ok := debug.ReadBuildInfo()
		if ok {
			l.Debug(info.Main)
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

			version = info.Main.Version
			if s := info.Main.Sum; s != "" {
				version += "(" + s + ")"
			}
			if platform != nil && platform.Version != "" {
				version += " / platform " + platform.Version
				if s := platform.Sum; s != "" {
					version += "(" + s + ")"
				}
			}
		}
	}

	kingpin.Version(version)

	kingpin.HelpFlag.Short('h')

	var flags Flags
	kingpin.Flag("grpc-addr", "gRPC listen address").Default(":443").StringVar(&flags.GRPCAddr)
	kingpin.Flag("debug-addr", "Debug listen address").Default(":8080").StringVar(&flags.DebugAddr)
	// kingpin.Flag("tls.addr", "Debug listen address").Default(":8080").StringVar(&flags.DebugAddr)

	return &flags
}
