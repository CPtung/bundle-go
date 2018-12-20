package monitor

import (
    "log"
)

type EventHandler struct {
    m   *RdsMgmt
}

func (self *EventHandler) UpdateDB(key string, event string) {
    switch event {
        case "set":
            GetMonitorDB().Add(key, 0)
        case "del":
            GetMonitorDB().Del(key)
        case "expired":
            GetMonitorDB().Del(key)
    }
}

func (self *EventHandler) Close() {
    self.m.Close()
}

func NewEventHandler() *EventHandler {
    handler := EventHandler{}
    if handler.m = NewRdsMgmt(); handler.m == nil {
        log.Println("Create redis mgmt failed")
        return nil
    }
    go handler.m.Subscribe(handler.UpdateDB)
    return &handler
}
