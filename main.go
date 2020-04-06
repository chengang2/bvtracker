package main

import (
	"bvtracker/api"
	"bvtracker/dao"
	"bvtracker/g"
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//token := c.GetHeader("authorization")
		token := c.GetHeader("token")
		if token != "bvtracker" {
			g.GinErrorResponse(c, 1, " header must contain bvtracker")
		}
	}
}

// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		//var headerKeys []string
		//for k, _ := range c.Request.Header {
		//	headerKeys = append(headerKeys, k)
		//}
		//headerStr := strings.Join(headerKeys, ", ")
		//if headerStr != "" {
		//	headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		//} else {
		//	headerStr = "access-control-allow-origin, access-control-allow-headers"
		//}
		//fmt.Println("headerStr==",headerStr)
		if origin != "" {
			//下面的都是乱添加的-_-~
			// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", origin)
			//根据 headerStr 加上额外自定义的参数
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, Access-Control-Request-Headers, Origin, Sec-Fetch-Mode, Accept-Language, Connection, Accept, Access-Control-Request-Method, Referer, User-Agent, Accept-Encoding,token,content-type")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			// c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			// c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()
	g.ParseConfig(*cfg)

	dao.Init()
	defer dao.SqlDB.Close()

	r := gin.New()

	r.Use(Cors()) //作用于之后的代码

	login := r.Group("/login")
	login.Use(TokenAuth())
	{
		login.POST("/phone_code", api.Get_phone_code)
		login.POST("/login", api.Login)
	}

	r.POST("/report_data", api.Report_data)
	r.POST("/report_manager", api.Report_manager)
	r.POST("/wetchat_code", api.Wetchat_code)
	r.GET("/wx_login", api.Exchange)
	r.Run(":"+g.Config().ServerPort)

}
