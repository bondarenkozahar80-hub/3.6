package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

func (r *Router) groupByHandler(c *gin.Context, groupFunc func(context.Context, time.Time, time.Time) (map[string]int64, error)) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	var err error
	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid from date"})
			return
		}
	}
	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid to date"})
			return
		}
	}

	resultRaw, err := groupFunc(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	result := make(map[string]float64)
	for k, v := range resultRaw {
		result[k] = float64(v) / 100.0
	}

	c.JSON(http.StatusOK, result)

}

func (r *Router) groupByDayHandler(c *gin.Context) {
	r.groupByHandler(c, r.tAnalyticsGetter.GroupByDay)
}

func (r *Router) groupByWeekHandler(c *gin.Context) {
	r.groupByHandler(c, r.tAnalyticsGetter.GroupByWeek)
}

func (r *Router) groupByMonthHandler(c *gin.Context) {
	r.groupByHandler(c, r.tAnalyticsGetter.GroupByMonth)
}

func (r *Router) groupByCategoryHandler(c *gin.Context) {
	r.groupByHandler(c, r.tAnalyticsGetter.GroupByCategory)
}
