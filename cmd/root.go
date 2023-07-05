package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rwxrob/klogin/internal/auth"
	"github.com/rwxrob/klogin/internal/clusters"
	"github.com/rwxrob/klogin/internal/run"
	"github.com/spf13/cobra"
	cli "k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func parseUser(u string) (uname, cluster string) {
	it := strings.SplitN(u, `@`, 2)
	switch len(it) {
	case 2:
		cluster = it[1]
		uname = it[0]
	case 1:
		cluster = it[0]
	}
	return
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "klogin",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		var uname, clname, pass string
		var cl clusters.Cluster
		var has bool

		if len(args) == 1 {
			uname, clname = parseUser(args[0])
		}

		// fetch *all* KUBECONFIG files (see api.Config)
		conf, err := cli.NewConfigFlags(true).ToRawKubeConfigLoader().RawConfig()
		if err != nil {
			return err
		}

		curctx := conf.CurrentContext
		ctx, has := conf.Contexts[curctx]
		if has && len(args) == 0 && curctx != "" {
			user := ctx.AuthInfo
			uname, clname = parseUser(user)
		}

		if clname == "" {
			clname = `prod` // default
		}

		if uname == "" {
			uname = run.Prompt(`Username: `)
		}

		pass = run.PromptHidden(`Password: `)
		fmt.Println()

		cl, has = clusters.Map[clname]
		if !has {
			return fmt.Errorf(`unsupported cluster name: %v`, clname)
		}

		grant, err := auth.ReqOIDCPass(
			uname, pass, cl.OIDCIssuerURL, cl.ClientID, cl.ClientSecret,
		)
		if err != nil {
			log.Fatal(err)
		}
		token, isstring := grant[`id_token`].(string)
		if !isstring {
			log.Fatal(`id_token not found or is not a string`)
		}

		// reused by ModifyConfig calls
		o := clientcmd.NewDefaultPathOptions()

		// update/add the cluster entry
		cluster := api.NewCluster()
		cluster.Server = cl.APIServerURL
		cluster.CertificateAuthorityData = cl.CA
		conf.Clusters[cl.Name] = cluster
		if err := clientcmd.ModifyConfig(o, conf, true); err != nil {
			return err
		}

		// update/add the user@cluster entry
		authinfo := api.NewAuthInfo()
		authinfo.Token = token
		delete(conf.AuthInfos, cl.Name) // cleanup old, unqualified entries
		authinfoid := strings.Join([]string{uname, cl.Name}, `@`)
		conf.AuthInfos[authinfoid] = authinfo
		if err := clientcmd.ModifyConfig(o, conf, true); err != nil {
			return err
		}

		// update/add context
		ctx = api.NewContext()
		ctx.Cluster = cl.Name
		ctx.AuthInfo = authinfoid
		ctx.Namespace = uname
		conf.Contexts[authinfoid] = ctx
		if err := clientcmd.ModifyConfig(o, conf, true); err != nil {
			return err
		}

		// set current context
		conf.CurrentContext = authinfoid
		if err := clientcmd.ModifyConfig(o, conf, true); err != nil {
			return err
		}

		run.Exec(`kubectl`, `config`, `get-contexts`)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	rootCmd.AddCommand(completionCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
