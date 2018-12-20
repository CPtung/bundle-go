package internal

type Tag struct {
    SRCNAME       string    `json:"srcName"`
    TAGNAME       string    `json:"tagName"`
    DATATYPE      string    `json:"dataType"`
    DATAUNIT      string    `json:"dataUnit"`
    DURATION      int       `json:"duration"`
    ACCESS        string    `json:"access"`
    DESCRIPTION   string    `json:"description"`
}

type TagList struct {
    TAGLIST  []Tag  `json:"tagList"`
}

type SysTagList struct {
    TAGLIST  []Tag  `json:"systemTagList"`
}

type SysProfile struct {
    DATA  SysTagList  `json:"data"`
}

type Source struct {
    SOURCENAME  string  `json:"sourceName"`
    HOSTNAME    string  `json:"hostName"`
}

type FieldbusConfig struct {
    PROTOCOLLIST  []string           `json:"protocolList"`
    SOURCELIST    map[string]Source  `json:"sourceList"`
}

type DeviceConfig struct {
    protocolDevices map[string]interface{}
}
