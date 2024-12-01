package handler

import (
	"http-server/danilkovalev/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)


func (h *Handler) CreateIntegration(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var res models.Integration

	if err := c.BindJSON(&res); err != nil {
        return
    }
	h.service.CreateIntegration(id, res)
	c.IndentedJSON(http.StatusCreated, res)
}

func (h *Handler) UpdateIntegration(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var res models.Integration

	if err := c.BindJSON(&res); err != nil {
        return
    }
	err := h.service.UpdateIntegration(id, res)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Integration not found"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Integration updated"})
	}
}

func (h *Handler) GetIntegrations(c *gin.Context) {
	integrations := h.service.GetIntegrations()
	c.IndentedJSON(http.StatusOK, integrations)
}

func (h *Handler) GetIntegration(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	integration, err := h.service.GetIntegration(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Integration not found"})
	} else {
		c.IndentedJSON(http.StatusOK, integration)
	}
}

func (h *Handler) DeleteIntegrations(c *gin.Context) {
	h.service.DeleteIntegrations()
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Integrations deleted"})
}

func (h *Handler) DeleteIntegration(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.service.DeleteIntegration(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Integration not found"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Integration deleted"})
	}
}
