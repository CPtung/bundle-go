package resman

import (
    "log"
    "time"
    "sync"
    "regexp"
    "io/ioutil"
    "encoding/json"

    "github.com/go-redis/redis"
)

const (
    SIG_EXIT        = 1
    PREFIX_KEYSPACE = "__keyspace@0__:"
    PREFIX_TAG      = "taghub:monitor:"
    CONFIG_PATH     = "./resource.json"
)

type ResEventHook func(string, string)

type ResourceCfg struct {
    RedisHost   string      `json:redishost`
    RedisPort   string      `json:redisport`
    Resources   []string    `json:resources`
}

type ResourceClient struct {
    chCancel    chan int
    wgCancel    sync.WaitGroup
    rdsDB       *redis.Client
    rdsPubSub   *redis.PubSub
    config      *ResourceCfg
}

func (self *ResourceClient)LoadConfig() {
    if bytes, err := ioutil.ReadFile(CONFIG_PATH); err == nil {
        self.config = &ResourceCfg{"", "6379", []string{}}
        if err := json.Unmarshal(bytes, self.config); err != nil {
            log.Printf("%v\n", err)
        } else {
            log.Printf("Host: %v\n", self.config.RedisHost)
            log.Printf("Port: %v\n", self.config.RedisPort)
            for _, r := range self.config.Resources {
                log.Printf("Resource: %v\n", r)
            }
        }
    } else {
        log.Println("config path not found")
    }
}

func (self *ResourceClient)Initial() {
    // Load Resource Config
    self.LoadConfig()

    // Init Redis Client
    self.chCancel = make(chan int)
    self.wgCancel = sync.WaitGroup{}
    self.rdsDB = redis.NewClient(&redis.Options{
        Addr:         self.config.RedisHost + ":" + self.config.RedisPort,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        PoolSize:     2,
        PoolTimeout:  30 * time.Second,
    })

    self.rdsPubSub = self.rdsDB.PSubscribe(PREFIX_KEYSPACE + PREFIX_TAG + "*")
}

func (self *ResourceClient)Subscribe(hook ResEventHook) {
    if _, err := self.rdsPubSub.Receive(); err != nil {
        log.Printf("Create pubsub receive failed(%v)\n", err)
        self.rdsPubSub = nil
        return
    }
    ch := self.rdsPubSub.Channel()
    for true {
        select {
        case <-self.chCancel:
            self.rdsPubSub.Close()
            self.wgCancel.Done()
            break
        case msg, ok := <-ch:
            if ok {
                re, _ := regexp.Compile("__keyspace@[0-9]__:" + PREFIX_TAG + "(.*:.*)")
                key := re.ReplaceAllString(msg.Channel, "$1")
                for _, t := range self.config.Resources {
                    if t == key {
                        hook(key, msg.Payload)
                        break
                    }
                }
            }
        }
    }
}

func (self *ResourceClient)Close() {
    if self.rdsPubSub != nil {
        self.wgCancel.Add(1)
        self.chCancel<-SIG_EXIT
        self.wgCancel.Wait()
    }
}
