package micore

import (
    "github.com/gin-gonic/gin"
)

type CoreBase interface {
    SetRouteGroup(routeGroup *gin.RouterGroup)
    Index()
    Stop()
}
