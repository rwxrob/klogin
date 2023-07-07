package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/rwxrob/klogin/internal/auth"
	"github.com/rwxrob/klogin/internal/clusters"
	"github.com/rwxrob/klogin/internal/run"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	cli "k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd/api"
)

var version = `v0.3.0`

//go:embed root.txt
var description string // easier to maintain in its own file

// set from StringVar flags
var userName, namespace string

var rootCmd = &cobra.Command{
	Use:               `klogin [USER@CLUSTER|CLUSTER]`,
	Short:             `Login to prod, dev, or inf clusters as specific user`,
	Long:              description,
	Version:           version,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: rootComplete,

	RunE: func(x *cobra.Command, args []string) error {

		var uname string  // user
		var clname string // dev, prod, inf
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
			/*
				uname, clname = auth.ParseTarget(args[0])
				if len(uname) == 0 {
					ctx, has = conf.Contexts[clname]
					if has {
						conf.CurrentContext = clname
						uname, clname = auth.ParseTarget(ctx.AuthInfo)
					}
				}
			*/
		default:
			return fmt.Errorf(`invalid arguments: %v`, args)
		}

		if len(userName) != 0 {
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

		if len(namespace) != 0 {
			ctx.Namespace = namespace
		}

		conf.Contexts[clname] = ctx

		// --------------  end of cleaning up CurrentContext -------------

		pass := run.PromptHidden(`Password: `)
		fmt.Println()

		// LoginROPC depends heavily on CurrentContext updated and pointing
		// to valid, supported clusters

		if err = auth.LoginROPC(&conf, pass); err != nil {
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

	// now add any contexts with AuthInfo that points to supported cluster

	for k, ctx := range conf.Contexts {
		_, clname := auth.ParseTarget(ctx.AuthInfo)
		if !slices.Contains(clusters.Names, clname) {
			continue
		}
		if strings.HasPrefix(k, in) {
			possible = append(possible, k)
		}
	}

	// now look for any AuthInfo matches
	/*
		for k, _ := range conf.AuthInfos {
			var name string
			f := strings.SplitN(k, `@`, 2)
			switch len(f) {
			case 1:
				name = f[0]
			case 2:
				name = f[1]
			}
			if len(name) != 0 && !slices.Contains(clusters.Names, name) {
				continue
			}
			if strings.HasPrefix(k, in) {
				possible = append(possible, k)
			}
		}
	*/

	return possible, 0

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.PersistentFlags().StringVar(&userName, "user", "", "Username to use for login")
	rootCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace to use for context")
	rootCmd.AddCommand(completionCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
