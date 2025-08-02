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
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	opts := cli.GetOptions()
	listenAddr := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
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
		wish.WithPublicKeyAuth(auth.Handler(keysDir, users)),
		wish.WithMiddleware(
			app.Middleware(users),
			activeterm.Middleware(),
			logging.Middleware(),
			elapsed.Middleware(),
		),
	)

	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	if server.ServerConfigCallback == nil {
		log.Info("setting config callback")
		server.ServerConfigCallback = func(ctx ssh.Context) *gossh.ServerConfig {
			config := &gossh.ServerConfig{}
			return config
		}
	}

	if server.ConnCallback == nil {
		log.Info("setting connection callback")
		server.ConnCallback = func(ctx ssh.Context, conn net.Conn) net.Conn {
			log.Info("connection callback")
			return conn
		}
	}

	forwardHandler := &ssh.ForwardedTCPHandler{}

	server.LocalPortForwardingCallback = func(ctx ssh.Context, host string, port uint32) bool {
		log.Info("accepting local forward to %s:%d", host, port)
		return true
	}

	server.ReversePortForwardingCallback = func(ctx ssh.Context, host string, port uint32) bool {
		log.Info("accepting reverse forward on %s:%d", host, port)
		return true
	}

	server.RequestHandlers = map[string]ssh.RequestHandler{
		"default": func(ctx ssh.Context, srv *ssh.Server, req *gossh.Request) (bool, []byte) {
			log.Info("default handler")
			return forwardHandler.HandleSSHRequest(ctx, srv, req)
		},
		"forward": func(ctx ssh.Context, srv *ssh.Server, req *gossh.Request) (bool, []byte) {
			log.Info("forward handler")
			return forwardHandler.HandleSSHRequest(ctx, srv, req)
		},
		"tcpip-forward": func(ctx ssh.Context, srv *ssh.Server, req *gossh.Request) (bool, []byte) {
			log.Info("tcpip-forward handler")
			return forwardHandler.HandleSSHRequest(ctx, srv, req)
		},
		"cancel-tcpip-forward": func(ctx ssh.Context, srv *ssh.Server, req *gossh.Request) (bool, []byte) {
			log.Info("cancel-tcpip-forward handler")
			return forwardHandler.HandleSSHRequest(ctx, srv, req)
		},
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", opts.Host, "port", opts.Port)

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
