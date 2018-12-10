package micore

import (
    "os"
    "log"
    "io/ioutil"
    "path/filepath"
)

type H map[string]interface{}

type Config struct {
    path string
    data []byte
}

func isExist(path string) bool {
    file, err := os.Open(path)
    defer file.Close()
    if err != nil {
        return false
    }
    return true
}

func (self *Config)Load(path string) []byte {
    self.path = path
    if ok := isExist(self.path); !ok {
        if err := os.MkdirAll(filepath.Dir(self.path), os.ModePerm); err != nil {
            os.RemoveAll(filepath.Dir(self.path))
            log.Fatalf("failed to create file path: %s", err)
            return self.data
	}

	file, err := os.Create(path)
        if err != nil {
            log.Fatalf("failed creating file: %s", err)
        }
        defer file.Close()

        len, err := file.Write([]byte("{}"))
        if err != nil && len > 0 {
            log.Fatalf("failed writing to file: %s", err)
        } else {
            self.data = []byte("{}")
        }

    } else {
        if bytes, err := ioutil.ReadFile(self.path); err == nil {
            self.data = bytes
        }
    }
    return self.data
}

func (self *Config)LoadWithDefault(path string, pattern []byte) []byte {
    self.path = path
    if ok := isExist(self.path); !ok && pattern != nil {
        if err := os.MkdirAll(filepath.Dir(self.path), os.ModePerm); err != nil {
            os.RemoveAll(filepath.Dir(self.path))
            log.Fatalf("failed to create file path: %s", err)
            return self.data
	}

	file, err := os.Create(path)
        if err != nil {
            log.Fatalf("failed creating file: %s", err)
        }
        defer file.Close()

        len, err := file.Write(pattern)
        if err != nil && len > 0 {
            log.Fatalf("failed writing to file: %s", err)
        } else {
            self.data = pattern
        }

    } else {
        if bytes, err := ioutil.ReadFile(self.path); err == nil {
            self.data = bytes
        }
    }
    return self.data
}

func (self *Config)GetAll() []byte {
    return self.data
}

func (self *Config)ReloadAll() []byte {
    if bytes, err := ioutil.ReadFile(self.path); err == nil {
        self.data = bytes
    }
    return self.data
}

func (self *Config)SetAll(data []byte) []byte {
    err := ioutil.WriteFile(self.path, data, os.ModeAppend)
    if (err != nil) {
        return nil
    } else {
        self.data = data
    }
    return self.data
}
