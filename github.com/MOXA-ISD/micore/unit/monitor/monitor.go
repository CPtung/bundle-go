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
    Source  string  `json:"srcName"`
    Tag     string  `json:"tagName"`
}

type MonitorRequest struct {
    Id      string  `json:"id"`
    List    []Tag   `json:"list"`
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
    if self.hub == nil {
        if self.hub = taghub.NewContext(); self.hub == nil {
            return http.StatusInternalServerError, micore.RespErr("Taghub server not ready.")
        }
    }
    var r MonitorRequest
    if err := json.Unmarshal(request.Body, &r); err != nil {
        return http.StatusBadRequest, err.Error()
    }
    for _, t := range r.List {
        if err := GetMonitorDB().Add(r.Id, fmt.Sprintf("%v:%v", t.Source, t.Tag), MAX_TTL_TIME); err != nil {
            return http.StatusBadRequest, micore.RespErr(err.Error())
        }
        if rc := self.hub.EnableMonitor(t.Source, t.Tag); rc != 0 {
            return http.StatusBadRequest, fmt.Sprintf("Enable monitor (%v, %v) failed", t.Source, t.Tag)
        }
    }
    return http.StatusOK, nil
}

func (self *Monitor)PostMonitorDisable(request micore.RequestData) (int, interface{}) {
    var r MonitorRequest
    if err := json.Unmarshal(request.Body, &r); err != nil {
        return http.StatusBadRequest, micore.RespErr(err.Error())
    }
    if self.hub == nil {
        if self.hub = taghub.NewContext(); self.hub == nil {
            return http.StatusInternalServerError, micore.RespErr("Taghub server not ready.")
        }
    }
    for _, t := range r.List {
        if err := GetMonitorDB().Del(r.Id, fmt.Sprintf("%v:%v", t.Source, t.Tag)); err != nil {
            return http.StatusBadRequest, micore.RespErr(err.Error())
        }
        if num := GetMonitorDB().NumOf(fmt.Sprintf("%v:%v", t.Source, t.Tag)); num > 0 {
            continue
        } else if rc := self.hub.DisableMonitor(t.Source, t.Tag); rc != 0 {
            return http.StatusBadRequest, fmt.Sprintf("Disable monitor (%v, %v) failed", t.Source, t.Tag)
        } else {
            log.Printf("disable tag....\n")
        }
    }
    return http.StatusOK, nil
}

func (self *Monitor)GetMonitorStats(request micore.RequestData) (int, interface{}) {
    if self.hub == nil {
        if self.hub = taghub.NewContext(); self.hub == nil {
            return http.StatusInternalServerError, micore.RespErr("Taghub server not ready.")
        }
    }

    if !self.hub.IsReady() {
        return http.StatusBadRequest, micore.RespErr("Taghub server not ready.")
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
                return http.StatusBadRequest, micore.RespErr(err.Error())
            }
            return http.StatusOK, buffer
        }
    // if a valid searching time period is given(1 ~ 86400000 ms), doing the time searching for each tag
    } else if (searchMs > 0) {
        if (searchMs > 86400000) {
             return http.StatusBadRequest, micore.RespErr("Time-Query only supports the range from 1 to 86400000 ms")
        }
        for _, tag := range request.Query["tag"] {
            result := self.hub.ReadByPeriod(sourceName, tag, searchMs)
            if result != "" {
                var States []Stat
                if err := json.Unmarshal([]byte(result), &States); err != nil {
                    log.Println("Parse MonitStat JSON error: ", err)
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
                    log.Println("Parse MonitStat JSON error: ", err)
                    continue
                }
                buffer = append(buffer, state)
            }
        }
        return http.StatusOK, buffer
    }
    return http.StatusBadRequest, micore.RespErr("No command has been handled.")
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
    self.SetEndpointHandler(micore.CRUD_GET, "tags/monitor/:source", self.GetMonitorStats)
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
