package singleton

type Singleton struct{}

var instance = &Singleton{}

func GetInstance() *Singleton {
    return instance
}
