package cli

import (
	"flag"
	"log"
	"os"
)

type Options struct {
	WorkingDir  string
	KeysDir     string
	HostKeyFile string
	Port        int
	DB          string
}

func GetOptions() Options {
	var opts Options
	cwd, err := os.Getwd()

	if err != nil {
		log.Println("Could not get working directory.")
	}

	flag.StringVar(&opts.WorkingDir, "d", cwd, "working directory")
	flag.StringVar(&opts.KeysDir, "k", ".", "user keys directory")
	flag.StringVar(&opts.HostKeyFile, "h", "host_key", "host key file")
	flag.IntVar(&opts.Port, "p", 22, "listen port")
	flag.StringVar(&opts.DB, "D", ":memory:", "database name or path")
	flag.Parse()

	return opts
}
