package gitutil

type Context struct {
	GryptName    string // name of the current executable (treat it sort of like "$0" in bash).  probably equal to "grypt" (may be an absolute path in tests to a temporary binary).
	GryptVersion string // version identifier of the currently operating version of grypt.  maybe be useful to serialize into headers.
	Keyring      string // name of the keyring this command is operating with.  any grypt command only operates with one keyring at a time.
	RepoDataDir  string
	RepoWorkDir  string
}

func DetectContext() Context {
	c := Context{}
	// TODO: better error handling here.  in general grypt should also operate even without a git repo.  but it should definitely not blow up.  and currently, it blows up.
	c.RepoDataDir = git("rev-parse", "--git-dir").Output()       // warning: this has ridiculously erratic formating.  it may be relative to your cwd, or it may be absolute (if your cwd is deeper than the git data dir).
	c.RepoDataDir = c.RepoDataDir[0 : len(c.RepoDataDir)-1]      // there's a trailing '\n' on the output of rev-parse.  remove it.
	c.RepoWorkDir = git("rev-parse", "--show-toplevel").Output() // appears to always produce an absolute path
	c.RepoWorkDir = c.RepoWorkDir[0 : len(c.RepoWorkDir)-1]      // there's a trailing '\n' on the output of rev-parse.  remove it.
	return c
}

// TODO: this needs some tests itself to make sure it does something meaningful in an empty dir.
