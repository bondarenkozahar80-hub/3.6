package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bondarenkozahar80-hub/3.6/internal/model"
	"github.com/wb-go/wbf/ginext"
)

func (r *Router) updateTransactionHandler(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	var tDTO TransactionDTO
	if err := c.ShouldBindJSON(&tDTO); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	// Конвертируем сумму в копейки
	amount := int64(tDTO.Amount * 100)

	// Парсим дату
	eventDate, err := time.Parse("2006-01-02", tDTO.EventDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid event_date"})
		return
	}

	t := &model.Transaction{
		ID:          id,
		Name:        tDTO.Name,
		Description: tDTO.Description,
		Amount:      amount,
		Type:        tDTO.Type,
		Category:    tDTO.Category,
		EventDate:   eventDate.Format("2006-01-02"),
	}

	if err := r.tCRUDer.Update(c.Request.Context(), t); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)

}
