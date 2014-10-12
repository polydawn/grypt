package gitutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"polydawn.net/pogo/gosh"
)

var git = gosh.Sh("git")

var br = []byte("\n")

func PutGitFilterConfig(ctx Context) {
	git("config", "filter.grypt.smudge", fmt.Sprintf("%s git-smudge %%f", ctx.GryptName))()
	git("config", "filter.grypt.clean", fmt.Sprintf("%s git-clean %%f", ctx.GryptName))()
	// TODO: experiment with the 'required' config parameter, see `man gitattributues`.  making ourselves 'required' but exit 0 when a key is absent could give us desirable behavior like erroring when the command path isn't working.
	// if you have these filter configs in place but not required and the grypt command missing, you get very ugly error messages from git anyway.
	//git("config", "filter.grypt.required", "true")()
	git("config", "diff.grypt.textconv", fmt.Sprintf("%s git-textconv %%f", ctx.GryptName))()
}

/*
	Manipulable view of a gitattributes file.

	We track this as a messy glob of byte slices so that we can save it without causing diffs;
	we are not necessarily the only actor in a gitattributes file.
	This is inefficent, but gitattributes files are also realistically never expected to be
	more than a few kilobytes, so multiple searches are not a cause for irritation.
*/
type Gitattribs struct {
	Lines []GitattribLine
}

type GitattribLine struct {
	Pattern string // the first part of the line, which identifies the fileset the rule acts on.  may be nil.  called the pattern because that's what `man gitattributes` calls it.
	Raw     []byte // the whole line in its original form (so we can save it again)
}

var rPattern = regexp.MustCompile("^[^\\s]*")

func ReadRepoGitAttribs(ctx Context) *Gitattribs {
	return ReadGitAttribsFile(filepath.Join(ctx.RepoWorkDir, ".gitattributes"))
}

func ReadGitAttribsFile(filename string) *Gitattribs {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &Gitattribs{}
		} else {
			panic(err)
		}
	}
	return ParseGitAttribs(raw)
}

func ParseGitAttribs(raw []byte) *Gitattribs {
	rawLines := bytes.Split(raw, br)
	ga := &Gitattribs{
		Lines: make([]GitattribLine, len(rawLines)),
	}
	for i, line := range rawLines {
		gapattern := rPattern.Find(line)
		ga.Lines[i] = GitattribLine{
			Pattern: string(gapattern),
			Raw:     line,
		}
	}
	return ga
}

func (ga *Gitattribs) Marshall() []byte {
	lines := make([][]byte, len(ga.Lines))
	for i, line := range ga.Lines {
		lines[i] = line.Raw
	}
	return bytes.Join(lines, br)
}

func (ga *Gitattribs) SaveFile(filename string) {
	if err := ioutil.WriteFile(filename, ga.Marshall(), 0644); err != nil {
		panic(err)
	}
}

func (ga *Gitattribs) SaveRepoGitAttribs(ctx Context) {
	ga.SaveFile(filepath.Join(ctx.RepoWorkDir, ".gitattributes"))
}

func (ga *Gitattribs) PutGryptEntry(path string) {
	// currently this is a naive implementation that assumes you have no other attributes for the files we're keeping secret; there is no attempt to retain existing attributes.
	// also, god have mercy on your soul if your secret files have whitespace characters in their path.  afaict the format of gitattributes files is woefullly unprepared for that concept (though i'd love to be corrected).
	putLine := []byte(fmt.Sprintf("%s filter=grypt diff=grypt", path))
	for i, line := range ga.Lines {
		if line.Pattern == path {
			ga.Lines[i].Raw = putLine
			return
		}
	}
	ga.Lines = append(ga.Lines, GitattribLine{Pattern: path, Raw: putLine})
}
