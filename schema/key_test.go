package schema

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestKeyMarshalling(t *testing.T) {
	Convey("Given a Key", t, func() {
		k := Key{
			Scheme:    Aes256sha256ctr{},
			cipherKey: []byte{12, 13, 14, 15},
			hmacKey:   []byte{45, 47, 48},
		}

		Convey("Marshalling should not explode", func() {
			serial, err := k.MarshalBinary()
			So(err, ShouldBeNil)

			Convey("Unmarshalling should return an equivalent key", func() {
				k2 := &Key{}
				err := k2.UnmarshalBinary(serial)
				So(err, ShouldBeNil)
				So(k2, ShouldResemble, k)
			})
		})
	})
}
