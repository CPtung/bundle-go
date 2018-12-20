package micore

import (
    "os"
    "log"
    "syscall"
    "context"
    "os/signal"

    "github.com/gin-gonic/gin"
)

type micore struct {
    /* gin router instance */
    router *gin.Engine
    /* micore instances */
    cores []CoreBase
}

var instance *micore = nil
func Build(basePath string, cores ...CoreBase) *micore {
    if instance == nil {
	instance = &micore{}
        instance.router = gin.Default()
        group := instance.router.Group(basePath)

        for _, core := range cores {
            core.SetRouteGroup(group)
            instance.cores = append(instance.cores, core)
        }
    }
    return instance
}

func(self *micore) Run(addr string) (err error) {
    defer func() {
        log.Println("Activate Cores failed")
    }()

    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    ctx, cancel := context.WithCancel(context.Background())
    go func() {
        for {
            select {
                case <-quit:
                    cancel()
                case <-ctx.Done():
                    self.Stop()
                    os.Exit(0)
            }
        }
    }()
    for _, core := range self.cores {
        core.Index()
    }
    return self.router.Run(addr)
}

func(self *micore) Stop() {
    for _, core := range self.cores {
        core.Stop()
    }
}
