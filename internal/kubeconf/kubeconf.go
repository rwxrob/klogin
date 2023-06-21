package kubeconf

import (
	"os"
	"path/filepath"
)

// Full path to the *main* $HOME/.kube/config file.
var File string

// Full path to the *main* $HOME/.kube directory.
var Dir string

func init() {
	home, _ := os.UserHomeDir()
	Dir = filepath.Join(home, `.kube`)
	File = filepath.Join(Dir, `config`)
}
