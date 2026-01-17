package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

func (r *Router) getAllSortedHandler(c *gin.Context) {
	sortField := c.Query("sortField")
	order := c.Query("order")

	tList, err := r.tCRUDer.GetAllSorted(c.Request.Context(), sortField, order)
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
