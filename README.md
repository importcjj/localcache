# localcache

[![Build Status](https://travis-ci.org/importcjj/localcache.svg?branch=master)](https://travis-ci.org/importcjj/localcache)
[![GoDoc](https://godoc.org/github.com/importcjj/localcache?status.svg)](https://godoc.org/github.com/importcjj/localcache)

Dangerous but concurrent local cache wrapper service with high performance.


#### How to use?

```golang


// new serivce
dangerous := localcache.NewService(&localcache.Options{})

// Register Upstream
dangerous.Register("getterName", 5*time.Second, func(key RequestKey, sink Sink) error {
     strVal, err := upstream()
     if err != nil {
         return err
    }
    sink.SetBytes([]byte(strVal))
    return nil
})

// Get Cache
cache := dangerous.Get(&Key{
    GetterName: "getterName",
    CacheName:  "cachename",
    Value:      context,
})

w.Write(cache.Bytes())
```

    


    