package gitutil

import (
	grypt "polydawn.net/grypt"
	"polydawn.net/pogo/gosh"
)

var git = gosh.Sh("git")

const exeName = "grypt" // not sure of the best way to make this more general.  taking '$0' isn't necessarily the most stable idea either.

func PutGitFilterConfig(ctx grypt.Context) {
	git("config", "filter.grypt.smudge", fmt.Sprintf("%s git-smudge", exeName))()
	git("config", "filter.grypt.clean", fmt.Sprintf("%s git-clean", exeName))()
	git("config", "diff.grypt.textconv", fmt.Sprintf("%s git-textconv", exeName))()
	// TODO: experiment with the 'required' config parameter, see `man gitattributues`.  making ourselves 'required' but exit 0 when a key is absent could give us desirable behavior like erroring when the command path isn't working.
}
