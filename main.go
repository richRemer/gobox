package main

import (
	"context"
	_ "embed"
	"errors"
	"local/gobox/app"
	"local/gobox/auth"
	"local/gobox/cli"
	"local/gobox/repo"
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
	opts := cli.GetOptions()
	listenAddr := net.JoinHostPort(host, strconv.Itoa(opts.Port))
	hostKeyFile := filepath.Join(opts.WorkingDir, opts.HostKeyFile)
	keysDir := filepath.Join(opts.WorkingDir, opts.KeysDir)
	users, err := repo.OpenUsers(opts.DB)

	if err != nil {
		log.Fatal("Could not open user database", "error", err)
	}

	defer users.Close()

	if err := users.InitSchemaIfNeeded(); err != nil {
		log.Fatal("Could not initialize database", "error", err)
	}

	server, err := wish.NewServer(
		wish.WithAddress(listenAddr),
		wish.WithHostKeyPath(hostKeyFile),
		wish.WithPublicKeyAuth(auth.Handler(keysDir)),
		wish.WithMiddleware(
			app.Middleware(),
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
	log.Info("Starting SSH server", "host", host, "port", opts.Port)

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
