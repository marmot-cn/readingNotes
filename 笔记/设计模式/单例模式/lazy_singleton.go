package singleton

import "sync"

type LazySingleton struct{}

var lazyInstance *LazySingleton
var mu sync.Mutex

func GetLazyInstance() *LazySingleton {
    if lazyInstance == nil {
        mu.Lock()
        defer mu.Unlock()
        
        if lazyInstance == nil {
            lazyInstance = &LazySingleton{}
        }
    }
    return lazyInstance
}
