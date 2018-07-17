package localcache

import (
	"testing"
)

func TestGet(t *testing.T) {
	service := NewService(&Options{})

	key := "key"
	value := "Hello, world"

	service.Register(key, 0, func(key RequestKey, desc Sink) error {
		desc.SetBytes([]byte(value))
		return nil
	})

	getkey := &Key{
		GetterName: key,
		CacheName:  key,
		Value:      nil,
	}

	cache := service.Get(getkey)
	if s := cache.String(); s != value {
		t.Fatalf("got %s, expect %s", value, s)
	}
}
