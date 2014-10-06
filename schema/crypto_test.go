package schema

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
	"testing/iotest"
)

var (
	plaintextSize = 1024

	plaintext = mkRand(plaintextSize)
	schemas   = []Schema{
		Aes256sha256ctr{},
	}
)

func mkRand(sz int) []byte {
	k := make([]byte, sz)
	io.ReadFull(rand.Reader, k)
	return k
}

func TestEncrypt(t *testing.T) {
	t.Logf("%25s: %.75s...\n", "plaintext", hex.EncodeToString(plaintext))
	for _, sch := range schemas {
		k, err := sch.NewKey(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		buf := new(bytes.Buffer)
		if err := sch.Encrypt(bytes.NewReader(plaintext), buf, k); err != nil {
			t.Fatal(err)
		}
		t.Logf("%25s: %.75s...\n", sch.Name(), hex.EncodeToString(buf.Bytes()))
	}
	return
}

func TestRoundTrip(t *testing.T) {
	t.Logf("%25s: %.75s...\n", "plaintext", hex.EncodeToString(plaintext))
	for _, sch := range schemas {
		// make key
		k, err := sch.NewKey(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}

		// encrypt
		ciphertext := new(bytes.Buffer)
		if err := sch.Encrypt(bytes.NewReader(plaintext), ciphertext, k); err != nil {
			t.Fatal(err)
		}

		// decrypt & verify
		reheatedtext := new(bytes.Buffer)
		err = sch.Decrypt(ciphertext, reheatedtext, k)
		t.Logf("%25s: %.75s...\n", sch.Name(), hex.EncodeToString(reheatedtext.Bytes()))
		if err != nil {
			t.Fatal(err)
		}

		// check cleartext match
		if !bytes.Equal(plaintext, reheatedtext.Bytes()) {
			t.Fatal(fmt.Errorf("cleartext match failure"))
		}
	}
}

func TestMACFailure(t *testing.T) {
	for _, sch := range schemas {
		k := Key{mkRand(sch.KeySize()), mkRand(sch.MACSize())}
		var err error
		buf := new(bytes.Buffer)
		if err := sch.Encrypt(bytes.NewReader(plaintext), iotest.TruncateWriter(buf, int64(plaintextSize-2)), k); err != nil {
			t.Fatal(err)
		}
		if err = sch.Decrypt(buf, ioutil.Discard, k); err == nil {
			t.Logf("This should have errored! %s", sch.Name())
			t.Fail()
		} else {
			t.Logf("%25s: %v\n", sch.Name(), err)
		}
	}
	return
}
