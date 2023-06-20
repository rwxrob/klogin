package internal

import (
	"encoding/json"
	"fmt"
	"log"
)

// RResp captures REST JSON responses as a simple key-value pair map
// allowing all values from that API response to be flexibly captured
// without needing to be marshaled into a specific struct. RResp
// implements the fmt.Stringer interface as long-form JSON for
// convenience when printing and storing. Note that REST APIs that
// return JSON lists are not supported.
type RResp map[string]any

// String implements the fmt.Stringer interface as JSON returning "null"
// (a valid JSON value) if anything goes wrong during the marshalling.
func (r RResp) String() string {
	byt, err := json.Marshal(r)
	if err != nil {
		return "null"
	}
	return string(byt)
}

func (r *RResp) Print() { fmt.Println(r) }
func (r *RResp) Log()   { log.Print(r) }
