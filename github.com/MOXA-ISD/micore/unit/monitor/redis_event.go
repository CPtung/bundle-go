package monitor

import (
    "log"
    "time"
    "sync"
    "regexp"

    "github.com/go-redis/redis"
)

const (
    REDIS_HOST="redis:6379"
    MONITOR_KEY="taghub:monitor"
)

type Event struct {
    Source      string
    Tag         string
    RedisEvent  string
}

type EventHook func(string, string)

type RdsMgmt struct {
    chExit      chan int
    c           *redis.Client
    ps          *redis.PubSub
    wg          *sync.WaitGroup
}

func NewRdsMgmt() *RdsMgmt {
    m := RdsMgmt{}
    m.c = redis.NewClient(&redis.Options{
                                Addr: REDIS_HOST,
                                PoolSize: 2,
                                PoolTimeout: 30 * time.Second,})
    if _, err := m.c.Ping().Result(); err != nil {
        log.Printf("PingPong redis server error(%v)\n", err)
        defer m.c.Close()
        return nil
    }
    m.wg = &sync.WaitGroup{}
    m.chExit = make(chan int, 1)
    return &m
}

func (self *RdsMgmt) SubRun(hook EventHook) {
    if _, err := self.ps.Receive(); err != nil {
        log.Printf("Create pubsub receive failed(%v)\n", err)
        return
    }
    ch := self.ps.Channel()
    for true {
        select {
        case <-self.chExit:
            self.ps.Close()
            defer self.wg.Done()
            return
        case msg, ok := <-ch:
            if ok {
                re, _ := regexp.Compile("__keyspace@[0-9]__:" + MONITOR_KEY + ":(.*:.*)")
                key := re.ReplaceAllString(msg.Channel, "$1")
                event := msg.Payload
                log.Printf("%v:%v\n", key, event)
                hook(key, event)
            }
        }
    }
}
func (self *RdsMgmt) Subscribe(hook EventHook) {
    if (self.c == nil) {
        return
    }
    self.ps = self.c.PSubscribe("__keyspace@0__:" + MONITOR_KEY + ":*")
    go self.SubRun(hook)
}

func (self *RdsMgmt) Close() {
    self.wg.Add(1)
    self.chExit <- 1
    self.wg.Wait()
}
