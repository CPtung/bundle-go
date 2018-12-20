package internal

import (
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/MOXA-ISD/micore/pkg"
)

const (
    VIR_CFG_DEFAULT = "{\"tagList\": []}"
    VIRTUAL_TAG_PATH = "/var/tagservice/conf.d/virtual/virtualTags.json"
)

type VirtualClient struct {
    micore.Config
}

func NewVirtualClient(rootPath string) *VirtualClient{
    api := &VirtualClient{}
    configPath := fmt.Sprintf("%s%s", rootPath, VIRTUAL_TAG_PATH)
    api.LoadWithDefault(configPath, []byte(VIR_CFG_DEFAULT))
    return api
}

func (self *VirtualClient) UpdateTag(data []byte) (int, interface{}) {
    var config TagList
    if err := json.Unmarshal(self.GetAll(), &config); err != nil {
        return http.StatusBadRequest, micore.H{"error": "cannot load config"}
    }

    var tag Tag
    if err := json.Unmarshal(data, &tag); err != nil {
        return http.StatusBadRequest, micore.H{"error": "unknown tag format"}
    }

    for i := range config.TAGLIST {
        _tag := config.TAGLIST[i]
        if _tag.SRCNAME == tag.SRCNAME && _tag.TAGNAME == tag.TAGNAME {
            config.TAGLIST[i] = tag
            if bytes, _ := json.Marshal(config); bytes != nil {
                self.SetAll(bytes)
            }
            return http.StatusOK, tag
        }
    }

    config.TAGLIST = append(config.TAGLIST, tag)
    if bytes, _ := json.Marshal(config); bytes != nil {
        self.SetAll(bytes)
    }
    return http.StatusOK, tag
}

func (self *VirtualClient) DeleteTag(data []byte) (int, interface{}) {
    var config TagList
    if err := json.Unmarshal(self.GetAll(), &config); err != nil {
        return http.StatusBadRequest, micore.H{"error": "cannot load config"}
    }

    var tag Tag
    if err := json.Unmarshal(data, &tag); err != nil {
        return http.StatusBadRequest, micore.H{"error": "unknown tag format"}
    }

    for i := range config.TAGLIST {
        _tag := config.TAGLIST[i]
        if _tag.SRCNAME == tag.SRCNAME && _tag.TAGNAME == tag.TAGNAME {
            config.TAGLIST = append(config.TAGLIST[:i], config.TAGLIST[i+1:]...)
            if bytes, _ := json.Marshal(config); bytes != nil {
                self.SetAll(bytes)
            }
            return http.StatusOK, tag
        }
    }

    return http.StatusBadRequest, micore.H{"error": "tag not found"}
}

func (self *VirtualClient) GetList() (int, interface{}) {
    var tagList TagList
    if err := json.Unmarshal(self.GetAll(), &tagList); err != nil {
        return http.StatusNoContent, http.StatusText(http.StatusNoContent)
    }
    return http.StatusOK, tagList
}
