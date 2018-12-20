package internal

import (
    "fmt"
    "log"
    "bytes"
    "net/http"
    "encoding/json"

    "github.com/tidwall/gjson"
    "github.com/MOXA-ISD/micore/pkg"
)

const (
    SYS_CFG_DEFAULT = "{\"tagList\": []}"
    SYSTEM_PROFILE_ENDPOINT = "/profile"
    SYSTEM_TAG_PATH = "/var/tagservice/conf.d/system/systemTags.json"
)

type SystemClient struct {
    micore.Config
}

func NewSystemClient(rootPath string) *SystemClient{
    api := &SystemClient{}
    configPath := fmt.Sprintf("%s%s", rootPath, SYSTEM_TAG_PATH)
    api.LoadWithDefault(configPath, []byte(SYS_CFG_DEFAULT))
    api.GetDeviceProfile()
    return api
}

func (self *SystemClient) GetDeviceProfile() {
    getProfileCmd := fmt.Sprintf("%v -d get -r %v", APICLI, SYSTEM_PROFILE_ENDPOINT)
    if status, output := micore.Exec(getProfileCmd); status != 200 {
        log.Printf("[system]: get device profile failed (%v)\n", status)
    } else {
        value := gjson.Get(output, "data.systemTagList")
        buffer := new(bytes.Buffer)
        if err := json.Compact(buffer, []byte(value.String())); err != nil {
            log.Println(err)
            buffer.WriteString("[]")
        }
        list := fmt.Sprintf("{\"tagList\": %v}", buffer.String())
        self.SetAll([]byte(list))
    }
}

func (self *SystemClient) GetList() (int, interface{}) {
    self.GetDeviceProfile()
    var tagList TagList
    if err := json.Unmarshal(self.GetAll(), &tagList); err != nil {
        return http.StatusNoContent, http.StatusText(http.StatusNoContent)
    }
    return http.StatusOK, tagList
}
