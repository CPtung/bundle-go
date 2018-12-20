package virtual

import (
    "github.com/MOXA-ISD/micore/pkg"
    "github.com/MOXA-ISD/micore/internal"
)

func (self *Virtual)GetVirtualTag(request micore.RequestData) (int, interface{}) {
    return self.client.GetList()
}

type Virtual struct {
    micore.CoreRoute
    client *internal.VirtualClient
}

func (self *Virtual) Index() {
    self.client = internal.NewVirtualClient("")

    // setup the mapping from route to result handler
    self.GenEndpointHandler()
    self.SetEndpointHandler(micore.CRUD_GET, "tags/virtual", self.GetVirtualTag)
}

func (self *Virtual) Stop() {
}

func New() *Virtual {
    self := Virtual{}
    // Return Instance
    return &self
}
