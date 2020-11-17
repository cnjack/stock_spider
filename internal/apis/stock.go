package apis

import (
	"net/http"

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
		"data": trends,
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
		"data": stocks,
	})
}
