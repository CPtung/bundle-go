package monitor

import (
    "fmt"
    "sync"
    "errors"

    "github.com/MOXA-ISD/micore/pkg"
)

const (
    MAX_TTL_TIME    = 65536
    STR_LAST_UPDATE = "last_update_ts"
)

type TTL    map[string]interface{}
type TABLE  map[string]TTL          // map[Key]TTL

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
    if _, ok := self.table[key][STR_LAST_UPDATE]; ok {
        return len(self.table[key]) - 1
    }
    return len(self.table[key])
}

func (self *MonitorDB)Add(id string, key string, timeout int) error {
    self.mtx.Lock()
    defer self.mtx.Unlock()
    if _, ok := self.table[key]; !ok {
        ttl := make(TTL)
        ttl[id] = timeout
        self.table[key] = ttl
    } else {
        // Save New User Setting
        if _, ok := self.table[key][id]; !ok {
            self.table[key][id] = timeout
            return nil
        }
        // Check if incoming timeout time is longer than ever
        var new_val int = timeout
        for key, val := range self.table[key] {
            if key != STR_LAST_UPDATE && new_val < val.(int) {
                new_val = val.(int)
            }
        }
        self.table[key][id] = new_val
    }
    self.table[key][STR_LAST_UPDATE] = micore.GetTimeStamp()
    return nil
}

func (self *MonitorDB)Del(id string, key string) error {
    self.mtx.Lock()
    defer self.mtx.Unlock()
    if id == "__super__" {
        delete(self.table, key)
        return nil
    } else if _, ok := self.table[key]; !ok {
        err := fmt.Sprintf("Cannot remove monitor because Key(%v) not exists", key)
        return errors.New(err)
    } else if _, ok := self.table[key][id]; !ok {
        err := fmt.Sprintf("Cannot remove monitor because UserId(%v) not exists", id)
        return errors.New(err)
    }
    delete(self.table[key], id)
    self.table[key][STR_LAST_UPDATE] = micore.GetTimeStamp()

    if len(self.table[key]) == 0 {
        delete(self.table, key)
    }
    return nil
}
