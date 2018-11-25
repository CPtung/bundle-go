package modules

import (
    "os"
    "log"
    "io/ioutil"

    "github.com/buger/jsonparser"
)

type Config struct {
    path string
    data []byte
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func isExist(path string) bool {
    file, err := os.Open(path)
    defer file.Close()
    check(err)
    return true
}

/*
func InitConfig(path string) *Config {
    if isExist(path) {
        return  &Config{ path, nil }
    }
    return nil
}*/

func (self *Config)Load(path string) []byte {

    if ok := isExist(path); ok {
        bytes, err := ioutil.ReadFile(path)
        check(err)
        self.path = path
        self.data = bytes
    }

    return self.data
}

func (self *Config)Get(keys ...string) interface{} {
    if value, _, _, err := jsonparser.Get(self.data, keys...); err == nil {
        return value
    }
    return nil
}

func (self *Config)Set(setValue []byte, keys ...string) interface{} {
    if value, err := jsonparser.Set(self.data, setValue, keys...); err == nil {
        return value
    }
    return nil
}

func (self *Config)Save() {
    defer func () {
        log.Fatal("Got failure on loading file from path: ")
    }()

    ioutil.WriteFile(self.path, self.data, os.ModeAppend)
}
