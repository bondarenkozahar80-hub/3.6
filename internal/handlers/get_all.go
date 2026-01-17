package handlers

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

func (r *Router) getAllTransactionsHandler(c *ginext.Context) {
	tList, err := r.tCRUDer.GetAll(c.Request.Context())
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
			"amount":      float64(t.Amount) / 100.0, //конвертация копеек в рубли
			"type":        t.Type,
			"category":    t.Category,
			"event_date":  t.EventDate,
		})
	}
	c.JSON(http.StatusOK, response)
}
