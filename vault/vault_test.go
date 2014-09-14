package vault

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestVaultHeadersRoundtrip(t *testing.T) {
	Convey("Given some valid serialized headers", t, func() {
		headers := Headers{
			Header_grypt_scheme: "rot13",
			"A":                 "b",
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

	Convey("Given some mix of valid and invalid serialized headers", t, func() {
		headers := Headers{
			Header_grypt_scheme: " rot13 ",
			"A":                 "b",
			"c":                 "d",
			"clearly not":       "d",
		}
		serial, err := Content{
			Headers: headers,
		}.MarshalBinary()
		So(err, ShouldBeNil)

		reheated := &Content{}
		err = reheated.UnmarshalBinary(serial)
		So(err, ShouldBeNil)

		Convey("Leading and trailing whitespace should be trimmed", func() {
			So(reheated.Headers[Header_grypt_scheme], ShouldEqual, "rot13")
		})

		Convey("Invalid header entries should be absent", func() {
			_, err := reheated.Headers["c"]
			So(err, ShouldNotBeNil)
			_, err = reheated.Headers["clearly not"]
			So(err, ShouldNotBeNil)
			// So(len(reheated.Headers), ShouldEqual, 2) // rong, because of the forced headers
		})
	})
}
