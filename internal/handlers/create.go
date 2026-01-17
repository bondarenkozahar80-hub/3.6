package handlers

import (
	"net/http"
	"time"

	"github.com/bondarenkozahar80-hub/3.6/internal/model"
	"github.com/wb-go/wbf/ginext"
)

func (r *Router) createTransactionHandler(c *ginext.Context) {
	var tDTO TransactionDTO
	if err := c.ShouldBindJSON(&tDTO); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	amount := int64(tDTO.Amount * 100)

	eventDate, err := time.Parse("2006-01-02", tDTO.EventDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid event_date"})
		return
	}

	t := &model.Transaction{
		Name:        tDTO.Name,
		Description: tDTO.Description,
		Amount:      amount,
		Type:        tDTO.Type,
		Category:    tDTO.Category,
		EventDate:   eventDate.Format("2006-01-02"), // хранится как строка
	}

	id, err := r.tCRUDer.Create(c.Request.Context(), t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ginext.H{"id": id})

}
