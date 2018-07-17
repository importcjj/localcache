package localcache

import (
	"bytes"
	"testing"
)

func TestMemSink(t *testing.T) {
	sink := &memSink{}

	var testcases = []struct {
		Val       interface{}
		ExpectStr string
		ExpectErr error
	}{
		{[]byte("stringVal"), "stringVal", nil},
		{"stringVal", "stringVal", nil},
		{1, "1", nil},
		{map[string]interface{}{"name": 1}, `{"name":1}`, nil},
	}

	for _, tc := range testcases {
		sink.Reset()
		err := sink.SetValue(tc.Val)
		if err != tc.ExpectErr {
			t.Fatalf("expect %v, but got %v", tc.ExpectErr, err)
		}

		if b := sink.Bytes(); !bytes.Equal([]byte(tc.ExpectStr), b) {
			t.Fatalf("expect %s, but got %s", tc.ExpectStr, string(b))
		}
	}
}
