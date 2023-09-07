package util

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// GetCh reads and returns exactly one byte from the io.Reader. Returns
// any errors from [io.Reader.Read] or io.EOF if nothing left to read.
// Note that this reads one byte (character) at a time, not runes (which
// can be multiple bytes).
func GetCh(in io.Reader) (byte, error) {
	buf := make([]byte, 1)
	if n, err := in.Read(buf); n == 0 || err != nil {
		if err != nil {
			return 0, err
		}
		return 0, io.EOF
	}
	return buf[0], nil
}

// IsTerminal returns true if file descriptor is an interactive
// terminal.
func IsTerminal(fd uintptr) bool {
	return terminal.IsTerminal(int(fd))
}

// ReadBytesFromTerm reads a token (up to 4096 bytes) from the terminal
// standard input and prints the mark for every markn bytes to the terminal.
func ReadBytesFromTerm(mark string, markn, max int) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	state, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	defer func() {
		terminal.Restore(int(os.Stdin.Fd()), state)
	}()
	if err != nil {
		return nil, err
	}
	for i := 0; i < max; i++ {
		v, e := GetCh(os.Stdin)
		if v <= 0 || e != nil || v == 13 || v == 10 || v == 3 || v == 4 {
			break
		}
		err = buf.WriteByte(v)
		if i%markn == 0 {
			fmt.Print(mark)
		}
	}
	return buf.Bytes(), err
}
