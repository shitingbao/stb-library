package sgin

import (
	"stb-library/app/storage/internal/biz"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type Sgin struct {
	bg  *biz.GreeterUsecase
	log *log.Helper
	g   *gin.Engine
}

// sgin 只作路由对应
func NewSgin(b *biz.GreeterUsecase, logger log.Logger) *gin.Engine {
	ginModel := gin.Default()
	s := &Sgin{
		bg:  b,
		log: log.NewHelper(logger),
		g:   ginModel,
	}
	s.setRoute()
	return s.g
}

func (s *Sgin) setRoute() {
	rg := s.g.Group("/stb").Use()
	{
		rg.GET("/h", s.helloworld)
	}
}

func (s *Sgin) GetGin() *gin.Engine {
	return s.g
}