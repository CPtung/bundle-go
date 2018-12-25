package fieldbus

import (
    "net/http"

    "github.com/MOXA-ISD/micore/pkg"
    "github.com/MOXA-ISD/micore/internal"
)

func (self *Fieldbus)GetFieldbusTag(request micore.RequestData) (int, interface{}) {
    return self.client.TagList()
}

func (self *Fieldbus)GetTemplate(request micore.RequestData) (int, interface{}) {
    return self.client.TmpGet(request.Param["type"], request.Param["name"])
}

func (self *Fieldbus)PostTemplate(request micore.RequestData) (int, interface{}) {
    return self.client.TmpAdd(request.Param["type"], request.Body)
}

func (self *Fieldbus)PutTemplate(request micore.RequestData) (int, interface{}) {
    status, output := self.client.TmpEdit(request.Param["type"], request.Body)
    if status == http.StatusOK {
        self.client.UpdateTagList(request.Param["type"])
    }
    return status, output
}

func (self *Fieldbus)DeleteTemplate(request micore.RequestData) (int, interface{}) {
    status, output := self.client.TmpRemove(request.Param["type"], request.Param["name"])
    if status == http.StatusOK {
        self.client.UpdateTagList(request.Param["type"])
    }
    return status, output
}

func (self *Fieldbus)GetTemplates(request micore.RequestData) (int, interface{}) {
    return self.client.TmpList(request.Param["type"])
}

func (self *Fieldbus)GetDevices(request micore.RequestData) (int, interface{}) {
    return self.client.DeviceList(request.Param["type"])
}

func (self *Fieldbus)AddDevice(request micore.RequestData) (int, interface{}) {
    status, output := self.client.DeviceAdd(request.Param["type"], request.Body)
    if status == http.StatusOK {
        self.client.UpdateTagList(request.Param["type"])
    }
    return status, output
}

func (self *Fieldbus)EditDevice(request micore.RequestData) (int, interface{}) {
    status, output := self.client.DeviceEdit(request.Param["type"], request.Body)
    if status == http.StatusOK {
        self.client.UpdateTagList(request.Param["type"])
    }
    return status, output
}

func (self *Fieldbus)DeleteDevices(request micore.RequestData) (int, interface{}) {
    ids, ok := request.Query["ids"]
    if !ok {
        return http.StatusBadRequest, micore.RespErr("Device Id not found")
    }
    status, output := self.client.MultiDeviceRemove(request.Param["type"], ids[0])
    if status == http.StatusOK {
        self.client.UpdateTagList(request.Param["type"])
    }
    return status, output
}

func (self *Fieldbus)GetTagStatus(request micore.RequestData) (int, interface{}) {
    return self.client.TagStatus(request.Param["type"])
}

func (self *Fieldbus)PostStart(request micore.RequestData) (int, interface{}) {
    return self.client.StartFieldbusControl(request.Param["type"])
}

func (self *Fieldbus)PostStop(request micore.RequestData) (int, interface{}) {
    return self.client.StopFieldbusControl(request.Param["type"])
}

func (self *Fieldbus)PostRestart(request micore.RequestData) (int, interface{}) {
    return self.client.RestartFieldbusControl(request.Param["type"])
}

func (self *Fieldbus)GetProtocols(request micore.RequestData) (int, interface{}) {
    return self.client.GetProtocols()
}

func (self *Fieldbus)PostProtocolEvent(request micore.RequestData) (int, interface{}) {
    if request.Param["type"] == "protocol" {
        return self.client.SetProtocolConfig(request.Body)
    }
    return http.StatusNotFound, `{"message": "Resource Not Found"}`
}

var configPath string = "/var/tagservice/conf.d/fieldbus/tag.json"

type Fieldbus struct {
    micore.CoreRoute
    client *internal.FieldbusClient
}

func (self *Fieldbus) Index() {
    self.client = internal.NewFieldbusClient("")

    // setup the mapping from route to result handler
    self.GenEndpointHandler()
    self.SetEndpointHandler(micore.CRUD_GET, "tags/fieldbus", self.GetFieldbusTag)
    self.SetEndpointHandler(micore.CRUD_GET, "tags/fieldbus/:type/status", self.GetTagStatus)
    self.SetEndpointHandler(micore.CRUD_POST,"tags/fieldbus/:type/start", self.PostStart)
    self.SetEndpointHandler(micore.CRUD_POST,"tags/fieldbus/:type/stop", self.PostStop)
    self.SetEndpointHandler(micore.CRUD_POST,"tags/fieldbus/:type/restart", self.PostRestart)
    self.SetEndpointHandler(micore.CRUD_GET, "tags/fieldbus/:type/templates", self.GetTemplates)
    self.SetEndpointHandler(micore.CRUD_POST,"tags/fieldbus/:type/templates/:name", self.PostTemplate)
    self.SetEndpointHandler(micore.CRUD_PUT, "tags/fieldbus/:type/templates/:name", self.PutTemplate)
    self.SetEndpointHandler(micore.CRUD_DEL, "tags/fieldbus/:type/templates/:name", self.DeleteTemplate)

    self.SetEndpointHandler(micore.CRUD_POST,"tags/fieldbus/:type/device", self.AddDevice)
    self.SetEndpointHandler(micore.CRUD_PUT, "tags/fieldbus/:type/device", self.EditDevice)
    self.SetEndpointHandler(micore.CRUD_GET, "tags/fieldbus/:type/devices", self.GetDevices)
    self.SetEndpointHandler(micore.CRUD_DEL, "tags/fieldbus/:type/devices", self.DeleteDevices)

    self.SetEndpointHandler(micore.CRUD_POST,"tags/fieldbus/:type/event", self.PostProtocolEvent)
    self.SetEndpointHandler(micore.CRUD_GET, "tags/fieldbusProtocols", self.GetProtocols)
}

func (self *Fieldbus) Stop() {
}

func New() *Fieldbus {
    self := Fieldbus{}
    // Init Config
    self.Load(configPath)
    // Return Instance
    return &self
}
