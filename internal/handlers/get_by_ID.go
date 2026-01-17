package handlers

import (
	"net/http"
	"strconv"

	"github.com/wb-go/wbf/ginext"
)

func (r *Router) getTransactionByIDHandler(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}
	t, err := r.tCRUDer.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}
	if t == nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "transaction not found in DB"})
		return
	}
	c.JSON(http.StatusOK, ginext.H{
		"id":          t.ID,
		"name":        t.Name,
		"description": t.Description,
		"amount":      float64(t.Amount) / 100.0, //конывертация копеек в рубли
		"type":        t.Type,
		"category":    t.Category,
		"event_date":  t.EventDate,
	})

}
