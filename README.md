# localcache
Dangerous and High Performance local cache service


## How to use

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
cache := dangerous.Get(&localcache.Key{
    GetterName: "getterName",
    CacheName:  "cachename",
    Value:      context,
})

w.Write(cache.Bytes())
```

    


    