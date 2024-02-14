package handler

import (
	"ewallet/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransferParams struct {
	To     string  `json:"to" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

func (h *Handler) CreateWallet(c *gin.Context) {
	newWallet, err := h.services.Create()
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "can't create wallet internal server error", err)
		return
	}

	c.JSON(http.StatusOK, newWallet)
}

func (h *Handler) Transfer(c *gin.Context) {
	srcId := c.Param("walletId")

	input := TransferParams{}

	err := c.BindJSON(&input)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "incorrect input values", err)
		return
	}

	err = h.services.Transfer(srcId, input.To, input.Amount)
	if err != nil {
		if err == service.ErrOutWalletNoExist {
			h.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
		} else if err == service.ErrInWalletNoExist ||
			err == service.ErrNoBalanceForTrans ||
			err == service.ErrNegativeTransAmount {
			h.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		} else {
			h.ErrorResponse(c, http.StatusInternalServerError, "transfer failed, internal server error", err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successful transfer"})
}

func (h *Handler) History(c *gin.Context) {
	walletId := c.Param("walletId")

	transfers, err := h.services.History(walletId)
	if err != nil {
		if err == service.ErrWalletNoExist {
			h.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
			return
		}

		h.ErrorResponse(c, http.StatusInternalServerError, "can't get history, internal server error", err)
		return
	}

	c.JSON(http.StatusOK, transfers)
}

func (h *Handler) Balance(c *gin.Context) {
	walletId := c.Param("walletId")

	balance, err := h.services.Balance(walletId)

	if err != nil {
		if err == service.ErrWalletNoExist {
			h.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
			return
		}

		h.ErrorResponse(c, http.StatusInternalServerError, "can't get balance, internal server error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": walletId, "balance": balance})
}
