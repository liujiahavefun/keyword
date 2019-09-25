package main

import (
    "runtime"
    "net/http"
    "time"
    "fmt"
    "encoding/json"
    "github.com/gin-gonic/gin"
)

type HttpHandlerDecorator func(gin.HandlerFunc) gin.HandlerFunc

func Decorator(h gin.HandlerFunc, decors... HttpHandlerDecorator) gin.HandlerFunc {
    for i := range decors {
        d := decors[len(decors)-1-i]
        h = d(h)
    }
    return h
}

func WithPrepareEnv(h gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Keys == nil {
            c.Keys = make(map[string]interface{})
        }

        h(c)
    }
}

func WithLogger(h gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func(t time.Time) {
            path := c.Request.URL.Path
            raw := c.Request.URL.RawQuery
            clientIP := c.ClientIP()
            method := c.Request.Method
            statusCode := c.Writer.Status()
            //comment := c.Errors.ByType(0).String()

            if raw != "" {
                path = path + "?" + raw
            }

            retcode := -1
            retmsg := ""
            retdata := ""
            if d, ok := c.Keys["jsonData"]; ok {
                if j, ok := (d.(gin.H)); ok {
                    if code, ok := j["errno"]; ok {
                        retcode = code.(int)
                    }
                    if msg, ok := j["errmsg"]; ok {
                        retmsg = msg.(string)
                    }
                    if data, ok := j["data"]; ok {
                        b, err := json.Marshal(data)
                        if err == nil {
                            retdata = string(b)
                        }
                    }
                }
            }

            fmt.Printf("[GODLP] %v |%3d|%13v|%15s|%-7s %s|%4d|%s|%s \n",
                t.Format("2006/01/02 - 15:04:05"),
                statusCode,
                time.Since(t),
                clientIP,
                method,
                path,
                retcode,
                retmsg,
                retdata,
            )
        }(time.Now())

        h(c)
    }
}

func WithJsonRender(h gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                //logErrorf("process request [%v] met exception: %v ", c.Request.URL, err)
                fmt.Printf("process request [%v] met exception: %v \n", c.Request.URL, err)
                fmt.Println(callStack())
                c.JSON(http.StatusInternalServerError, gin.H{
                    "errno": 1000,
                    "errmsg": "server error",
                    "data": make(map[string]interface{}),
                })
            }
        }()

        h(c)

        if data, exist := c.Keys["jsonData"]; exist {
            c.JSON(http.StatusOK, data)
        }else {
            fmt.Printf("process request [%v] finished, but no response \n", c.Request.URL)
            c.JSON(http.StatusInternalServerError, gin.H{
                "errno": 1000,
                "errmsg": "server error",
                "data": make(map[string]interface{}),
            })
        }
    }
}

func callStack() (msg string) {
    for skip := 0; ; skip++ {
        pc, file, line, ok := runtime.Caller(skip)
        if !ok {
            break
        }
        f := runtime.FuncForPC(pc)
        msg += fmt.Sprintf("frame = %v, file = %v, line = %v, func = %v\n", skip, file, line, f.Name())
    }
    return msg
}
