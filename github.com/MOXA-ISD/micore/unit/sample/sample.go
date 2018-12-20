package sample

import (
    "log"
    "net/http"

    "github.com/MOXA-ISD/micore/pkg"
)

func (self *Sample)GetValue(request micore.RequestData) (int, interface{}){
    return http.StatusOK, micore.H{}
}

var configPath string = "/var/sample/sample.json"

type Sample struct {
    micore.CoreRoute
}

func (self *Sample) Index() {
    // setup the mapping from route to result handler
    self.GenEndpointHandler()
    self.SetEndpointHandler(micore.CRUD_GET, "sample/:name", self.GetValue)
}

func (self *Sample) Stop() {
    log.Println("Sample Stop")
}

func New() *Sample {
    self := &Sample{}
    // Init Config
    self.Load(configPath)
    // Return Instance
    return self
}
