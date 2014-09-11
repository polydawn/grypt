package cli

import (
	grypt "polydawn.net/grypt"
)

func KeepSecret(ctx grypt.Context, keyring string, files []string) {
	PutGitFilterConfig(ctx)
}
