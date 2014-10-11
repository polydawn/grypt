package pem

import (
	"encoding/pem"
	"testing"
	"strings"

	. "github.com/smartystreets/goconvey/convey"
)

/*
	Do The Right Thing to a multiline string literal indented with your code.

	More literally: take a string and consume the leading line break,
	and leading tabs in count matching the first line's indentation,
*/
func lit(s string) string {
	// the leading line is probably the break before the content, so skip that
	s = strings.TrimPrefix(s, "\n")
	
	// figure out how many tabs are on the first line.  this is the indentation baseline we expect to strip for the rest.
	depth := 0
	for ; s[depth] == '\t' ; depth++ {}
	baseline := strings.Repeat("\t", depth)
	
	// strip up to that number of indents from each line.
	lines := strings.Split(s, "\n")
	linecount := len(lines)
	output := make([]string, linecount)
	for n, s := range lines {
		if strings.HasPrefix(s, baseline) {
			output[n] = s[depth:]
		} else {
			output[n] = s
		}
	}

	// the last line might have fewer, and we don't want those, but we still probably do want a trailing break
	if strings.Count(output[linecount-1], "\t") == len(output[linecount-1])  {
		output[linecount-1] = ""
	}

	return strings.Join(output, "\n")
}

func TestPemFormatBasics(t *testing.T) {
	Convey("Given some bananas", t, func() {
		block := &pem.Block{
			Type: "GRYPT CIPHERTEXT HEADER",
			Headers: map[string]string{},
			// pem.Block.Bytes is a zero value for us, we're not gonna use b64
		}

		Convey("We should have a party", func() {
			So(string(pem.EncodeToMemory(block)), ShouldEqual, lit(`
				-----BEGIN GRYPT CIPHERTEXT HEADER-----
				-----END GRYPT CIPHERTEXT HEADER-----
			`))
		})
	})
}

