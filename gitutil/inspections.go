package gitutil

import (
	"strings"
)

/*
	Produces a map of filenames currently staged, mapped to the git object ID of the staged blob
	(i.e., you can `git show <objid>` to get the raw content as staged).

	Some parsing here is currently sketchy.  In particular it is not safe in the face of filenames with tabs or linebreaks.
	This method is thus intended for testing and debugging more so than inclusion in production-ready flows.
*/
func ListStagedFiles() map[string]string {
	diffIndexLines := strings.Split(git("diff-index", "--cached", "HEAD").Output(), "\n")
	stagedObjectIds := make(map[string]string)
	for _, line := range diffIndexLines {
		splat := strings.Split(line, " ")
		if len(splat) != 5 {
			continue
		}
		filename := strings.Split(splat[4], "\t")[1]
		stagedObjectIds[filename] = splat[3]
	}
	return stagedObjectIds
}

/*
	Much like `ListStagedFiles()`, but maps the filenames to the contents of the staged files instead of just the blob objectsIds.

	This method does not pass through any smudge/clean/textconv filters -- it just produces the raw blob object, as it stands in git's dircache.
*/
func ListStagedFileContents() map[string][]byte {
	stagedContents := make(map[string][]byte)
	stagedObjectIds := ListStagedFiles()
	for filename, objid := range stagedObjectIds {
		stagedContents[filename] = []byte(git("show", objid).Output()) // FIXME: this cast to byte slice is ridiculous, should just read it as bytes the first time
	}
	return stagedContents
}
