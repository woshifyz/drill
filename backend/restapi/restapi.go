package restapi

import (
	"drill/config"
	"drill/restapi/middleware"
	_ "drill/static/statik"
	"fmt"
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
)

type RestApi struct {
	port int
}

func NewRestApi(port int) *RestApi {
	return &RestApi{
		port: port,
	}
}

func RenderSuccess(c *gin.Context, resp interface{}) {
    finalResp := gin.H{
        "code":      0,
        "error_msg": "",
        "result":    resp,
    }
    c.JSON(http.StatusOK, finalResp)
}

func (r *RestApi) Start() {
	//gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(cors.Default())
	g.Use(gin.Recovery())
	g.Use(middleware.CatchAllMiddleware)

	v1 := g.Group("/api/v1")
	{
		v1.POST("/util/env", func(context *gin.Context) {
			wsUrl := config.GlobalConfig.PublicWsUrl
			RenderSuccess(context, gin.H{"ws_url": wsUrl})
		})
	}

	fsLoc, _ := fs.New()
	g.StaticFS("/", fsLoc)

	fmt.Printf("start http door at localhost:%d\n", r.port)
	g.Run(fmt.Sprintf("localhost:%d", r.port))
}
