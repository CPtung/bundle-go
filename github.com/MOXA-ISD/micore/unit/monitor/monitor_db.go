package monitor

import (
    "log"
    "sync"
)

type TABLE  map[string]int

type Operation interface {
    NumOf(string) int
    AddToDB(string, string, int) int
    DelFromDB(string, string, int) int
}

type MonitorDB struct {
    mtx sync.Mutex
    table TABLE
    Operation
}

var instance *MonitorDB
var once sync.Once

func GetMonitorDB() *MonitorDB {
    once.Do(func() {
        instance = &MonitorDB{}
        instance.mtx = sync.Mutex{}
        instance.table = make(TABLE)
    })
    return instance
}

func (self *MonitorDB)NumOf(key string) int {
    self.mtx.Lock()
    defer self.mtx.Unlock()
    if _, ok := self.table[key]; !ok {
        return 0
    }
    return 1
}

func (self *MonitorDB)Add(key string, ttl int) int {
    self.mtx.Lock()
    defer self.mtx.Unlock()
    if _, ok := self.table[key]; !ok {
        self.table[key] = ttl
    } else if (self.table[key] & ttl) != 0 {
        self.table[key] = 0
    }
    return 0;
}

func (self *MonitorDB)Del(key string) int {
    self.mtx.Lock()
    defer self.mtx.Unlock()
    if _, ok := self.table[key]; !ok {
        log.Printf("Cannot remove monitor(%v) becase key not exists\n", key)
        return -1
    }
    delete(self.table, key)
    return 0
}
