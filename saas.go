// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	from = "github.com/percona-platform/platform"
	to   = "github.com/percona-platform/saas"
)

// copyAndPatchFile copies files src to dst,
// changing all occurrences of "github.com/percona-platform/platform" to "github.com/percona-platform/saas"
func copyAndPatchFile(src, dst string) error {
	b, err := ioutil.ReadFile(src) //nolint:gosec
	if err != nil {
		return err
	}

	b = bytes.Replace(b, []byte(from), []byte(to), -1)

	if err = os.MkdirAll(filepath.Dir(dst), 0o755); err != nil { //nolint:gosec
		return err
	}

	return ioutil.WriteFile(dst, b, 0o644)
}

// runInDir runs command name with args in dir and returns stdout.
func runInDir(dir, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...) //nolint:gosec
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	log.Print(strings.Join(cmd.Args, " "))
	return cmd.Output()
}

func main() {
	flag.Parse()

	const root = "../saas"
	if _, err := os.Stat(root); err != nil {
		log.Fatal(err)
	}

	// remove directories
	for _, d := range []string{"api", "gen", "pkg"} {
		path := filepath.Join(root, d)
		log.Printf("Removing %s ...", path)
		if err := os.RemoveAll(path); err != nil {
			log.Fatal(err)
		}
	}

	// copy and patch files
	for _, src := range []string{
		"api/auth", "api/check", "api/telemetry",
		"gen/auth", "gen/check", "gen/telemetry",
		"pkg/check", "pkg/logger", "pkg/starlark",
	} {
		err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			var copy bool
			for _, s := range []string{".go", ".proto"} {
				if strings.HasSuffix(path, s) {
					copy = true
				}
			}
			for _, s := range []string{"_test.go", "_fuzz.go"} {
				if strings.HasSuffix(path, s) {
					copy = false
				}
			}

			if !copy {
				return nil
			}

			dst := filepath.Join(root, path)
			log.Printf("%s -> %s", path, dst)
			return copyAndPatchFile(path, dst)
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	// install and tidy to check if we have anything
	_, err := runInDir(root, "go", "install", "-v", "./...")
	if err != nil {
		log.Fatal(err)
	}
	_, err = runInDir(root, "go", "mod", "tidy")
	if err != nil {
		log.Fatal(err)
	}

	// check dependencies
	b, err := runInDir(root, "go", "list", "-json", "./...")
	if err != nil {
		log.Fatal(err)
	}
	type packageInfo struct {
		Dir  string
		Deps []string
	}
	d := json.NewDecoder(bytes.NewReader(b))
	for {
		var info packageInfo
		err = d.Decode(&info)
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, dep := range info.Deps {
			if strings.Contains(dep, from) {
				log.Fatalf("%s depends on platform module:\n%s", info.Dir, strings.Join(info.Deps, "\n"))
			}
		}
	}
}
