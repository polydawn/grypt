package cli

import (
	grypt "polydawn.net/grypt"
	"polydawn.net/grypt/gitutil"
	"polydawn.net/pogo/gosh"
)

var git = gosh.Sh("git")

func KeepSecret(ctx grypt.Context, keyring string, files []string) {
	gitattrs := gitutil.ReadRepoGitAttribs(ctx)
	for _, file := range files {
		gitattrs.PutGryptEntry(file)
	}
	gitattrs.SaveRepoGitAttribs(ctx)
	git(gosh.DefaultIO)("add", "--", ".gitattributes")()
	for _, file := range files {
		git(gosh.DefaultIO)("add", "--", file)(gosh.Opts{OkExit: []int{0, 128}})()
	}
}
