package main

import (
    "fmt"
    "github.com/MOXA-ISD/taghub/src/modules/group"
    "github.com/gin-gonic/gin"
)

var basePath = "/api/v1"

func main() {
    router := gin.Default()
    bundleGroup := group.New(router.Group(basePath))
    if bundleGroup == nil {
        fmt.Println("Failed to create bundle group.")
    } else {
        router.Run(":8080")
    }
}
