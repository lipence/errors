package errors

import (
	"bytes"
	"encoding/json"
	"sync"
)

var bytesBuffers = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func getBytesBuffer() *bytes.Buffer {
	return bytesBuffers.Get().(*bytes.Buffer)
}

func returnBytesBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bytesBuffers.Put(buf)
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
	var buffer = getBytesBuffer()
	defer returnBytesBuffer(buffer)
	var encoder = json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(source); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buffer.Bytes(), []byte{'\n'}), nil
}
