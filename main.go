package main

import (
	"os"

	"github.com/zepif/EtherUSDC/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
