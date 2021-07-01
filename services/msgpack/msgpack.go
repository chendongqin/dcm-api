package msgpack

// author: linyulin
// 封装msgpack常用方法

import (
	msgp "github.com/vmihailenco/msgpack/v4"
)

func UnpackInt32(b []byte) (int32, error) {
	i := new(int32)
	err := Unmarshal(b, i)
	return *i, err
}

func UnpackInt64(b []byte) (int64, error) {
	i := new(int64)
	err := Unmarshal(b, i)
	return *i, err
}

func UnpackString(b []byte) (string, error) {
	i := new(string)
	err := Unmarshal(b, i)
	return *i, err
}

func UnpackFloat64(b []byte) (float64, error) {
	i := new(float64)
	err := Unmarshal(b, i)
	return *i, err
}

func UnpackFloat32(b []byte) (float32, error) {
	i := new(float32)
	err := Unmarshal(b, i)
	return *i, err
}

func UnpackBool(b []byte) (bool, error) {
	boolean := new(bool)
	err := Unmarshal(b, boolean)
	return *boolean, err
}

func Unmarshal(b []byte, i interface{}) error {
	return msgp.Unmarshal(b, i)
}

func Marshal(v interface{}) ([]byte, error) {
	return msgp.Marshal(v)
}
