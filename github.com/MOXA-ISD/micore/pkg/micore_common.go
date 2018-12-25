package micore

type Resp map[string]interface{}

func RespErr(payload... interface{}) Resp {
    resp := make(Resp)
    if len(payload) > 0 {
        resp["message"] = payload[0]
    }
    return resp
}

func RespBody(key string, payload interface{}) Resp {
    resp := make(Resp)
    resp[key] = payload
    return resp
}
