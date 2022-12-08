// Package fuzz provides fuzzing helpers.
package fuzz

import (
	"crypto/sha1" //nolint:gosec
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

//nolint:gochecknoglobals
var corpusM sync.Mutex

// AddToCorpus adds data to go-fuzz corpus.
func AddToCorpus(prefix string, b []byte) {
	corpusM.Lock()
	defer corpusM.Unlock()

	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("runtime.Caller failed")
	}
	dir := filepath.Join(filepath.Dir(file), "fuzzdata", "corpus")
	if err := os.MkdirAll(dir, 0o750); err != nil {
		panic(err)
	}

	// go-fuzz uses SHA1 for non-cryptographic hashing
	file = fmt.Sprintf("%040x", sha1.Sum(b)) //nolint:gosec
	if prefix != "" {
		prefix = strings.ReplaceAll(prefix, " ", "_")
		prefix = strings.ReplaceAll(prefix, "/", "_")
		file = prefix + "-" + file
	}

	path := filepath.Join(dir, file)
	if err := os.WriteFile(path, b, 0o640); err != nil { //nolint:gosec
		panic(err)
	}
}
