package localcache

import (
	"encoding/json"
)

import (
	"errors"
	"time"
)

// ErrGetterUndefined ...
var ErrGetterUndefined = errors.New("cache getter undefined")

// GetterFunc ...
type GetterFunc func(key RequestKey, sink Sink) error

type getter struct {
	Func   GetterFunc
	Expire time.Duration
}

// Sink ...
type Sink interface {
	SetValue(interface{}) error
	SetBytes([]byte) (int, error)

	Bytes() []byte
}

type memSink struct {
	Key  RequestKey
	data []byte
}

func (sink *memSink) Bytes() []byte {
	return sink.data
}

func (sink *memSink) SetValue(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	sink.data = data
	return nil
}

func (sink *memSink) SetBytes(b []byte) (int, error) {
	sink.data = b
	return len(b), nil
}
