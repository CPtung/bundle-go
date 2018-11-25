package group

import (
    "github.com/MOXA-ISD/taghub/src/modules"
    "github.com/MOXA-ISD/taghub/src/tag"
    "github.com/gin-gonic/gin"
)

type BundleGroup struct {
    /* bundles instance */
    TagBundle modules.BundleInterface

    /* internal */
    ginCtx *gin.RouterGroup
}

var instance *BundleGroup = nil
func New(ctx *gin.RouterGroup) *BundleGroup {
    if instance == nil {
        instance = &BundleGroup{ tag.New(ctx) , ctx }
    }
    return instance
}
