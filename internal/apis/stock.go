package apis

import (
	"net/http"
	"stock/pkg/spiders"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TrendRequest struct {
	Day        int    `json:"day" form:"day" binding:"lte=2"`
	Code       string `json:"code" form:"code" binding:"required"`
	ShowBefore bool   `json:"show_before" form:"show_before"`
}

func (c *Controller) Trend(ctx *gin.Context) {
	params := new(TrendRequest)
	if err := ctx.Bind(params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  err.Error(),
		})
		return
	}
	trends, err := c.service.Trend(params.Code, params.Day, params.ShowBefore)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":        params.Code,
			"day":         params.Day,
			"show_before": params.ShowBefore,
		}).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "service internal error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"list": trends,
	})
}

type SearchRequest struct {
	Key string `json:"key" form:"key"`
}

func (c *Controller) Search(ctx *gin.Context) {
	params := new(SearchRequest)
	if err := ctx.Bind(params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  err.Error(),
		})
		return
	}
	stocks, err := c.service.Search(params.Key)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code": params.Key,
		}).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "service internal error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"list": stocks,
	})
}

type KLineRequest struct {
	Code      string       `json:"code" form:"code" binding:"required"`
	Type      spiders.Type `json:"type" form:"type"`
	StartTime time.Time    `json:"start_time" form:"start_time" binding:"required" time_format:"2006-01-02 15:04:05"`
	EndTime   time.Time    `json:"end_time" form:"end_time" time_format:"2006-01-02 15:04:05"`
}

func (c *Controller) KLine(ctx *gin.Context) {
	params := new(KLineRequest)
	if err := ctx.Bind(params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  err.Error(),
		})
		return
	}
	if params.Type == "" {
		params.Type = spiders.OneHour
	}
	if params.EndTime.IsZero() {
		params.EndTime = time.Now()
	}
	kline, err := c.service.KLine(params.Code, params.Type, params.StartTime, params.EndTime)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":       params.Code,
			"type":       params.Type,
			"start_time": params.StartTime,
			"end_time":   params.EndTime,
		}).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "service internal error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"data": kline,
	})
}

type StockRequest struct {
	Code string `json:"code" form:"code"`
}

func (c *Controller) Stock(ctx *gin.Context) {
	params := new(StockRequest)
	if err := ctx.Bind(params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  err.Error(),
		})
		return
	}
	stock, err := c.service.Stock(params.Code)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code": params.Code,
		}).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "service internal error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"data": stock,
	})
}

type MultiStockRequest struct {
	Codes []string `json:"codes" form:"codes[]"`
}

func (c *Controller) MultiStock(ctx *gin.Context) {
	params := new(MultiStockRequest)
	if err := ctx.Bind(params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  err.Error(),
		})
		return
	}
	stocks, err := c.service.MultiStock(params.Codes)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code": params.Codes,
		}).Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "service internal error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"list": stocks,
	})
}
