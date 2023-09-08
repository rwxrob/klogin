package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/rwxrob/klogin/internal/auth"
	"github.com/rwxrob/klogin/internal/clusters"
	"github.com/rwxrob/klogin/internal/run"
	"github.com/rwxrob/klogin/internal/util"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	cli "k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd/api"
)

var version = `v0.3.1`

//go:embed root.txt
var description string

// set from Cobra StringVar flags
var userName, namespace string

var rootCmd = &cobra.Command{
	Use:               `klogin [CLUSTER|CONTEXT]`,
	Short:             `Login to supported clusters`,
	Long:              description,
	Version:           version,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: rootComplete,

	RunE: func(x *cobra.Command, args []string) error {

		var uname, clname string // uname@clname
		var has bool
		var ctx *api.Context

		// ---  grab configuration and update current context as needed ---

		// fetches *all* KUBECONFIG files (see api.Config)
		conf, err := cli.NewConfigFlags(true).ToRawKubeConfigLoader().RawConfig()
		if err != nil {
			return err
		}

		// handle arguments or infer

		switch len(args) {

		case 0:
			ctx, has = conf.Contexts[conf.CurrentContext]
			if has {
				uname, clname = auth.ParseTarget(ctx.AuthInfo)
			}

		case 1:
			ctx, has = conf.Contexts[args[0]]
			if has {
				uname, clname = auth.ParseTarget(ctx.AuthInfo)
				if !slices.Contains(clusters.Names, clname) {
					return fmt.Errorf(`unsupported cluster: %v (not one of %v)`, clname, clusters.Names)
				}
				conf.CurrentContext = args[0]
			}

		default:
			return fmt.Errorf(`invalid arguments: %v`, args)
		}

		if len(userName) != 0 { // from --user
			uname = userName
		}

		if len(clname) == 0 {
			clname = clusters.Default
		}

		if !slices.Contains(clusters.Names, clname) {
			return fmt.Errorf(`unsupported cluster: %v (not one of %v)`, clname, clusters.Names)
		}

		if len(uname) == 0 {
			uname = run.Prompt(`Username: `)
		}

		if ctx == nil {
			ctx = api.NewContext()
			conf.CurrentContext = clname
		}

		ctx.AuthInfo = uname + `@` + clname
		ctx.Cluster = clname

		if len(ctx.Namespace) == 0 {
			ctx.Namespace = uname
		}

		if len(namespace) != 0 { // from --namespace
			ctx.Namespace = namespace
		}

		conf.Contexts[clname] = ctx

		// --------------  end of cleaning up CurrentContext -------------

		fmt.Printf(
			"Please obtain a token from the following URL and paste it here:\n%v\n",
			clusters.Map[clname].LoginURL,
		)

		byt, err := util.ReadBytesFromTerm(`*`, 10, 4096)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println()
		token := strings.TrimSpace(string(byt))

		// LoginAuth depends heavily on CurrentContext updated and pointing
		// to valid, supported clusters in memory

		if err = auth.LoginAuth(&conf, token); err != nil {
			return err
		}

		return run.Exec(`kubectl`, `config`, `get-contexts`)
	},
}

// rootComplete adds any possible unqualified cluster/context name first
// to the list of possible completions, then it looks for anything in
// the user/credentials/AuthInfo context for a match as well.
func rootComplete(cmd *cobra.Command, args []string, in string) ([]string, cobra.ShellCompDirective) {
	possible := []string{}

	// first add the unqualified cluster names, if any

	for _, name := range clusters.Names {
		if strings.HasPrefix(name, in) {
			possible = append(possible, name)
		}
	}

	conf, err := cli.NewConfigFlags(true).ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return possible, 0
	}

	// add any contexts with AuthInfo that points to supported clusters

	for k, ctx := range conf.Contexts {
		_, clname := auth.ParseTarget(ctx.AuthInfo)
		if !slices.Contains(clusters.Names, clname) {
			continue
		}
		if strings.HasPrefix(k, in) {
			possible = append(possible, k)
		}
	}

	return possible, 0

}

// Execute is the main entry point from the command called by main.main().
func Execute() {
	rootCmd.PersistentFlags().StringVar(&userName, "user", "", "Username to use for login")
	rootCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace to use for context")
	rootCmd.AddCommand(completionCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
