package kubeconf

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/rwxrob/kubectl-login/internal/run"
)

// Context represents a current Kubernetes context as output from the
// kubectl config get-contexts output. Note that the format of the
// configuration file is different in that the name itself is
// (unfortunately) not under the context values. For this reason,
// marshaling and unmarshaling are not recommended since the resulting
// string is incompatible with the content of a kube config file.
type Context struct {
	Name      string
	Cluster   string
	User      string // referred to as AuthInfo elsewhere
	Namespace string
}

// CurContextName returns only the name of the current context. See
// CurContext if the entire Context struct is wanted instead.
func CurContextName() string {
	return strings.TrimSpace(run.OutQuiet(`kubectl`, `config`, `current-context`))
}

// CurContext returns the Context object from Contexts for the
// CurContextName or a nil if not found.
func CurContext() *Context {
	ctx, has := Contexts()[CurContextName()]
	if !has {
		return nil
	}
	return &ctx
}

// Contexts safely returns all the current user contexts (normally
// returned from kubectl config get-contexts) as a map by parsing the
// kubectl view config -o jsonpath='jsonpath={.contexts}' output.
func Contexts() map[string]Context {
	contexts := map[string]Context{}

	out := run.OutQuiet(`kubectl`, `config`, `view`, `-o`, `jsonpath={.contexts}`)

	// slice of "named context" structs
	holder := []struct {
		Name    string
		Context struct {
			Cluster   string
			User      string
			Namespace string
		}
	}{}

	err := json.Unmarshal([]byte(out), &holder)
	if err != nil {
		log.Print(err)
		return contexts
	}

	for _, c := range holder {
		if !(len(c.Name) > 0) {
			continue
		}
		contexts[c.Name] = Context{
			Name:      c.Name,
			Cluster:   c.Context.Cluster,
			User:      c.Context.User,
			Namespace: c.Context.Namespace,
		}
	}

	return contexts
}
