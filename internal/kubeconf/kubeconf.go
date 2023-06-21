package kubeconf

import (
	"os"
	"path/filepath"
)

// Dir returns the full path to the *main* $HOME/.kube directory.
func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, `.kube`)
}

// File returns the full path to the *main* $HOME/.kube/config file.
func File() string { return filepath.Join(Dir(), `config`) }
