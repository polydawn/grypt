package vault

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestVaultHeadersRoundtrip(t *testing.T) {
	Convey("Given some serialized headers", t, func() {
		headers := Headers{
			Header_grypt_scheme: "rot13",
			"a":                 "b",
			"c":                 "d",
		}
		serial, err := Content{
			Headers: headers,
		}.MarshalBinary()
		So(err, ShouldBeNil)

		Convey("We should get unmarshal the same headers back", func() {
			reheated := &Content{}
			err := reheated.UnmarshalBinary(serial)
			So(err, ShouldBeNil)
			So(reheated.Headers, ShouldResemble, headers)
		})
	})
}
