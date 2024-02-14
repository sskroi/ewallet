package handler

import (
	"ewallet/pkg/service"
	"log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func New(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	apiV1Wallet := router.Group("/api/v1/wallet")
	{
		apiV1Wallet.POST("", h.CreateWallet)
		apiV1Wallet.POST("/:walletId/send", h.Transfer)
		apiV1Wallet.GET("/:walletId/history", h.History)
		apiV1Wallet.GET("/:walletId", h.Balance)
	}

	return router
}

func (h *Handler) ErrorResponse(c *gin.Context, statusCode int, msg string, err error) {
	log.Printf("ERROR: %s", err.Error())

	c.AbortWithStatusJSON(statusCode, gin.H{"message": msg})
}
