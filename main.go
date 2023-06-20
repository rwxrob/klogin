package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/kubectl-login/internal"
)

const reminder = `
Updated config file: %v
Cluster CA file:     %v
Dashboard:
  1) kubectl proxy
  2) http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/

`

func main() {

	home, _ := os.UserHomeDir()
	kubedir := filepath.Join(home, `.kube`)
	kubeconfig := filepath.Join(kubedir, `config`)

	os.Mkdir(kubedir, 0755)

	cluster := `prod` // default

	if len(os.Args) > 1 {
		cluster = os.Args[1]
	} else {
		current := strings.TrimSpace(internal.OutQuiet(`kubectl`, `config`, `current-context`))
		if current != "" {
			cluster = current
		}
	}

	ci, has := internal.Clusters[cluster]
	if !has {
		log.Fatalf(`unable to locate info for cluster: %v`, cluster)
	}

	fmt.Printf("Please enter login information for '%v' cluster:\n", cluster)

	user := internal.Prompt(`Username: `)
	pass := internal.PromptHidden(`Password: `)
	fmt.Println()
	fmt.Println()

	grant, err := internal.ReqOIDCPassAuth(
		user, pass, ci.OIDCIssuerURL, ci.ClientID, ci.ClientSecret,
	)
	if err != nil {
		log.Fatal(err)
	}

	cafile := filepath.Join(kubedir, cluster+`.crt`)
	if err := os.WriteFile(cafile, ci.CA, 0600); err != nil {
		log.Fatalf(`failed attempting to write cert to %v`, cafile)
	}

	internal.Exec(`kubectl`, `config`, `set-cluster`, ci.Name,
		`--server`, ci.APIServerURL,
		`--certificate-authority`, cafile,
		`--embed-certs`,
	)

	internal.Exec(
		`kubectl`, `config`, `set-context`, ci.Name,
		`--cluster`, ci.Name,
		`--user`, ci.Name,
		`--namespace`, user, // some may prefer `default`
	)

	internal.Exec(`kubectl`, `config`, `use-context`, ci.Name)

	token, _ := grant[`id_token`].(string)
	internal.Exec(`kubectl`, `config`, `set-credentials`, ci.Name, `--token`, token)

	fmt.Printf(reminder, kubeconfig, cafile)

}
