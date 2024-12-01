package handler

import (
	"http-server/danilkovalev/internal/service"
	"net/http"
	"os"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() {
	router := gin.Default()

	router.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("https://www.amocrm.ru/oauth?client_id=%s&mode=popup", os.Getenv("CLIENT")))
	})

	router.POST("/account", h.CreateAccount)
    router.POST("/integration/:id", h.CreateIntegration)

	router.DELETE("/accounts", h.DeleteAccounts)
	router.DELETE("/account/:id", h.DeleteAccount)
	router.DELETE("/integrations", h.DeleteIntegrations)
	router.DELETE("/integration/:id", h.DeleteIntegration)

	router.PATCH("/account/:id", h.UpdateAccount)
	router.PATCH("/integration/:id", h.UpdateIntegration)

    router.GET("/accounts", h.GetAllAccounts)
	router.GET("/account/:id", h.GetAccount)
    router.GET("/integrations", h.GetIntegrations)
	router.GET("/integration/:id", h.GetIntegration)

	router.GET("/authorization/", h.Authorization)

	router.POST("/key/", h.AddUnisenderKey)
	router.GET("/contact_webhook/", h.ChangeContact)

	router.Run(":8080")
}