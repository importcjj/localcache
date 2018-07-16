package singlecache

import (
	"time"
)

// Options ...
type Options struct {
	Broker     Broker
	GetTimeout time.Duration
}

func (o *Options) merge(opts *Options) {
	if opts.Broker != nil {
		o.Broker = opts.Broker
	}

	if opts.GetTimeout != 0 {
		o.GetTimeout = opts.GetTimeout
	}
}

func newDefaultOptions() *Options {
	return &Options{
		Broker:     &EmptyBroker{},
		GetTimeout: 10 * time.Second,
	}
}
