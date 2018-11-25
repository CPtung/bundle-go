package tag

import (
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/MOXA-ISD/taghub/src/modules"
)

/* ********************************************************* */
/* ******************* Routing Functions ******************* */
/* ********************************************************* */
func GetTags(request modules.RequestData) (int, interface{}) {
    if _type, ok := request.Param["type"]; ok {
        if _type == "all" {
            fmt.Println("get all tags")
        } else if _type == "fieldbus" {
            fmt.Println("get fieldbus tags")
        } else if _type == "system" {
            fmt.Println("get system tags")
        } else if _type == "virtual" {
            fmt.Println("get virtual tags")
        } else {
            fmt.Printf("get type: %v", _type)
        }
        return 200, nil
    }
    return 400, map[string]string{"error": "no supported type"}
}

/* ********************************************************* */
/* ******************* Basis Session ********************** */
/* ********************************************************* */

var configPath string = "./data/tag.json"

type Bundle struct {
    modules.Config
    modules.BundleBase
}

func New(ctx *gin.RouterGroup) *Bundle {
    self := Bundle{}
    // Init Config
    self.Load(configPath)
    // Init Bundle
    self.Index(ctx)
    // Return Instance
    return &self
}

func (self *Bundle) Index(router *gin.RouterGroup) {
    // setup the mapping from route to result handler
    self.GenEndpointHandler()
    self.SetEndpointHandler(router, modules.CRUD_GET, "tags/:type", GetTags)
}
