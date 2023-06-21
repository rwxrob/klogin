package kubeconf_test

import (
	"fmt"

	"github.com/rwxrob/kubectl-login/internal/kubeconf"
)

func ExampleCurContext() {
	fmt.Println(kubeconf.CurContext())

	// Output:
	// prod
}

func ExampleContexts() {
	contexts := kubeconf.Contexts()
	fmt.Println(contexts[`prod`].Name)
	fmt.Println(contexts[`prod`].Cluster)
	fmt.Println(contexts[`prod`].User)
	fmt.Println(contexts[`prod`].Namespace)

	// Output:
	// prod
	// prod
	// prod
	// rwxrob
}
