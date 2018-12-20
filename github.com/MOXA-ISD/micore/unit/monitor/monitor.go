package monitor

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "strconv"

    "github.com/MOXA-ISD/micore/pkg"
    "github.com/MOXA-ISD/micore/unit/monitor/taghub"
)

type Tag struct {
    source  string  `json:"srcName"`
    tag     string  `json:"tagName"`
}

type MonitorRequest struct {
    id      string  `json:"id"`
    list    []Tag   `json:"list"`
}

type Stat struct {
    SOURCENAME      string      `json:"srcName"`
    TAGNAME         string      `json:"tagName"`
    TIMESTAMP       int64       `json:"ts"`
    VALUETYPE       string      `json:"dataType"`
    VALUE           interface{} `json:"dataValue"`
}

func GetMs(strMs []string) int64 {
    var ms int64 = -1
    if len(strMs) > 0 {
        if num, err := strconv.ParseInt(strMs[0], 10, 64); err != nil {
            ms = int64(num)
        }
    }
    return ms
}

func (self *Monitor)PostMonitorEnable(request micore.RequestData) (int, interface{}) {
    var r MonitorRequest
    if err := json.Unmarshal(request.Body, &r); err != nil {
        return http.StatusBadRequest, err.Error()
    }
    if self.hub == nil {
        if self.hub = taghub.NewContext(); self.hub == nil {
            return http.StatusInternalServerError, "Taghub server not ready."
        }
    }
    for _, t := range r.list {
        if rc := self.hub.EnableMonitor(t.source, t.tag); rc != 0 {
            return http.StatusBadRequest, fmt.Sprintf("Enable monitor (%v, %v) failed", t.source, t.tag)
        } else {
            GetMonitorDB().Add(fmt.Sprintf("%v:%v", t.source, t.tag), 0)
        }
    }
    return http.StatusOK, nil
}

func (self *Monitor)PostMonitorDisable(request micore.RequestData) (int, interface{}) {
    var r MonitorRequest
    if err := json.Unmarshal(request.Body, &r); err != nil {
        return http.StatusBadRequest, err.Error()
    }
    if self.hub == nil {
        if self.hub = taghub.NewContext(); self.hub == nil {
            return http.StatusInternalServerError, `"message":"Taghub server not ready."`
        }
    }
    for _, t := range r.list {
        if rc := self.hub.DisableMonitor(t.source, t.tag); rc == 0 {
            return http.StatusBadRequest, fmt.Sprintf("Disable monitor (%v, %v) failed", t.source, t.tag)
        } else {
            GetMonitorDB().Del(fmt.Sprintf("%v:%v", t.source, t.tag))
        }
    }
    return http.StatusOK, nil
}

func (self *Monitor)GetMonitorStats(request micore.RequestData) (int, interface{}) {
    if self.hub == nil {
        if self.hub = taghub.NewContext(); self.hub == nil {
            return http.StatusInternalServerError, `"message":"Taghub server not ready."`
        }
    }

    if !self.hub.IsReady() {
        return http.StatusBadRequest, `"message":"Taghub server not ready."`
    }

    var buffer []Stat
    // Get the name of source whichs for searching tags
    sourceName := request.Param["source"]
    searchMs := GetMs(request.Query["ms"])
    _, foundTag := request.Query["tag"]

    // if no query tags, searching all tags by the given source name
    if !foundTag {
        result := self.hub.ReadBySource(sourceName)
        if result != "" {
            if err := json.Unmarshal([]byte(result), &buffer); err != nil {
                return http.StatusBadRequest, err.Error()
            }
            return http.StatusOK, buffer
        }
    // if a valid searching time period is given(1 ~ 86400000 ms), doing the time searching for each tag
    } else if (searchMs > 0) {
        if (searchMs > 86400000) {
             return http.StatusBadRequest, "Time-Query only supports the range from 1 to 86400000 ms"
        }
        for _, tag := range request.Query["tag"] {
            result := self.hub.ReadByPeriod(sourceName, tag, searchMs)
            if result != "" {
                var States []Stat
                if err := json.Unmarshal([]byte(result), &States); err != nil {
                    log.Println("Parse Stat Json error: ", err)
                    continue
                }
                buffer = append(buffer, States...)
            }
        }
        return http.StatusOK, buffer
    // if no one of above query rules is applied, searching the newest value of each tag
    } else {
        for _, tag := range request.Query["tag"] {
            result := self.hub.Read(sourceName, tag)
            if result != "" {
                var state Stat
                if err := json.Unmarshal([]byte(result), &state); err != nil {
                    log.Println("Parse Stat Json error: ", err)
                    continue
                }
                buffer = append(buffer, state)
            }
        }
        return http.StatusOK, buffer
    }
    return http.StatusBadRequest, "No command has been handled."
}

type Monitor struct {
    micore.CoreRoute
    hnd *EventHandler
    hub *taghub.Context
}

func (self *Monitor) Index() {
    if self.hnd = NewEventHandler(); self.hnd == nil {
        log.Println("Create Monitor Event Handler failed")
    }
    if self.hub = taghub.NewContext(); self.hub == nil {
        log.Println("Create Taghub Context failed")
    }

    // setup the mapping from route to result handler
    self.GenEndpointHandler()
    self.SetEndpointHandler(micore.CRUD_POST, "tags/stats/:source", self.GetMonitorStats)
    self.SetEndpointHandler(micore.CRUD_POST, "tags/monitor/enable",  self.PostMonitorEnable)
    self.SetEndpointHandler(micore.CRUD_POST, "tags/monitor/disable", self.PostMonitorDisable)
}

func (self *Monitor) Stop() {
    if self.hnd != nil {
        self.hnd.Close()
    }
}

func New() *Monitor {
    self := Monitor{}
    // Return Instance
    return &self
}
