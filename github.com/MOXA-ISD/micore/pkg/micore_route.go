package micore

import (
    "log"
    "strings"
    "net/http"

    "github.com/gin-gonic/gin"
)

type RouteEntry struct {
    RoutePath  string
    Method     uint8
    callbackFunc CallbackFunc
}

func (entry *RouteEntry) Route(c *gin.Context) {
    // Create a Request Data Object
    requestData := RequestData{ entry.RoutePath, entry.Method, nil, nil, 0, nil, 0 }

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
    if code, data := entry.callbackFunc(requestData); code == http.StatusOK {
        c.JSON(http.StatusOK, data)
    } else {
        c.JSON(http.StatusBadRequest, data)
    }
}


//////////////////////////////////////////////////////////////
type EndpointHandler map[string]RouteEntry

type CoreRoute struct {
    Config
    routeGroup  *gin.RouterGroup
    _handler    EndpointHandler
}

func (self *CoreRoute) SetRouteGroup(routeGroup *gin.RouterGroup) {
    self.routeGroup = routeGroup
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

func (self *CoreRoute) GenEndpointHandler() {
    self._handler = EndpointHandler{}
}

func (self *CoreRoute) SetEndpointHandler(_method uint8, _route string, _callback CallbackFunc) {
    entry := RouteEntry{ _route, _method, _callback }
    switch (_method) {
        case CRUD_GET:
            self.routeGroup.GET(_route, entry.Route)
        case CRUD_POST:
            self.routeGroup.POST(_route, entry.Route)
        case CRUD_PUT:
            self.routeGroup.PUT(_route, entry.Route)
        case CRUD_DEL:
            self.routeGroup.DELETE(_route, entry.Route)
    }
    _resource := GetResource(_method, _route)
    self._handler[_resource] = entry
}

func (self *CoreRoute) GetEndpointHandler(_method uint8, _route string) RouteEntry {
    _resource := GetResource(_method, _route)
    _entry := self._handler[_resource]
    return _entry
}
