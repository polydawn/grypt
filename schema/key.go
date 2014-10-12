package schema

import (
	"bytes"
	"encoding/binary"
	"io"
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
	cipherKey []byte
	hmacKey   []byte
}

func (k Key) MarshalBinary() (data []byte, err error) {
	output := &bytes.Buffer{}

	binary.Write(output, binary.BigEndian, int8(1))

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

	var key1len int32
	if err := binary.Read(input, binary.BigEndian, &key1len); err != nil {
		return err
	}
	_, err := io.ReadAtLeast(input, k.cipherKey, int(key1len))
	if err != nil {
		return err
	}

	var key2len int32
	if err := binary.Read(input, binary.BigEndian, &key2len); err != nil {
		return err
	}
	_, err = io.ReadAtLeast(input, k.hmacKey, int(key1len))
	if err != nil {
		return err
	}

	return nil
}
