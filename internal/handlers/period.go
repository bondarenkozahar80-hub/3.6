package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

func (r *Router) getTransactionsByPeriodHandler(c *gin.Context) {
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

	tList, err := r.tCRUDer.GetByPeriod(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	var response []ginext.H
	for _, t := range tList {
		response = append(response, ginext.H{
			"id":          t.ID,
			"name":        t.Name,
			"description": t.Description,
			"amount":      float64(t.Amount) / 100.0,
			"type":        t.Type,
			"category":    t.Category,
			"event_date":  t.EventDate,
		})
	}
	c.JSON(http.StatusOK, response)

}
