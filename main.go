package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/piodrs/sea/internal/app"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: sea file")
		os.Exit(1)
	}

	if err := app.Run(flag.Arg(0)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
