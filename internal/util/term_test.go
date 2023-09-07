package util_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/rwxrob/klogin/internal/util"
)

func ExampleGetCh() {

	var b byte
	var err error
	var r io.Reader

	r = strings.NewReader(`s4ж`)
	for {
		b, err = util.GetCh(r)
		fmt.Printf("byte: %x error: %v\n", b, err)
		if err == io.EOF {
			break
		}
	}
	fmt.Println(`done`)

	// Output:
	// byte: 73 error: <nil>
	// byte: 34 error: <nil>
	// byte: d0 error: <nil>
	// byte: b6 error: <nil>
	// byte: 0 error: EOF
	// done

}
