package main

import (
	"github.com/backend/cache-benchmark/mock"
	"github.com/backend/localcache"
	"github.com/go-redis/redis"
	"net/http"
	"time"
)

var (
	service *localcache.Service
)

var (
	r *redis.Client
)

type redisBroker struct {
	r *redis.Client
}

func (broker *redisBroker) Get(key string) ([]byte, error) {
	return broker.r.Get(key).Bytes()
}

func (broker *redisBroker) Set(key string, b []byte) error {
	return broker.r.Set(key, b, 0).Err()
}

func Prepare() {
	r = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	service = localcache.NewService(&localcache.Options{
		Broker: &redisBroker{r: r},
	})

	service.Register("/getstring", 5*time.Second, func(key localcache.RequestKey, dest localcache.Sink) error {
		req := key.GetValue().(*http.Request)
		t := req.FormValue("t")
		dest.SetBytes([]byte(mock.ReadDataFromUpStream(t)))
		return nil
	})
}

func HandlerFunc(rw http.ResponseWriter, req *http.Request) {
	// fmt.Println(req.URL.Path, req.RequestURI)
	cache := service.Get(&localcache.Key{
		GetterName: req.URL.Path,
		CacheName:  req.RequestURI,
		Value:      req,
	})
	rw.Header().Set("Last-Modified", cache.LastModified())
	rw.Header().Set("ETag", cache.MD5())
	rw.Write(cache.Bytes())
}

func main() {
	Prepare()
	http.HandleFunc("/getstring", HandlerFunc)
	http.ListenAndServe(":6060", nil)
}
