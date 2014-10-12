package vault

import (
	"bytes"
	"crypto/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"polydawn.net/grypt/schema"
)

func TestVaultRoundtrip(t *testing.T) {
	Convey("Given a cleartext and choice of cipher scheme", t, func() {
		sch := schema.Aes256sha256ctr{}
		k, err := sch.NewKey(rand.Reader)
		So(err, ShouldBeNil)
		str := "cleartext! :D"
		cleartext := bytes.NewBuffer([]byte(str))

		Convey("Vault should produce a ciphertext stream", func() {
			ciphertext := &bytes.Buffer{}
			SealEnvelope(cleartext, ciphertext, k)

			Convey("Vault should be able to return the cleartext", func() {
				reheated := &bytes.Buffer{}
				headers := OpenEnvelope(ciphertext, reheated, k)

				So(headers[Header_grypt_version], ShouldEqual, "1.0")
				So(headers[Header_grypt_scheme], ShouldEqual, sch.Name())
				So(headers[Header_grypt_keyring], ShouldEqual, "default")

				So(len(reheated.Bytes()), ShouldEqual, len([]byte(str)))
				So(string(reheated.Bytes()), ShouldEqual, str)
			})
		})
	})
}
