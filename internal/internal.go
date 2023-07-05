package internal

import (
	"encoding/json"
	"fmt"
	"log"
)

// AsJSON outputs whatever is passed as JSON (null if <nil>). Errors during
// unmashalling are logged and null returned.
func AsJSON(a any) string {
	if a == nil {
		return `null`
	}
	byt, err := json.Marshal(a)
	if err != nil {
		log.Print(err)
		return `null`
	}
	return string(byt)
}

// PrintAsJSON calls AsJSON and prints to standard output with
// fmt.Println.
func PrintAsJSON(a any) { fmt.Println(AsJSON(a)) }
