package micore

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

