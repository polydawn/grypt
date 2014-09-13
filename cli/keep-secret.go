package cli

import (
	"os"
	grypt "polydawn.net/grypt"
	"polydawn.net/grypt/gitutil"
	"polydawn.net/pogo/gosh"
)

var git = gosh.Sh("git")

func KeepSecret(ctx grypt.Context, keyring string, files []string) {
	// put git config.  probably already exists, but this should be an effectively idempotent set in that case.
	gitutil.PutGitFilterConfig(ctx)

	// check up front that all the secret files exist
	// this is racey with other checks later, but those later checks are done by git and come back to use as undifferentiated exit codes, so there's only so much we can do here.
	// exit if any of the secret files don't exist
	for _, file := range files {
		_, err := os.Stat(file)
		if err != nil {
			panic(err) // TODO: should probably exit politely with a well-known status code for file-not-found.  but needs infra: we can't literally os.Exit here because tests.
		}
	}

	// TODO: check if any of these files are staged or committed already, because that probably means the cleartext is in git history
	// we should abort and warn about that situation because it's very much not what you wanted.

	// write new '.gitattributes' config file
	gitattrs := gitutil.ReadRepoGitAttribs(ctx)
	for _, file := range files {
		gitattrs.PutGryptEntry(file)
	}
	gitattrs.SaveRepoGitAttribs(ctx)

	// stage it all
	git(gosh.DefaultIO)("add", "--", ".gitattributes")()
	for _, file := range files {
		git(gosh.DefaultIO)("add", "--", file)(gosh.Opts{OkExit: []int{0, 128}})()
	}
}
