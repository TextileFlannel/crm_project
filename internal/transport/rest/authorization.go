package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


func (h *Handler) Authorization(c *gin.Context) {
	id := c.Query("client_id")
	code := c.Query("code")
	referer := c.Query("referer")

	err := h.service.Authorization(code, id, referer)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
	}else{
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Successful authorization"})
	}
}
