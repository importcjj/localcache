package singlecache

import (
	"testing"
)

func TestGet(t *testing.T) {
	service := NewService(&Options{})

	key := "key"
	value := "Hello, world"

	service.Register(key, 0, func(key string, desc Sink) error {
		desc.SetBytes([]byte(value))
		return nil
	})

	cache := service.Get(key)
	if s := cache.String(); s != value {
		t.Fatalf("got %s, expect %s", value, s)
	}

}
