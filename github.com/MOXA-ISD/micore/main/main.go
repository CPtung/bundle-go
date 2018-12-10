package main

import (
    "log"

    "github.com/MOXA-ISD/micore/pkg"
    "github.com/MOXA-ISD/micore/unit/sample"
)

var basePath = "/api/v1"

func main() {
    miCores := micore.Build(basePath, sample.New())

    if miCores == nil {
        log.Println("Failed to create core services.")
    } else {
        miCores.Run(":8088")
    }
}
