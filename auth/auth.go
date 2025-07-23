package auth

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
)

func Handler(dir string) ssh.PublicKeyHandler {
	return func(ctx ssh.Context, key ssh.PublicKey) bool {
		file := filepath.Join(dir, ctx.User()+"_keys")
		bytes, err := os.ReadFile(file)

		if err != nil {
			return false
		}

		keys := map[string]ssh.PublicKey{}

		for len(bytes) > 0 {
			pub, _, _, rest, err := ssh.ParseAuthorizedKey(bytes)

			if err != nil {
				log.Error("Failed to parse %s: %v", file, err)
				return false
			}

			keys[string(pub.Marshal())] = pub
			bytes = rest
		}

		_, ok := keys[string(key.Marshal())]
		return ok
	}
}
