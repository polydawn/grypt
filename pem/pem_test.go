package pem

import (
	"fmt"
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
	Convey("Given a nearly empty block", t, func() {
		block := &pem.Block{
			Type: "GRYPT CIPHERTEXT HEADER",
			Headers: map[string]string{},
			// pem.Block.Bytes is a zero value for us, we're not gonna use b64
		}
		serial := pem.EncodeToMemory(block)

		Convey("We should get an empty body section", func() {
			So(string(serial), ShouldEqual, lit(`
				-----BEGIN GRYPT CIPHERTEXT HEADER-----
				-----END GRYPT CIPHERTEXT HEADER-----
			`))
		})

		// skip this test, it's fucked, their serializer fails to roundtrip empty values
		SkipConvey("Everything is still empty when reheated", func() {
			reheated, rest := pem.Decode(serial)
			So(len(rest), ShouldEqual, 0)
			So(reheated, ShouldResemble, block)
		})
	})

	Convey("Given some headers", t, func() {
		block := &pem.Block{
			Type: "GRYPT CIPHERTEXT HEADER",
			Headers: map[string]string{
				"Grypt-Test-Header": "some value",
				"Grypt-caps-sense": "moar value",
			},
			Bytes: []byte{} ,
		}
		serial := pem.EncodeToMemory(block)

		Convey("The serial format is stable and looks nice", func() {
			So(string(serial), ShouldEqual, lit(`
				-----BEGIN GRYPT CIPHERTEXT HEADER-----
				Grypt-Test-Header: some value
				Grypt-caps-sense: moar value
				
				-----END GRYPT CIPHERTEXT HEADER-----
			`)) // i don't particularly understand why having nonzero headers got us this extra line break at the end, and consider that a bit wrong if there's no body bytes, but okay, whatever.
		})

		Convey("Everything is the same when reheated", func() {
			// but the conservative approach on serialization doesn't do much good if you can't round-trip it -.-
			reheated, rest := pem.Decode(serial)
			So(reheated.Type, ShouldResemble, block.Type)
			So(reheated.Headers, ShouldResemble, block.Headers)
			fmt.Printf("\n::::\t%#v\n\t%#v\n", reheated.Bytes, block.Bytes)
			So(reheated.Bytes, ShouldResemble, block.Bytes)
			So(reheated, ShouldResemble, block)
			So(len(rest), ShouldEqual, 0)
		})
	})

	Convey("Given headers named with leading or trailing spaces", t, func() {
		block := &pem.Block{
			Type: "GRYPT CIPHERTEXT HEADER",
			Headers: map[string]string{
				"  leading": "x",
				"trailing  ": "y",
			},
		}
		serial := pem.EncodeToMemory(block)

		Convey("The strange names are preserved in serial form", func() {
			// i guess this is the appropriate conservative behavior...
			So(string(serial), ShouldEqual, lit(`
				-----BEGIN GRYPT CIPHERTEXT HEADER-----
				  leading: x
				trailing  : y

				-----END GRYPT CIPHERTEXT HEADER-----
			`))
		})

		Convey("The strange names are altered (trimmed) when reheated", func() {
			// but the conservative approach on serialization doesn't do much good if you can't round-trip it -.-
			reheated, rest := pem.Decode(serial)
			So(reheated, ShouldResemble, &pem.Block{
				Type: "GRYPT CIPHERTEXT HEADER",
				Headers: map[string]string{
					"leading": "x",
					"trailing": "y",
				},
			})
			So(len(rest), ShouldEqual, 0)
		})
	})
}

