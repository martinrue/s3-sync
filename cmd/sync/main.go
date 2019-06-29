package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	dir     = flag.String("dir", "", "")
	bucket  = flag.String("bucket", "", "")
	help    = flag.Bool("help", false, "")
	version = flag.Bool("version", false, "")
)

func main() {
	usageAndExit := func(code int) {
		fmt.Fprint(os.Stderr, "usage: sync --dir=<directory> --bucket=<bucket>\n")
		os.Exit(code)
	}

	flag.Usage = func() {
		usageAndExit(1)
	}

	flag.Parse()

	if *help {
		usageAndExit(0)
	}

	if *version {
		fmt.Println("v0.0.1")
		os.Exit(0)
	}

	if *dir == "" || *bucket == "" {
		usageAndExit(1)
	}
}
