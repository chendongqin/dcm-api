package msgpack

import (
	"bytes"
	"github.com/vmihailenco/msgpack/v4"
)

type Writer struct {
	buf     bytes.Buffer
	encoder *msgpack.Encoder
	err     error
}

func NewWriter() *Writer {
	writer := &Writer{}
	writer.encoder = msgpack.NewEncoder(&writer.buf)
	return writer
}

func (w *Writer) Pack(v interface{}) *Writer {
	switch v.(type) {
	case int:
		w.err = w.encoder.EncodeInt(int64(v.(int)))
	case int8:
		w.err = w.encoder.EncodeInt(int64(v.(int8)))
	case int32:
		w.err = w.encoder.EncodeInt(int64(v.(int32)))
	case int64:
		w.err = w.encoder.EncodeInt(v.(int64))
	default:
		w.err = w.encoder.Encode(v)
	}
	return w
}

func (w *Writer) ToByteArray() []byte {
	return w.buf.Bytes()
}

func (w *Writer) ToString() string {
	return w.buf.String()
}

func (w *Writer) GetError() error {
	return w.err
}

func (w *Writer) IsEmpty() bool {
	return w.buf.Len() <= 0
}

func (w *Writer) Reset() *Writer {
	w.err = nil
	w.buf.Reset()
	return w
}
