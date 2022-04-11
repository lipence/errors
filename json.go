package errors

import (
	"bytes"
	"encoding/json"
	"sync"
)

var jsonEncBuffers = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GetBuffer() *bytes.Buffer {
	return jsonEncBuffers.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	jsonEncBuffers.Put(buf)
}

func toJSONMarshalable(err error) error {
	if _, ok := err.(json.Marshaler); ok {
		return err
	}
	return JsonErr{err}
}

type JsonErr struct {
	error
}

func (e JsonErr) MarshalJSON() ([]byte, error) {
	if e.error == nil {
		return []byte("null"), nil
	}
	return marshalJSONWithoutEscape(e.Error())
}

func marshalJSONWithoutEscape(source interface{}) ([]byte, error) {
	var buffer = GetBuffer()
	defer PutBuffer(buffer)
	var encoder = json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(source); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
