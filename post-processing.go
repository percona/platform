//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	platformRepo = "github.com/percona-platform/platform"
	saasRepo     = "github.com/percona-platform/saas"
	saasRoot     = "../saas"
)

var generatedImportRe = regexp.MustCompile(`(?mi)[\n]^.*github_com_mwitkow.*$`)

func saasFilePatch(content []byte) []byte {
	return bytes.Replace(content, []byte(platformRepo), []byte(saasRepo), -1)
}

// copyAndPatchFile copies a file src to dst, applying a specified patch function
func copyAndPatchFile(src, dst string, patchFunc func([]byte) []byte) error {
	b, err := ioutil.ReadFile(src) //nolint:gosec
	if err != nil {
		return err
	}

	b = patchFunc(b)

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

func removeDirs(root string, directories ...string) {
	for _, d := range directories {
		path := filepath.Join(root, d)
		log.Printf("Removing %s ...", path)
		if err := os.RemoveAll(path); err != nil {
			log.Fatal(err)
		}
	}
}

// makeProcessDirsFunc returns a function that applies patch function to included files and copies them to the root directory.
func makeProcessDirsFunc(root string, patchFunc func([]byte) []byte, includeFiles []string, excludeFiles []string) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var copy bool
		for _, s := range includeFiles {
			if strings.Contains(path, "internal") {
				// TODO: Improve internal packages handling
				panic("internal packages should not be copied to saas repo")
			}
			if strings.HasSuffix(path, s) {
				copy = true
			}
		}
		for _, s := range excludeFiles {
			if strings.HasSuffix(path, s) {
				copy = false
			}
		}

		if !copy {
			return nil
		}

		dst := filepath.Join(root, path)
		log.Printf(" %s -> %s", path, dst)
		return copyAndPatchFile(path, dst, patchFunc)
	}
}

func walk(processDirsFunc func(string, os.FileInfo, error) error, directories ...string) {
	log.Print("Copying and patching files:")
	for _, src := range directories {
		err := filepath.Walk(src, processDirsFunc)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func processSaas() {
	if _, err := os.Stat(saasRoot); err != nil {
		log.Fatal(err)
	}

	removeDirs(saasRoot, "api", "gen", "pkg")

	processDirsFunc := makeProcessDirsFunc(saasRoot, saasFilePatch, []string{".go", ".proto"}, []string{"_test.go", "_fuzz.go"})

	walk(processDirsFunc,
		"api/auth", "api/check", "api/telemetry",
		"gen/auth", "gen/check", "gen/telemetry",
		"pkg/check", "pkg/logger", "pkg/starlark",
		"pkg/alert", "pkg/common",
	)

	// install and tidy to check if we have anything
	_, err := runInDir(saasRoot, "go", "install", "-v", "./...")
	if err != nil {
		log.Fatal(err)
	}
	_, err = runInDir(saasRoot, "go", "mod", "tidy")
	if err != nil {
		log.Fatal(err)
	}

	// check dependencies
	b, err := runInDir(saasRoot, "go", "list", "-json", "./...")
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
			if strings.Contains(dep, platformRepo) {
				log.Fatalf("%s depends on platform module:\n%s", info.Dir, strings.Join(info.Deps, "\n"))
			}
		}
	}
}

func main() {
	const saasProject = "saas"

	availableProjects := []string{
		saasProject,
	}

	availableProjectsStr := strings.Join(availableProjects, " | ")

	project := flag.String("project", "", fmt.Sprintf("project to run post-processing for (%s)", availableProjectsStr))

	flag.Parse()

	if flag.NFlag() > 1 {
		flag.PrintDefaults()
		log.Fatal("Too many arguments, use only one")
	}

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		log.Fatal("You have to provide one argument")
	}

	switch *project {
	case saasProject:
		processSaas()
	default:
		flag.PrintDefaults()
		log.Fatal("Provide the target project name")
	}
}
