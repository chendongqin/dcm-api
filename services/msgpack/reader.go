package msgpack

import (
	"bytes"
	"github.com/vmihailenco/msgpack/v4"
)

type Decoder struct {
	*msgpack.Decoder
}

func NewDecoder(input []byte) *Decoder {
	return &Decoder{msgpack.NewDecoder(bytes.NewReader(input))}
}
