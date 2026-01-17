package handlers

import (
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
)

func (r *Router) getAnalyticsHandler(c *ginext.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from time.Time = time.Time{} // минимальная дата, если не передана
	var to time.Time = time.Now()    // текущая дата, если не передана
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

	sumKopeyki, err := r.tAnalyticsGetter.GetSum(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	avgKopeyki, err := r.tAnalyticsGetter.GetAvg(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	count, err := r.tAnalyticsGetter.GetCount(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	medianKopeyki, err := r.tAnalyticsGetter.GetMedian(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	p90Kopeyki, err := r.tAnalyticsGetter.GetPercentile90(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"sum":          float64(sumKopeyki) / 100.0,
		"avg":          float64(avgKopeyki) / 100.0,
		"count":        count,
		"median":       medianKopeyki / 100.0,
		"percentile90": p90Kopeyki / 100.0,
	})
}
