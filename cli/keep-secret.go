package cli

import (
	grypt "polydawn.net/grypt"
	"polydawn.net/grypt/gitutil"
)

func KeepSecret(ctx grypt.Context, keyring string, files []string) {
	gitutil.PutGitFilterConfig(ctx)
}
