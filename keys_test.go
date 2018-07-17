package localcache

import (
	"testing"
)

func TestKey(t *testing.T) {
	expect := &Key{
		GetterName: "GetterName",
		CacheName:  "CacheName",
		Value:      map[string]string{"x": "y"},
	}

	var k interface{} = expect

	requestK, ok := k.(RequestKey)
	if !ok {
		t.Fatalf("expect %v, but got %v", true, ok)
	}

	if val := requestK.GetGetterName(); val != expect.GetterName {
		t.Fatalf("expect %v, but got %v", expect, val)
	}

	if val := requestK.GetCacheName(); val != expect.CacheName {
		t.Fatalf("expect %v, but got %v", expect, val)
	}

	val := requestK.GetValue()
	m, ok := val.(map[string]string)
	if !ok {
		t.Fatalf("expect %v, but got %v", true, ok)
	}

	if m["x"] != "y" {
		t.Fatalf("expect %v, but got %v", "y", m["x"])
	}
}
