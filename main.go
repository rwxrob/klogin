package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rwxrob/kubectl-login/internal/auth"
	"github.com/rwxrob/kubectl-login/internal/clusters"
	"github.com/rwxrob/kubectl-login/internal/kubeconf"
	"github.com/rwxrob/kubectl-login/internal/run"
)

func main() {

	dir := kubeconf.Dir()
	os.Mkdir(dir, 0755)
	contexts := kubeconf.Contexts()
	current := kubeconf.CurContextName()

	if len(os.Args) > 1 {
		current = os.Args[1]
	}

	if current == "" {
		current = `prod`
	}

	context, hascontext := contexts[current]
	cluster, hascluster := clusters.Map[current]

	if !hascontext && !hascluster {
		log.Fatalf(`unable to locate info for context/cluster: %v`, current)
	}

	if !hascontext && hascluster {
		context = kubeconf.Context{
			Name:      current,
			Cluster:   current,
			User:      current,
			Namespace: run.Prompt(`Username: `),
		}
	}

	if hascontext && !hascluster {
		log.Fatalf(`context (%v) unsupported by this login plugin`, current)
	}

	pass := run.PromptHidden(`Password: `)
	fmt.Println()

	grant, err := auth.ReqOIDCPass(
		context.Namespace, pass, cluster.OIDCIssuerURL, cluster.ClientID, cluster.ClientSecret,
	)
	if err != nil {
		log.Fatal(err)
	}

	cafile, err := os.CreateTemp(``, `kubectl-login`)
	cert := cafile.Name()
	defer os.Remove(cert)
	if _, err := cafile.Write(cluster.CA); err != nil {
		log.Fatalf(`failed attempting to write cert to %v`, cafile)
	}

	run.Exec(`kubectl`, `config`, `set-cluster`, cluster.Name,
		`--server`, cluster.APIServerURL,
		`--certificate-authority`, cert,
		`--embed-certs`,
	)

	run.Exec(
		`kubectl`, `config`, `set-context`, context.Name,
		`--cluster`, cluster.Name,
		`--user`, context.User,
		`--namespace`, context.Namespace,
	)

	run.Exec(`kubectl`, `config`, `use-context`, context.Name)

	token, _ := grant[`id_token`].(string)
	run.Exec(`kubectl`, `config`, `set-credentials`, context.User, `--token`, token)

}
