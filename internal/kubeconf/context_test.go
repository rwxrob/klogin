package kubeconf_test

import (
	"fmt"

	"github.com/rwxrob/kubectl-login/internal/kubeconf"
)

func ExampleCurContextName() {
	fmt.Println(kubeconf.CurContextName())
	// Output:
	// prod
}

func ExampleCurContext() {
	fmt.Println(kubeconf.CurContext().Name)
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
