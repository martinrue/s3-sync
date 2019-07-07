package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/martinrue/s3-sync/sync"
)

var (
	dir     = flag.String("dir", "", "")
	bucket  = flag.String("bucket", "", "")
	region  = flag.String("region", "eu-west-1", "")
	ext     = flag.String("ext", "", "")
	silent  = flag.Bool("silent", false, "")
	help    = flag.Bool("help", false, "")
	version = flag.Bool("version", false, "")
)

func main() {
	log := func(message string, args ...interface{}) {
		if !*silent {
			fmt.Fprintln(os.Stderr, fmt.Sprintf(message, args...))
		}
	}

	usageAndExit := func(code int) {
		log("usage: sync --dir=<directory> --bucket=<bucket> --region=<region> --ext=<extensions> --silent")
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
		log("v0.0.1")
		os.Exit(0)
	}

	if *dir == "" || *bucket == "" {
		usageAndExit(1)
	}

	syncer := &sync.Syncer{
		Store: &sync.S3Store{Bucket: *bucket, Region: *region},
		Log:   log,
	}

	json, err := syncer.Run(*dir, *ext)
	if err != nil {
		log("error: %v", err)
		os.Exit(1)
	}

	fmt.Println(json)
}
