package main

import (
	"polydawn.net/pogo/gosh"
)

var git = gosh.Sh("git")

type Context struct {
	RepoDataDir string
	RepoWorkDir string
}

func DetectContext() Context {
	c := Context{}
	// TODO: better error handling here.  in general grypt should also operate even without a git repo.  but it should definitely not blow up.  and currently, it blows up.
	c.RepoDataDir = git("rev-parse", "--git-dir").Output()       // warning: this has ridiculously erratic formating.  it may be relative to your cwd, or it may be absolute (if your cwd is deeper than the git data dir).
	c.RepoWorkDir = git("rev-parse", "--show-toplevel").Output() // appears to always produce an absolute path
	return c
}
