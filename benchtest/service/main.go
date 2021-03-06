package main

import (
	"fmt"
	"github.com/importcjj/localcache"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
)

const (
	addr = ":6060"
)

var (
	r = rand.New(rand.NewSource(time.Now().Unix()))
)

var (
	dangerous *localcache.Service
)

func init() {
	dangerous = localcache.NewService(&localcache.Options{})
	dangerous.Register("/getstring", 5*time.Second, cacheGetter)
}

func cacheGetter(key localcache.RequestKey, sink localcache.Sink) error {
	s, err := upstream()
	if err != nil {
		return err
	}

	sink.SetBytes([]byte(s))
	return nil
}

func getHandler(rw http.ResponseWriter, req *http.Request) {
	cache := dangerous.Get(&localcache.Key{
		GetterName: req.URL.Path,
		CacheName:  req.RequestURI,
		Value:      req,
	})

	rw.Write(cache.Bytes())
}

func upstream() (string, error) {
	cost := r.Intn(200)
	time.Sleep(time.Duration(cost) * time.Millisecond)

	if cost > 100 {
		return "", fmt.Errorf("Error=> handle costs %d ms", cost)
	}

	return fmt.Sprintf("OK =>just costs %d ms", cost), nil
}

func main() {
	http.HandleFunc("/getstring", getHandler)
	http.ListenAndServe(addr, nil)
}
