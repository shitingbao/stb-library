package sgin

import (
	"net/http"
	"os"
	"path"
	v1 "stb-library/api/storage/v1"
	"stb-library/app/storage/internal/biz"
	"stb-library/lib/response"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type Sgin struct {
	v1.UnimplementedStorageServer
	user      *biz.UserUseCase
	transform *biz.TransformUseCase
	center    *biz.CentralUseCase
	log       *log.Helper
	g         *gin.Engine
}

func NewGinEngine() *gin.Engine {
	return gin.Default()
}

// sgin 只作路由对应
func NewSgin(ginModel *gin.Engine, u *biz.UserUseCase, c *biz.CentralUseCase, logger log.Logger) *Sgin {
	s := &Sgin{
		user:   u,
		center: c,
		log:    log.NewHelper(logger),
		g:      ginModel,
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	s.setRoute(dir)
	return s
}

func (s *Sgin) setRoute(dir string) {
	s.g.Use(cross)
	rg := s.g.Group("/api")
	{
		rg.POST("/login", s.login)
		rg.GET("/logout", s.logout)
		rg.POST("/register", s.register)

		rg.GET("/userinfo", s.getUserInfo)
		rg.POST("/upload", s.upload)

		rg.GET("/central", s.sayHello)
	}

	s.g.StaticFS("assets", http.Dir(path.Join(dir, "assets")))
}

func (s *Sgin) sayHello(ctx *gin.Context) {
	user := &biz.UserResult{}
	if err := ctx.Bind(user); err != nil {
		response.JsonErr(ctx, err, nil)
		return
	}

	n, err := s.center.SayHello(user.UserName)

	if err != nil {
		response.JsonErr(ctx, err, nil)
		return
	}
	response.JsonOK(ctx, n)
}

func cross(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "127.0.0.1,124.70.156.31,socket1.cn") //设置允许跨域的请求地址
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
}
