package localcache

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

// Cache ...
type Cache interface {
	CreateTime() time.Time
	Scan(v interface{}) error
	GetError() error
	String() string
	Bytes() []byte

	LastModified() string
	MD5() string
}

type cache struct {
	Data      []byte
	ErrString string
	Timestamp int64
	MD5       string
}

// CacheWrapper ...
type CacheWrapper struct {
	cache *cache
}

func newCacheWrapper() *CacheWrapper {
	return &CacheWrapper{
		cache: &cache{},
	}
}

// GetError ...
func (wrapper *CacheWrapper) GetError() error {
	if len(wrapper.cache.ErrString) > 0 {
		return errors.New(wrapper.cache.ErrString)
	}

	return nil
}

// MD5 ...
func (wrapper *CacheWrapper) MD5() string {
	return wrapper.cache.MD5
}

// LastModified ...
func (wrapper *CacheWrapper) LastModified() string {
	return wrapper.CreateTime().Format(time.RFC1123)
}

// Bytes ...
func (wrapper *CacheWrapper) Bytes() []byte {
	if len(wrapper.cache.ErrString) > 0 {
		return []byte(wrapper.cache.ErrString)
	}

	return wrapper.cache.Data
}

func (wrapper *CacheWrapper) String() string {
	if len(wrapper.cache.ErrString) > 0 {
		return wrapper.cache.ErrString
	}

	return string(wrapper.cache.Data)
}

// Scan ...
func (wrapper *CacheWrapper) Scan(v interface{}) error {
	return json.Unmarshal(wrapper.Bytes(), v)
}

// CreateTime ...
func (wrapper *CacheWrapper) CreateTime() time.Time {
	return time.Unix(wrapper.cache.Timestamp, 0)
}

func (wrapper *CacheWrapper) setError(err error) {
	if err != nil {
		wrapper.cache.ErrString = err.Error()
	} else {
		wrapper.cache.ErrString = ""
	}
}

func (wrapper *CacheWrapper) setCacheData(b []byte) {
	if len(b) > 0 {
		wrapper.cache.Data = b
		h := md5.New()
		h.Write(b)
		wrapper.cache.MD5 = hex.EncodeToString(h.Sum(nil))
	}

	wrapper.cache.Timestamp = time.Now().Unix()
}

func (wrapper *CacheWrapper) checkExpire(expire time.Duration) bool {
	if expire <= 0 {
		return false
	}

	return wrapper.CreateTime().Add(expire).Before(time.Now())
}

func (wrapper *CacheWrapper) encode() []byte {
	var buffer bytes.Buffer
	err := gob.NewEncoder(&buffer).Encode(wrapper.cache)
	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (wrapper *CacheWrapper) decode(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	var buffer = bytes.NewBuffer(b)
	err := gob.NewDecoder(buffer).Decode(wrapper.cache)
	return err
}
