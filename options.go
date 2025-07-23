package main

import (
	"flag"
	"log"
	"os"
)

type options struct {
	WorkingDir  string
	KeysDir     string
	HostKeyFile string
	Port        int
}

func getopts() options {
	var opts options
	cwd, err := os.Getwd()

	if err != nil {
		log.Println("Could not get working directory.")
	}

	flag.StringVar(&opts.WorkingDir, "d", cwd, "working directory")
	flag.StringVar(&opts.KeysDir, "k", ".", "user keys directory")
	flag.StringVar(&opts.HostKeyFile, "h", "host_key", "host key file")
	flag.IntVar(&opts.Port, "p", 22, "listen port")
	flag.Parse()

	return opts
}
