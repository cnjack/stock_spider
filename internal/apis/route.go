package apis

import (
	"net/http"
	"stock/internal/services"
	"stock/pkg/spiders"

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
	service := services.NewService(&spiders.EastMoneyProvider{})
	ctl := NewController(service)

	router.GET("health", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})

	gRouter := router.Group("/api")

	gRouter.GET("trend", ctl.Trend)

	gRouter.GET("kline", ctl.KLine)
	gRouter.GET("search", ctl.Search)
	gRouter.GET("stock", ctl.Stock)
	gRouter.GET("multi_stock", ctl.MultiStock)

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
