package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net/http"
)

func CatchAllMiddleware(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			resp := gin.H{
				"code":      500,
				"error_msg": fmt.Sprint(r),
				"result":      nil,
			}
			c.Render(http.StatusOK, render.JSON{Data: resp})
		}
	}()
	c.Next()
}
