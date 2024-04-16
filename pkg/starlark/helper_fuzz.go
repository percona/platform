package starlark

import (
	"encoding/json"

	"github.com/percona/platform/pkg/fuzz"
)

// Those helpers are shared by tests and fuzz code,
// and can't be in either because the other user will not find them.
// _fuzz.go suffix is required for `make saas`.

type fuzzData struct {
	Script string                   `json:"s,omitempty"`
	Input  []map[string]interface{} `json:"i,omitempty"`
}

// addToFuzzCorpus adds data to go-fuzz corpus.
func addToFuzzCorpus(name, script string, input []map[string]interface{}) {
	fd := &fuzzData{
		Script: script,
		Input:  input,
	}
	b, err := json.Marshal(fd)
	if err != nil {
		panic(err)
	}

	fuzz.AddToCorpus(name, b)
}
