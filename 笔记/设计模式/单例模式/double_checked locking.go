package singleton

import (
	"sync"
)

type Singleton struct {
	// 你的数据字段
}

var instance *Singleton
var mu sync.Mutex

// GetInstance 使用双重检测来获取单例对象的实例
func GetInstance() *Singleton {
	if instance == nil { // 第一次检查，如果实例不存在，则锁定
		mu.Lock()
		defer mu.Unlock()

		if instance == nil { // 第二次检查，如果实例仍然不存在，则创建
			instance = &Singleton{}
		}
	}
	return instance
}
