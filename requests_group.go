package localcache

import (
	"errors"
)

var (
	// ErrCacheBrokerTimeout ...
	ErrCacheBrokerTimeout = errors.New("cache broker timeout")
	// ErrCacheGetterTimeout ...
	ErrCacheGetterTimeout = errors.New("cache getter timeout")
)

type request struct {
	Key string
}

type requestGroup struct {
	cacheWrapper *CacheWrapper
	timeout      bool
	closeQ       chan struct{}
}

func newRequestGroup() *requestGroup {
	group := &requestGroup{
		cacheWrapper: newCacheWrapper(),
		closeQ:       make(chan struct{}),
	}

	return group
}

func (group *requestGroup) HandleGet(service *Service, key RequestKey) {
	getter, ok := service.getters[key.GetGetterName()]
	if !ok {
		group.cacheWrapper.setError(ErrGetterUndefined)
		service.getEventFinishedQ <- &event{Key: key}
		return
	}

	timeout := service.registerGetTimeout(key, service.GetTimeout)
	err := group.getFromBroker(service, key, getter)
	if !timeout.Stop() {
		service.Logger.Warn("get from broker timeout")
		return
	}

	if err == nil {
		if group.cacheWrapper.checkExpire(getter.Expire) {
			service.registerUpdateEvent(key)
		}
	} else {
		service.Logger.Warnf("get from broker failed with %v, will try upstream", err)

		timeout = service.registerGetTimeout(key, service.GetTimeout)
		err = group.getFromGetter(service, key, getter)
		if !timeout.Stop() {
			service.Logger.Warn("get from upstream timeout")
			return
		}
	}

	service.getEventFinishedQ <- &event{Key: key}
}

func (group *requestGroup) handleUpdate(service *Service, key RequestKey, getter *getter) {
	timeout := service.registerUpdateTimeout(key, service.GetTimeout)
	err := group.getFromGetter(service, key, getter)
	if !timeout.Stop() {
		service.Logger.Warn("get from upstream timeout")
		return
	}

	if err != nil {
		service.Logger.Warn("get from upstream failed", err)
	}

	service.updateEventFinishQ <- &event{Key: key}
}

func (group *requestGroup) getFromBroker(service *Service, key RequestKey, getter *getter) error {
	b, err := service.Broker.Get(key.GetCacheName())
	group.cacheWrapper.setError(err)
	if err != nil {
		return err
	}

	return group.cacheWrapper.decode(b)
}

func (group *requestGroup) getFromGetter(service *Service, key RequestKey, getter *getter) error {
	dest := &memSink{Key: key}
	err := getter.Func(key, dest)
	group.cacheWrapper.setError(err)
	if err != nil {
		return err
	}
	group.cacheWrapper.setCacheData(dest.Bytes())

	return service.Broker.Set(key.GetCacheName(), group.cacheWrapper.encode())
}

// Finish ...
func (group *requestGroup) Finish() {
	close(group.closeQ)
}

// Timeout ...
func (group *requestGroup) Timeout() {
	group.timeout = true
	group.cacheWrapper.setError(ErrBrokerUndefined)
	group.Finish()
}

// Cache ...
func (group *requestGroup) Cache() Cache {
	return group.cacheWrapper
}
