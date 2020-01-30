package app

import (
	"fmt"
	"runtime/debug"

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
		info, ok := debug.ReadBuildInfo()
		if ok {
			version = fmt.Sprintf("%s (%s)", info.Main.Version, info.Main.Sum)
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
