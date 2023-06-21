package kubeconf

import (
	"os"
	"path/filepath"
)

// Dir returns the full expected path to the *main* $HOME/.kube
// directory even if it does not yet exist.
func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, `.kube`)
}

// File returns the full expected path to the *main* $HOME/.kube/config
// file even if it does not yet exist.
func File() string { return filepath.Join(Dir(), `config`) }
