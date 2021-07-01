package msgpack

import (
	"github.com/astaxie/beego/logs"
	"github.com/bmizerany/assert"
	"testing"
)

func TestUnpackInt32(t *testing.T) {
	var i int32 = 51
	b, err := Marshal(i) //1 bytes to got.
	a, err := UnpackInt32(b)
	if err != nil {
		logs.Info(err)
	}
	assert.Equal(t, i, a)
}

func TestUnpackInt64(t *testing.T) {
	var i int64 = 65535
	b, err := Marshal(i)
	a, err := UnpackInt64(b)
	if err != nil {
		logs.Info(err)
	}
	assert.Equal(t, i, a)
}

func TestUnpackInt32ByBytes(t *testing.T) {
	var v1, v2, v3 int8
	v1 = -51
	v2 = -1
	v3 = -1
	//if v1 < 0 then byte(v1) == byte(256 + int(v1))
	b := []byte{byte(v1), byte(v2), byte(v3)}
	a, err := UnpackInt32(b)
	if err != nil {
		logs.Info(err)
	}
	assert.Equal(t, int32(65535), a)
}

func TestUnpackInt64ByBytes(t *testing.T) {
	var v1, v2, v3 int8
	v1 = -51
	v2 = -1
	v3 = -1
	b := []byte{byte(v1), byte(v2), byte(v3)}
	a, err := UnpackInt64(b)
	if err != nil {
		logs.Info(err)
	}
	assert.Equal(t, int64(65535), a)
}
