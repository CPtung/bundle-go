package main

import (
    "log"

    "github.com/MOXA-ISD/micore/pkg"
    "github.com/MOXA-ISD/micore/unit/fieldbus"
    "github.com/MOXA-ISD/micore/unit/system"
    "github.com/MOXA-ISD/micore/unit/virtual"
    "github.com/MOXA-ISD/micore/unit/monitor"
)

var basePath = "/api/v1"

func main() {
    miCores := micore.Build(basePath,
                        fieldbus.New(),
                        system.New(),
                        virtual.New(),
                        monitor.New())

    if miCores == nil {
        log.Println("Failed to create core services.")
    } else {
        miCores.Run(":8086")
    }
}
