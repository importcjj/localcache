package localcache

import (
	"testing"
)

func TestBroker(t *testing.T) {
	var testcase = []struct {
		val interface{}
	}{
		{&UndefinedBroker{}},
		{&EmptyBroker{}},
	}

	for _, tc := range testcase {
		_, ok := tc.val.(Broker)
		if !ok {
			t.Fatalf("expect %v, got %v", true, ok)
		}
	}
}
