// +build ignore

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// copyFile copies files src to dst,
func copyFile(src, dst string) error {
	b, err := ioutil.ReadFile(src) //nolint:gosec
	if err != nil {
		return err
	}

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

	const targetDir = "../saas-ui/packages/platform-ui/src/core"

	if _, err := os.Stat(targetDir); err != nil {
		log.Fatal(err)
	}

	// remove directories
	for _, d := range []string{"gen"} {
		path := filepath.Join(targetDir, d)
		log.Printf("Removing %s ...", path)
		if err := os.RemoveAll(path); err != nil {
			log.Fatal(err)
		}
	}

	// copy and patch files
	for _, src := range []string{
		"gen/web/auth",
	} {
		err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			var copy bool
			for _, s := range []string{".js", ".ts"} {
				if strings.HasSuffix(path, s) {
					copy = true
				}
			}

			if !copy {
				return nil
			}

			dst := filepath.Join(targetDir, path)
			log.Printf("%s -> %s", path, dst)
			return copyFile(path, dst)
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
