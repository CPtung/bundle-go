package internal

import (
    "os"
    "fmt"
    "log"
    "net/http"
    "encoding/json"

    "github.com/MOXA-ISD/micore/pkg"
)

const (
    // apicli settings
    APICLI = "apicli"
    TAG_CONFIG_DEFAULT = "{\"tagList\": []}"
    PAYLOAD_PATH = "/tmp/payload.json"

    // fieldbus settings
    FIELDBUS_CONFIG_DEFAULT = "{}"
    DEVICE_CONFIG_PATH = "/var/tagservice/conf.d/fieldbus/devices.json"
    FIELDBUS_CONFIG_PATH = "/var/tagservice/conf.d/fieldbus/protocols.json"
    FIELDBUS_TAG_PATH = "/var/tagservice/conf.d/fieldbus/fieldbusTags.json"
    PROTOCOL_CONFIG_DEFAULT = "{\"protocolList\":[], \"sourceList\": {}}"
)


type ProtocolMgmt struct {
    micore.Config
}

type DeviceMgmt struct {
    micore.Config
}

type FieldbusClient struct {
    micore.Config
    deviceMgmt *DeviceMgmt
    protocolMgmt *ProtocolMgmt
    payloadPath string
}

func ContainOf(item interface{}, slice interface{}) (error, bool) {
    list := reflect.ValueOf(slice)
    switch reflect.TypeOf(list).Kind()a {
        case reflect.Array, reflect.Slice:
            for idx := 0; idx < list.Len(); ++idx {
                if list.Index(idx).Interface() == item {
                    return nil, true
                }
            }
        case reflect.Map:
            if list.MapIndex(reflect.ValueOf(item)).IsValid() {
                return nil, true
            }
    }
    return errors.New("item not found"), false
}

func NewFieldbusClient(rootPath string) *FieldbusClient{
    // Init fieldbus api & tag config
    api := &FieldbusClient{}
    tagFilePath := fmt.Sprintf("%s%s", rootPath, FIELDBUS_TAG_PATH)
    api.LoadWithDefault(tagFilePath, []byte(FIELDBUS_CONFIG_DEFAULT))

    // Init fieldbus protocol config
    api.protocolMgmt = &ProtocolMgmt{}
    configPath := fmt.Sprintf("%s%s", rootPath, FIELDBUS_CONFIG_PATH)
    api.protocolMgmt.LoadWithDefault(configPath, []byte(PROTOCOL_CONFIG_DEFAULT))

    // Init fieldbus device list
    api.deviceMgmt = &DeviceMgmt{}
    devicePath := fmt.Sprintf("%s%s", rootPath, DEVICE_CONFIG_PATH)
    api.deviceMgmt.LoadWithDefault(devicePath, []byte(FIELDBUS_CONFIG_DEFAULT))

    api.payloadPath = fmt.Sprintf("%s%s", rootPath, PAYLOAD_PATH)

    return api
}

func (self *FieldbusClient) WritePayload(data []byte) {
    file, err := os.Create(self.payloadPath)
    if err != nil {
        log.Printf("failed to create payload file: %s", err)
    }
    defer file.Close()

    len, err := file.Write(data)
    if err != nil && len > 0 {
        log.Printf("failed to write payload file: %s", err)
    }
}

func (self *FieldbusClient) GetProtocols() (int, interface{}) {
    var config FieldbusConfig
    if err := json.Unmarshal(self.protocolMgmt.GetAll(), &config); err != nil {
    return http.StatusBadRequest, micore.H{"message": "failed to get protocol list"}
    }
    return http.StatusOK, config.PROTOCOLLIST
}

func (self *FieldbusClient) GetHostName(protocol string) interface{} {
    var config FieldbusConfig
    if err := json.Unmarshal(self.protocolMgmt.GetAll(), &config); err != nil {
        log.Printf("warning: (%v) got failure on parsing protocol list\n", protocol)
        return false
    }
    for _, value := range config.PROTOCOLLIST {
        if value == protocol {
            return config.SOURCELIST[value].SOURCENAME
        }
    }
    return nil
}

func (self *FieldbusClient) Invoke(action string, protocol string) (int, interface{}) {
    hostName := self.GetHostName(protocol)
    if hostName == nil {
        log.Printf("warning: protocol(%v) is not in the supporting list\n", protocol)
        return http.StatusBadRequest, micore.H{"warning": "protocol is not int the supporting list"}
    }
    cmd := fmt.Sprintf("%v --host %v %v", "fbctlcli", hostName, action)
    //log.Println(cmd)
    status, output := micore.Exec(cmd)
    var response micore.H
    if status == 200 {
        if err := json.Unmarshal([]byte(output), &response); err != nil {
            log.Println("response payload is not a json format")
            return status, string(output)
        }
    } else {
        log.Printf("output: %v\n", output)
        response = micore.H{"message": output}
    }
    return status, response
}

func (self *FieldbusClient) StartFieldbusControl(protocol string) (int, interface{}) {
    action := " --start"
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) StopFieldbusControl(protocol string) (int, interface{}) {
    action := " --stop"
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) RestartFieldbusControl(protocol string) (int, interface{}) {
    action := " --restart"
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) UpdateTagList(protocol string) (int, interface{}) {
    action := fmt.Sprintf(" -p %v --tag_list", protocol)
    status, output := self.Invoke(action, protocol)
    if status == http.StatusOK {
        self.ReloadAll()
    }
    return status, output
}

