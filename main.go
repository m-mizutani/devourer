package main

import (
	"os"

	"github.com/m-mizutani/devourer/pkg/cli"
)

func main() {
	if err := cli.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
