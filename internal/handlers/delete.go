package handlers

import (
	"net/http"
	"strconv"

	"github.com/wb-go/wbf/ginext"
)

func (r *Router) deleteTransactionHandler(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}
	err = r.tCRUDer.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
