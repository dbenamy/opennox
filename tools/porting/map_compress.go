package main

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy/cnxz"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		panic("source and destination paths required")
	}
	if err := cnxz.CompressFile(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
