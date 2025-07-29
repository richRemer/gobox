package auth

import (
	"errors"
	"io/fs"
	"local/gobox/app"
	"local/gobox/repo"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
)

func Handler(keysDir string, users *repo.UserRepo) ssh.PublicKeyHandler {
	return func(ctx ssh.Context, key ssh.PublicKey) bool {
		keyFile := filepath.Join(keysDir, ctx.User()+"_keys")
		pem := key.Marshal()

		if loginAdmin(keyFile, key) {
			user, err := users.FindOrRegister(ctx.User())

			if err != nil {
				log.Error("Could not register admin", "user", ctx.User(), "error", err)
			}

			user.Role = app.RoleAdmin
			ctx.SetValue("user", user)
		} else {
			user, err := users.FindByPublicKey(string(pem))

			if err != nil {
				user = app.User{}
			}

			ctx.SetValue("user", user)
		}

		ctx.SetValue("publicKey", string(pem))

		return true
	}
}

func loginAdmin(keyFile string, key ssh.PublicKey) bool {
	bytes, err := os.ReadFile(keyFile)

	if errors.Is(err, fs.ErrNotExist) {
		return false
	} else if err != nil {
		log.Error("Could not read key", "file", keyFile, "error", err)
	}

	keys := map[string]ssh.PublicKey{}

	for len(bytes) > 0 {
		pub, _, _, rest, err := ssh.ParseAuthorizedKey(bytes)

		if err != nil {
			log.Error("Could not read key", "file", keyFile, "error", err)
			return false
		}

		keys[string(pub.Marshal())] = pub
		bytes = rest
	}

	_, ok := keys[string(key.Marshal())]
	return ok
}
