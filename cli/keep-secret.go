package cli

import (
	grypt "polydawn.net/grypt"
	"polydawn.net/grypt/gitutil"
)

func KeepSecret(ctx grypt.Context, keyring string, files []string) {
	gitattrs := gitutil.ReadRepoGitAttribs(ctx)
	println(gitattrs)
}
