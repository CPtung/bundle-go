package main

import (
    "os"
    "log"
    "syscall"
    "os/signal"

    "github.com/MOXA-ISD/micore/sample/internal"
)

func Exit() chan os.Signal {
    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    return quit
}

func Hook(key string, event string) {
    log.Printf("resource (%v) go event (%v)\n", key, event)
}

func main() {
    // Constructor
    c := &resman.ResourceClient{}
    // Initial Resource Client
    c.Initial()
    // Start subscribe
    go c.Subscribe(Hook)

    // Wait Exit Signal
    <-Exit()
    // Destructor
    c.Close()
}
