package apis

import (
	"stock/internal/services"
	"stock/pkg/spiders/eastmoney"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Route(port string) {
	router := gin.Default()
	service := services.NewService(&eastmoney.EastMoneyProvider{})
	ctl := NewController(service)

	router.GET("trend", ctl.Trend)

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
