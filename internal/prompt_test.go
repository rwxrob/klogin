package internal_test

import (
	"fmt"

	"github.com/rwxrob/kubectl-login/internal"
)

func ExampleIsInteractive() {
	fmt.Println(internal.IsInteractive)
	// Output:
	// false
}