func (self *FieldbusClient) TmpList(protocol string) (int, interface{}) {
    action := fmt.Sprintf(" -p %v --tmp_list", protocol)
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) TmpGet(protocol string, name string) (int, interface{}) {
    action := fmt.Sprintf(" -p %v --tmp_get=%v", protocol, name)
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) TmpAdd(protocol string, data []byte) (int, interface{}) {
    self.WritePayload(data)
    action := fmt.Sprintf(" -p %v --tmp_add=%v", protocol, self.payloadPath)
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) TmpEdit(protocol string, data []byte) (int, interface{}) {
    self.WritePayload(data)
    action := fmt.Sprintf(" -p %v --tmp_edit=%v", protocol, self.payloadPath)
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) TmpRemove(protocol string, name string) (int, interface{}) {
    remove := []byte(fmt.Sprintf("{\"templateName\": \"%v\"}", name))
    self.WritePayload(remove)
    action := fmt.Sprintf(" -p %v --tmp_remove=%v", protocol, self.payloadPath)
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) DeviceAdd(protocol string, data []byte) (int, interface{}) {
    self.WritePayload(data)
    action := fmt.Sprintf(" -p %v --device_add=%v", protocol, self.payloadPath)
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) DeviceEdit(protocol string, data []byte) (int, interface{}) {
    edit := []byte(fmt.Sprintf("{\"deviceList\": [%v]}", string(data)))
    self.WritePayload(edit)
    action := fmt.Sprintf(" -p %v --device_edit=%v", protocol, self.payloadPath)
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) DeviceRemove(protocol string, data []byte) (int, interface{}) {
    remove := []byte(fmt.Sprintf("{\"deviceList\": [%v]}", string(data)))
    self.WritePayload(remove)
    action := fmt.Sprintf(" -p %v --device_remove=%v", protocol, self.payloadPath)
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) MultiDeviceRemove(protocol string, data []byte) (int, interface{}) {
    self.WritePayload(data)
    action := fmt.Sprintf(" -p %v --device_remove=%v", protocol, self.payloadPath)
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) DeviceUpdate(protocol string) (int, interface{}) {
    action := fmt.Sprintf(" -p %v --device_list", protocol)
    status, output := self.Invoke(action, protocol)
    if status == http.StatusOK {
        self.deviceMgmt.ReloadAll()
    }
    return status, output
}

func (self *FieldbusClient) DeviceList(protocol string) (int, interface{}) {
    deviceList := self.deviceMgmt.GetAll()
    var response map[string][]interface{}
    if err := json.Unmarshal(deviceList, &response); err != nil {
        return http.StatusBadRequest, micore.H{"message": "Load device list failed"}
    }
    list := make([]interface{}, 0)
    for key, value := range response {
        if key == protocol || key == "all" {
            list = append(list, value...)
        }
    }
    return http.StatusOK, micore.H{"deviceList": list}
}

func (self *FieldbusClient) TagStatus(protocol string) (int , interface{}) {
    action := fmt.Sprintf(" -p %v --status", protocol)
    return self.Invoke(action, protocol)
}

func (self *FieldbusClient) TagList() (int, interface{}) {
    tagList := self.GetAll()
    var response map[string][]interface{}
    if err := json.Unmarshal(tagList, &response); err != nil {
        return http.StatusBadRequest, micore.H{"message": "Load tag list failed"}
    }

    list := make([]interface{}, 0)
    for _, value := range response {
        list = append(list, value...)
    }
    return http.StatusOK, micore.H{"tagList": list}
}

func (self *FieldbusClient) SetProtocolConfig(data []byte) (int, interface{}) {
    var enrollMessage map[string]string
    if err := json.Unmarshal(data, &enrollMessage); err != nil {
        return http.StatusBadRequest, `"message": "enroll message format error"`
    }
    protocolName, pfound := enrollMessage["protocolName"]
    serviceName, sfound := enrollMessage["serviceName"]
    hostName, hfound := enrollMessage["hostName"]
    if !pfound || !sfound || !hfound {
        return http.StatusBadRequest, `"message": "enroll message content error"`
    }

    var config FieldbusConfig
    if err := json.Unmarshal(self.protocolMgmt.GetAll(), &config); err != nil {
        return http.StatusBadRequest, micore.H{"message": "failed to get protocol config"}
    }

    if _, ok := ContainOf(protocolName, config.PROTOCOLLIST); !ok {
        return http.StatusBadRequest, micore.H{"message": "protocol has been added"}
    }

    config.PROTOCOLLIST = append(config.PROTOCOLLIST, protocolName)
    config.SOURCELIST[protocolName] = Source{ fmt.Sprintf("/%v/control", serviceName), hostName }
    bytes, err := json.Marshal(config)
    if err != nil {
        return http.StatusBadRequest, `"message": "enroll protocol error"`
    }
    self.protocolMgmt.SetAll(bytes)

    return self.DeviceUpdate(protocolName)
}
