package localcache

import (
	"errors"
	"github.com/importcjj/localcache/log"
	"os"
	"time"
)

// ErrServiceClosed ...
var ErrServiceClosed = errors.New("cache service closed")

// Service ...
type Service struct {
	*Options

	getters map[string]*getter

	getEventsQ        chan *getEvent
	getEventFinishedQ chan *event
	getPendingGroups  map[string]*requestGroup

	updateEventQ        chan *event
	updateEventFinishQ  chan *event
	updatePendingGroups map[string]*requestGroup

	Logger *log.Logger

	closeQ  chan struct{}
	running bool
}

// NewService returns a new cache service.
func NewService(opts *Options) *Service {
	options := newDefaultOptions()
	options.merge(opts)

	logger := log.New(os.Stderr)

	serv := &Service{
		Options: options,
		getters: make(map[string]*getter),

		getEventsQ:        make(chan *getEvent, 100000),
		getEventFinishedQ: make(chan *event, 1000),
		getPendingGroups:  make(map[string]*requestGroup),

		updateEventQ:        make(chan *event, 1000),
		updateEventFinishQ:  make(chan *event, 1000),
		updatePendingGroups: make(map[string]*requestGroup),

		Logger: logger,

		closeQ: make(chan struct{}),
	}

	go serv.loop()

	return serv
}

func (serv *Service) loop() {
	for {
		select {
		case getEvent := <-serv.getEventsQ:
			serv.startGetEvent(getEvent)
		case getfinish := <-serv.getEventFinishedQ:
			serv.finishGetEvent(getfinish)
		case update := <-serv.updateEventQ:
			serv.startUpdateEvent(update)
		case updateFinish := <-serv.updateEventFinishQ:
			serv.finishUpdateEvent(updateFinish)
		case <-serv.closeQ:
			return
		}
	}
}

// Close ...
func (serv *Service) Close() {
	close(serv.closeQ)
}

func (serv *Service) registerUpdateEvent(key RequestKey) {
	serv.Logger.Debug("register update event")
	serv.updateEventQ <- &event{Key: key}
}

func (serv *Service) registerGetEvent(evt *getEvent) {
	select {
	case <-serv.closeQ:
		close(evt.request)
	case serv.getEventsQ <- evt:
	}
}

func (serv *Service) timeoutUpdateEvent(key RequestKey) {
	serv.updateEventFinishQ <- &event{Key: key, Kind: EvtKindTimeout}
}

func (serv *Service) registerUpdateTimeout(key RequestKey, dur time.Duration) *time.Timer {
	return time.AfterFunc(dur, func() {
		serv.updateEventFinishQ <- &event{Key: key, Kind: EvtKindTimeout}
	})
}

func (serv *Service) startUpdateEvent(evt *event) {
	serv.Logger.Debug("start update event")
	getter, ok := serv.getters[evt.Key.GetGetterName()]
	if !ok {
		return
	}

	_, ok = serv.updatePendingGroups[evt.Key.GetCacheName()]
	if ok {
		return
	}

	rg := newRequestGroup()
	serv.updatePendingGroups[evt.Key.GetCacheName()] = rg
	go rg.handleUpdate(serv, evt.Key, getter)
}

func (serv *Service) finishUpdateEvent(evt *event) {
	rg, ok := serv.updatePendingGroups[evt.Key.GetCacheName()]
	if !ok {
		panic("update requests group not found")
	}

	switch evt.Kind {
	case EvtKindTimeout:
		rg.Timeout()
	default:
		rg.Finish()
	}

	delete(serv.updatePendingGroups, evt.Key.GetCacheName())
}

// registerGetTimeout ...
func (serv *Service) registerGetTimeout(key RequestKey, dur time.Duration) *time.Timer {
	return time.AfterFunc(dur, func() {
		serv.getEventFinishedQ <- &event{Key: key, Kind: EvtKindTimeout}
	})
}

func (serv *Service) finishGetEvent(evt *event) {
	rg, ok := serv.getPendingGroups[evt.Key.GetCacheName()]
	if !ok {
		panic("get requests group not found")
	}

	switch evt.Kind {
	case EvtKindTimeout:
		rg.Timeout()
	default:
		rg.Finish()
	}

	delete(serv.getPendingGroups, evt.Key.GetCacheName())
}

func (serv *Service) startGetEvent(evt *getEvent) {
	group, ok := serv.getPendingGroups[evt.Key.GetCacheName()]
	if !ok {
		group = newRequestGroup()
		serv.getPendingGroups[evt.Key.GetCacheName()] = group
		go group.HandleGet(serv, evt.Key)
	}

	evt.request <- group
}

// Register ...
func (serv *Service) Register(key string, expire time.Duration, getterFunc GetterFunc) error {
	serv.getters[key] = &getter{
		Func:   getterFunc,
		Expire: expire,
	}
	return nil
}

// Get Cache by key
func (serv *Service) Get(key RequestKey) Cache {
	evt := &getEvent{Key: key, request: make(chan *requestGroup)}
	serv.registerGetEvent(evt)

	req, ok := <-evt.request
	if !ok {
		goto serviceClosed
	}

	select {
	case <-serv.closeQ:
		goto serviceClosed
	case <-req.closeQ:
		return req.Cache()
	}

serviceClosed:
	cache := newCacheWrapper()
	cache.setError(ErrServiceClosed)
	return cache
}
