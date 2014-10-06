package schema

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"testing"
	"testing/iotest"
)

var (
	plaintextSize = 1024
	out           [][]byte

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
		out = append(out, buf.Bytes())
	}
	return
}

func TestDecrypt(t *testing.T) {
	t.Logf("%25s: %.75s...\n", "plaintext", hex.EncodeToString(plaintext))
	for i, sch := range schemas {
		k := Key{mkRand(sch.KeySize()), mkRand(sch.MACSize())}
		x := new(bytes.Buffer)
		if err := sch.Decrypt(bytes.NewReader(out[i]), x, k); err != nil {
			t.Fatal(err)
		}
		t.Logf("%25s: %.75s...\n", sch.Name(), hex.EncodeToString(x.Bytes()))
		if !bytes.Equal(plaintext, x.Bytes()) {
			t.Fail()
		}
	}
	return
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
