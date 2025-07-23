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
	DB          string
	Host        string
	Port        int
}

func GetOptions() Options {
	var opts Options
	cwd, err := os.Getwd()

	if err != nil {
		log.Fatal("Could not get working directory.")
	}

	flag.StringVar(&opts.WorkingDir, "d", cwd, "working directory")
	flag.StringVar(&opts.KeysDir, "k", ".", "user keys directory")
	flag.StringVar(&opts.HostKeyFile, "H", "host_key", "host key file")
	flag.StringVar(&opts.DB, "D", ":memory:", "database name or path")
	flag.StringVar(&opts.Host, "h", "localhost", "listen address")
	flag.IntVar(&opts.Port, "p", 22, "listen port")
	flag.Parse()

	return opts
}
