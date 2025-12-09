package ugin

import "github.com/gin-gonic/gin"

type Handler interface {
	Routes(r gin.IRouter)
}

type Middleware interface {
	Name() string
	Middleware() gin.HandlerFunc
}
