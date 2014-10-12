package schema

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io"
	"os"
)

/*
	Key struct stores the two byte slices for most symmetric crypto operations:
	the cipher key and the hmac key.

	This is a simplifying assumption for all the interfaces we currently use, but may break
	for other kinds of (very) exotic cipher suites we don't yet support.

	Keys implement serialization -- roughly.  Since we don't actually yet have a need for
	complex keys or polymorphic types to represent them, they're just a sequence of length delimited
	byte slices.  Verification of type can pretty much either A) not be done or B) done by just
	trying to use it; thus, there is no metadata included that pretends otherwise.  We do however
	include a leading byte for "version" (or "type") to give an upgrade path should this change in the future.

	@implements encoding.BinaryMarshaler
	@implements encoding.BinaryUnmarshaler
*/
type Key struct {
	Scheme    Schema
	cipherKey []byte
	hmacKey   []byte
}

func (k Key) MarshalBinary() (data []byte, err error) {
	output := &bytes.Buffer{}

	binary.Write(output, binary.BigEndian, int8(1))

	schemeBytes := []byte(k.Scheme.Name())
	binary.Write(output, binary.BigEndian, int32(len(schemeBytes)))
	output.Write(schemeBytes)

	binary.Write(output, binary.BigEndian, int32(len(k.cipherKey)))
	output.Write(k.cipherKey)

	binary.Write(output, binary.BigEndian, int32(len(k.hmacKey)))
	output.Write(k.hmacKey)

	return output.Bytes(), nil
}

func (k *Key) UnmarshalBinary(data []byte) error {
	input := bytes.NewBuffer(data)

	var version int8
	if err := binary.Read(input, binary.BigEndian, &version); err != nil {
		return err
	}

	var schemelen int32
	if err := binary.Read(input, binary.BigEndian, &schemelen); err != nil {
		return err
	}
	schemeBytes := make([]byte, schemelen)
	_, err := io.ReadAtLeast(input, schemeBytes, int(schemelen))
	if err != nil {
		return err
	}
	k.Scheme = ParseSchema(string(schemeBytes))

	var key1len int32
	if err := binary.Read(input, binary.BigEndian, &key1len); err != nil {
		return err
	}
	k.cipherKey = make([]byte, key1len)
	_, err = io.ReadAtLeast(input, k.cipherKey, int(key1len))
	if err != nil {
		return err
	}

	var key2len int32
	if err := binary.Read(input, binary.BigEndian, &key2len); err != nil {
		return err
	}
	k.hmacKey = make([]byte, key2len)
	_, err = io.ReadAtLeast(input, k.hmacKey, int(key2len))
	if err != nil {
		return err
	}

	return nil
}

/*
	base64 encode and write key 'k' to file 'f'
*/
func WriteKey(f string, k Key) error {
	bits, err := k.MarshalBinary()
	if err != nil {
		return err
	}
	file, err := os.Create(f)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := base64.NewEncoder(base64.StdEncoding, file)
	_, err = bytes.NewBuffer(bits).WriteTo(enc)
	if err != nil {
		return err
	}
	enc.Close()
	return nil
}

/*
	read and decode a key from file 'f'
*/
func ReadKey(f string) (Key, error) {
	k := Key{}
	bits := new(bytes.Buffer)
	file, err := os.Open(f)
	if err != nil {
		return Key{}, err
	}
	defer file.Close()
	dec := base64.NewDecoder(base64.StdEncoding, file)
	_, err = bits.ReadFrom(dec)
	if err != nil {
		return Key{}, err
	}
	err = k.UnmarshalBinary(bits.Bytes())
	if err != nil {
		return Key{}, err
	}
	return k, nil
}
