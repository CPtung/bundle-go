package system

import (
    "github.com/MOXA-ISD/micore/pkg"
    "github.com/MOXA-ISD/micore/internal"
)

func (self *System)GetSystemTag(request micore.RequestData) (int, interface{}) {
    return self.client.GetList()
}

type System struct {
    micore.CoreRoute
    client *internal.SystemClient
}

func (self *System) Index() {
    self.client = internal.NewSystemClient("")

    // setup the mapping from route to result handler
    self.GenEndpointHandler()
    self.SetEndpointHandler(micore.CRUD_GET, "tags/system", self.GetSystemTag)
}

func (self *System) Stop() {
}

func New() *System {
    self := System{}
    // Return Instance
    return &self
}
