package localcache

import (
	"errors"
)

var ErrBrokerUndefined = errors.New("sorry! localcache's broker undefined")

// Broker ...
type Broker interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
}

type UndefinedBroker struct{}

func (*UndefinedBroker) Get(_ string) ([]byte, error) { return nil, ErrBrokerUndefined }
func (*UndefinedBroker) Set(_ string, _ []byte) error { return ErrBrokerUndefined }

type EmptyBroker struct{}

func (*EmptyBroker) Get(_ string) ([]byte, error) { return nil, ErrBrokerUndefined }
func (*EmptyBroker) Set(_ string, _ []byte) error { return nil }
