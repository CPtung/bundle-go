package modules

import (
    "log"
    "strings"
    "net/http"

    "github.com/gin-gonic/gin"
)

const (
    CRUD_NONE = 0x0
    CRUD_GET  = 0x1
    CRUD_POST = 0x2
    CRUD_PUT  = 0x4
    CRUD_DEL  = 0x8
)

// Routing Parameters
type Parameters map[string]string
// Routing Query Parameters
type ExtraParams map[string][]string
// A Http Request Object
type RequestData struct {
    RoutePath   string
    Method      uint8
    Body        []byte
    Param       Parameters
    ParamSize   int
    Query       ExtraParams
    QuerySize   int
}

// virtual function to handle http requests for user 
type CallbackFunc func(RequestData) (int, interface{})

type RouteEntry struct {
    RoutePath  string
    Method     uint8
    callbackFunc     CallbackFunc
}

func (entry *RouteEntry) Route(c *gin.Context) {
    // Create a Request Data Object
    requestData := RequestData{ entry.RoutePath, entry.Method, nil, nil, 0, nil, 0}

    // Get Parameters
    if requestData.ParamSize = len(c.Params); requestData.ParamSize > 0 {
        requestData.Param = make(Parameters, requestData.ParamSize)
        for _, param := range c.Params {
            requestData.Param[param.Key] = param.Value
        }
    }

    // Get Query array
    if query := c.Request.URL.Query(); len(query) > 0 {
        requestData.QuerySize = len(query)
        requestData.Query = ExtraParams{}
        for key, value := range query {
            requestData.Query[key] = value
        }
    }

    // Get Request body
    if body, _err := c.GetRawData(); _err == nil {
        requestData.Body = body
    } else {
        log.Println("Failed to parse json body")
    }

    // Get Response
    if code, data := entry.callbackFunc(requestData); code == 200 {
        c.JSON(http.StatusOK, data)
    } else {
        c.JSON(400, data)
    }
}


//////////////////////////////////////////////////////////////
type EndpointHandler map[string]RouteEntry

type BundleBase struct {
    _handler    EndpointHandler
}

type BundleInterface interface {
    Index(router *gin.RouterGroup)
}

func GetResource(_method uint8, _route string) string {
    var resource strings.Builder
    switch (_method) {
        case CRUD_GET:
            resource.WriteString("get")
        case CRUD_POST:
            resource.WriteString("post")
        case CRUD_PUT:
            resource.WriteString("put")
        case CRUD_DEL:
            resource.WriteString("delete")
    }
    resource.WriteString(_route)
    return resource.String()
}
/*
func (self *BundleBase) LoadConfig(path string) {
    if self._config = InitConfig(path); self._config == nil {
        log.Printf("Bundle failed to load config\n")
    }
}*/

func (self *BundleBase) GenEndpointHandler() {
    self._handler = EndpointHandler{}
}

func (self *BundleBase) SetEndpointHandler(router *gin.RouterGroup, _method uint8, _route string, _callback CallbackFunc) {
    entry := RouteEntry{ _route, _method, _callback }
    switch (_method) {
        case CRUD_GET:
            router.GET(_route, entry.Route)
        case CRUD_POST:
            router.POST(_route, entry.Route)
        case CRUD_PUT:
            router.PUT(_route, entry.Route)
        case CRUD_DEL:
            router.DELETE(_route, entry.Route)
    }
    _resource := GetResource(_method, _route)
    self._handler[_resource] = entry
}

func (self *BundleBase) GetEndpointHandler(_method uint8, _route string) RouteEntry {
    _resource := GetResource(_method, _route)
    _entry := self._handler[_resource]
    return _entry
}
