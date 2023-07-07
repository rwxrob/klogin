package run

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// Exec checks for existence of first argument as an executable on the
// system and then runs it with exec.Command.Run  exiting in a way that
// is supported across all architectures that Go supports. The stdin,
// stdout, and stderr are connected directly to that of the calling
// program. Sometimes this is insufficient and the UNIX-specific SysExec
// is preferred. See exec.Command.Run for more about distinguishing
// ExitErrors from others.
func Exec(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing name of executable")
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	cmd := exec.Command(path, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Prompt prints the given message reads the string by calling Read. The
// argument signature is identical as that passed to fmt.Printf().
func Prompt(form string, args ...any) string {
	fmt.Printf(form, args...)
	return Read()
}

// PromptHidden prints the given message if the terminal IsInteractive
// and reads the string by calling ReadHidden (which does not echo to
// the screen). The argument signature is identical and passed to to
// fmt.Printf().
func PromptHidden(form string, args ...any) string {
	fmt.Printf(form, args...)
	return ReadHidden()
}

// Read reads a single line of input and chomps the \r?\n. Also see
// ReadHidden.
func Read() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// ReadHidden disables the cursor and echoing to the screen and reads
// a single line of input. Leading and trailing whitespace are removed.
// Also see Read.
func ReadHidden() string {
	byt, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(byt))
}
