package water

import "github.com/gin-gonic/gin"

type Endpoint func(ctx *gin.Context, req interface{}) (interface{}, error)
