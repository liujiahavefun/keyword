package main

import (
    "flag"
    "fmt"
    "github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
    c.Keys["jsonData"] = gin.H{
        "errno": 0,
        "errmsg": "ok",
        "data": "no fucking data",
    }
}

func main() {
    fmt.Println("starting godlp!")
    flag.Parse()

    gin.SetMode(gin.DebugMode)
    router := gin.Default()

    if gin.Mode() == gin.DebugMode {
        router.GET("/ping", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "message": "fuck off!",
            })
        })
    }

    apiGroup := router.Group("/api")
    {
        apiGroup.GET("/hello", Decorator(hello, WithLogger, WithPrepareEnv, WithJsonRender))

        kwGroup := apiGroup.Group("/dlp")
        {
            kwGroup.GET("/kw", Decorator(keywordMatchHandler, WithLogger, WithPrepareEnv, WithJsonRender))
            kwGroup.POST("/kw", Decorator(keywordMatchHandler, WithLogger, WithPrepareEnv, WithJsonRender))
            kwGroup.GET("/fuzzy", Decorator(fuzzyMatchHandler, WithLogger, WithPrepareEnv, WithJsonRender))
            kwGroup.POST("/fuzzy", Decorator(fuzzyMatchHandler, WithLogger, WithPrepareEnv, WithJsonRender))
            kwGroup.POST("/file", Decorator(fileKeywordMatchHandler, WithLogger, WithPrepareEnv, WithJsonRender))
        }
    }

    fmt.Println("godlp started!")
    router.Run("0.0.0.0:9100")
}
