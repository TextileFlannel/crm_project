package handler

import (
	"http-server/danilkovalev/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)


func (h *Handler) GetAllAccounts(c *gin.Context) {
	accounts := h.service.GetAllAccounts()
	c.IndentedJSON(http.StatusOK, accounts)
}


func (h *Handler) CreateAccount(c *gin.Context) {
	var res models.Account

	if err := c.BindJSON(&res); err != nil {
        return
    }
	h.service.CreateAccount(res)
    c.IndentedJSON(http.StatusCreated, res)
}


func (h *Handler) GetAccount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	account, err := h.service.GetAccount(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Account not found"})
	} else {
		c.IndentedJSON(http.StatusOK, account)
	}
}


func (h *Handler) DeleteAccounts(c *gin.Context) {
	err := h.service.DeleteAccounts()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Account not found"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Accounts deleted"})
	}
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.service.DeleteAccount(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Account not found"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Account deleted"})
	}
}

func (h *Handler) UpdateAccount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var res models.Account

	if err := c.BindJSON(&res); err != nil {
        return
    }
	err := h.service.UpdateAccount(id, res)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Account not found"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Account updated"})
	}
}


func (h *Handler) AddUnisenderKey(c *gin.Context) {
	unisenderKey := c.PostForm("unisender_key")
	accountIDStr := c.PostForm("account_id")

	accountID, _ := strconv.Atoi(accountIDStr)

	unisender := models.Unisender{
		AccountID: accountID,
		Key:       unisenderKey,
	}

	err := h.service.AddUnisenderKey(unisender.AccountID, unisender.Key)
	
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}else{
		c.IndentedJSON(http.StatusOK, unisender)
	}
}

