package gitutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	grypt "polydawn.net/grypt"
	"polydawn.net/pogo/gosh"
	"regexp"
)

var git = gosh.Sh("git")

var br = []byte("\n")

const exeName = "grypt" // not sure of the best way to make this more general.  taking '$0' isn't necessarily the most stable idea either.

func PutGitFilterConfig(ctx grypt.Context) {
	git("config", "filter.grypt.smudge", fmt.Sprintf("%s git-smudge %%f", exeName))()
	git("config", "filter.grypt.clean", fmt.Sprintf("%s git-clean %%f", exeName))()
	git("config", "diff.grypt.textconv", fmt.Sprintf("%s git-textconv %%f", exeName))()
	// TODO: experiment with the 'required' config parameter, see `man gitattributues`.  making ourselves 'required' but exit 0 when a key is absent could give us desirable behavior like erroring when the command path isn't working.
}

func PutAttributesEntry(ctx grypt.Context) {
	// currently this is a naive implementation that assumes you have no other attributes for the files we're keeping secret; there is no attempt to retain existing attributes.
	// also, god have mercy on your soul if your secret files have whitespace characters in their path.  afaict the format of gitattributes files is woefullly unprepared for that concept (though i'd love to be corrected).

	gitattrs, err := ioutil.ReadFile(filepath.Join(ctx.RepoWorkDir, ".gitattributes"))
	if err != nil {
		panic(err)
	}
	gitattrLines := bytes.Split(gitattrs, br)

	for _, line := range gitattrLines {
		r, _ := regexp.Compile("^(.*)\\s")
		gapattern := r.Find([]byte(line))
		println(gapattern) // appease the golang compiler while i keep working -.-
	}
}


/*
	Manipulable view of a gitattributes file.

	We track this as a string so that we can save it without causing diffs;
	we are not necessarily the only actor in a gitattributes file.
	This is inefficent, but gitattributes files are also realistically never expected to be
	more than a few kilobytes, so multiple searches are not a cause for irritation.
*/
type gitattribs struct {
	x byte[]
}

func (ga *gitattribs) PutGryptEntry(path string) {

}
