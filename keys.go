package localcache

// RequestKey ...
type RequestKey interface {
	GetGetterName() string
	GetCacheName() string
	GetValue() interface{}
}

// Key ...
type Key struct {
	GetterName string
	CacheName  string
	Value      interface{}
}

// GetGetterName ...
func (k *Key) GetGetterName() string {
	return k.GetterName
}

// GetCacheName ...
func (k *Key) GetCacheName() string {
	return k.CacheName
}

// GetValue ...
func (k *Key) GetValue() interface{} {
	return k.Value
}
