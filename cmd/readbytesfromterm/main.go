package main

import (
	"fmt"

	"github.com/rwxrob/klogin/internal/util"
)

func main() {
	it, err := util.ReadBytesFromTerm(`*`, 10, 4096)
	fmt.Println(it, err)
}
