package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/elapsed"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "localhost"
)

func main() {
	var configdir string
	var port int

	flag.StringVar(&configdir, "dir", ".", "configuration directory")
	flag.StringVar(&configdir, "C", ".", "configuration directory")
	flag.IntVar(&port, "port", 22, "listen port for incoming connection")
	flag.IntVar(&port, "p", 22, "listen port for incoming connection")
	flag.Parse()

	server, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, strconv.Itoa(port))),
		wish.WithHostKeyPath(filepath.Join(configdir, "host_key")),
		wish.WithPublicKeyAuth(auth(configdir)),
		wish.WithMiddleware(
			middleware(),
			activeterm.Middleware(),
			logging.Middleware(),
			elapsed.Middleware(),
		),
	)

	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()

	if err := server.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}
