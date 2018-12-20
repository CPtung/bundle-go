package taghub
/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -ltaghub -ljansson
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <taghub.h>
*/
import "C"
import "log"
import "unsafe"

type Context struct {
    capi *C.taghuber
}

const TAGHUB_CONN_HOST   string = "redis"
const TAGHUB_CONN_PORT      int = 6379
const TAGHUB_CONN_TIMEOUT   int = 3000    // in millisecond

var ctx *Context = nil

func NewContext() *Context {
    if ctx == nil {
        ctx = &Context{};
        if ctx.Init() != 0 {
            return nil
        }
    }
    return ctx
}

func (c *Context) IsReady() bool {
    defer func() {
        log.Println("Taghub server not ready")
    }()
    return  nil != c.capi
}

func (c *Context) Init() int {
    /* taghuber configuration */
    redis_server_domain := C.CString(TAGHUB_CONN_HOST)
    C.taghuber_config(redis_server_domain, 6379, 3000)

    /* taghuber initialization */
    c.capi = C.taghuber_new()
    if c.capi == nil {
        log.Println("Failed to create taghuber instance.")
        return 1
    }
    return 0
}

func (c *Context) Read(source, tag string) string {
    var result string = ""

    if c.capi == nil {
        log.Println("Null Taghuber Instance.")
        return ""
    }

    szSource := C.CString(source)
    szTag := C.CString(tag)
    if buffer := C.taghuber_read(c.capi, szSource, szTag); buffer != nil {
        result = C.GoString(buffer)
        if buffer != nil {
            defer C.free(unsafe.Pointer(buffer)); buffer = nil
        }
    }

    return result
}

func (c *Context) ReadBySource(source string) string {
    var result string = ""

    if c.capi == nil {
        log.Println("Null Taghuber Instance.")
        return ""
    }

    szSource := C.CString(source)
    if buffer := C.taghuber_read_by_source(c.capi, szSource); buffer != nil {
        result = C.GoString(buffer)
        if buffer != nil {
            defer C.free(unsafe.Pointer(buffer)); buffer = nil
        }
    }

    return result
}

func (c *Context) ReadByPeriod(source string, tag string, ms int64) string {
    var result string = ""

    if c == nil {
        log.Println("Context not initialized yet.")
        return ""
    }

    if c.capi == nil {
        log.Println("Null Taghuber Instance.")
        return ""
    }

    szTag := C.CString(tag)
    szSource := C.CString(source)
    if buffer := C.taghuber_read_by_period(c.capi, szSource, szTag, C.longlong(ms)); buffer != nil {
        result = C.GoString(buffer)
        if buffer != nil {
            defer C.free(unsafe.Pointer(buffer)); buffer = nil
        }
    }

    return result
}

func (c *Context) EnableMonitor(source string, tag string) int {

    if c.capi == nil {
        log.Println("Null Taghuber Instance.")
        return -1
    }

    strSource := C.CString(source)
    strTag := C.CString(tag)
    if rc, err := C.taghuber_monitor(c.capi, strSource, strTag, 1, 0); err != nil {
        log.Printf("EnableMonitor caught exception: (%v)\n", err)
        return int(rc)
    }
    return 0
}

func (c *Context) DisableMonitor(source string, tag string) int {

    if c.capi == nil {
        log.Println("Null Taghuber Instance.")
        return -1
    }

    strSource := C.CString(source)
    strTag := C.CString(tag)
    if rc, err := C.taghuber_monitor(c.capi, strSource, strTag, 0, 0); err != nil {
        log.Printf("DisableMonitor caught exception: (%v)\n", err)
        return int(rc)
    }
    return 0
}

func (c *Context) DeInit() {
    if (c != nil && c.capi != nil) {
        C.taghuber_delete(&(c.capi))
    }
}
