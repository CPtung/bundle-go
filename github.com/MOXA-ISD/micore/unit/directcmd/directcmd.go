package directcmd

import (
    "net/http"
    "github.com/MOXA-ISD/micore/pkg"
)

func (self *DirectCmd)GetValue(request micore.RequestData) (int, interface{}){
    return http.StatusOK, micore.RespBody()
}

var configPath string = "/var/tagservice/conf.d/directcmd/directcmd.json"

type DirectCmd struct {
    micore.CoreRoute
}

func (self *DirectCmd) Index() {
    // setup the mapping from route to result handler
    self.GenEndpointHandler()
    self.SetEndpointHandler(micore.CRUD_GET, "sample/:name", self.GetValue)
}

func (self *DirectCmd) Stop() {
}

func New() *DirectCmd {
    self := DirectCmd{}
    // Init Config
    self.Load(configPath)
    // Return Instance
    return &self
}
