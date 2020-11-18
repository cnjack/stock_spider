package apis

import (
	"stock/internal/services"
	"stock/pkg/spiders/eastmoney"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Route(port string) {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	service := services.NewService(&eastmoney.EastMoneyProvider{})
	ctl := NewController(service)
	gRouter := router.Group("/api")

	gRouter.GET("trend", ctl.Trend)
	gRouter.GET("kline", ctl.KLine)
	gRouter.GET("search", ctl.Search)

	if err := router.Run(port); err != nil {
		logrus.Panicln(err)
	}
}

type Controller struct {
	service *services.StockImpl
}

func NewController(service *services.StockImpl) *Controller {
	return &Controller{
		service: service,
	}
}
