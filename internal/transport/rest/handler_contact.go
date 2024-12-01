package handler

import (
	"fmt"
	"http-server/danilkovalev/internal/models"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)


func (h *Handler) ChangeContact(c *gin.Context) {
    values, err := parseRequestBodyGin(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    changeType, err := determineChangeType(values)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    contact, err := extractContactData(values, changeType)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.service.ChangeContact(contact, changeType); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change contact"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":      contact.AccountID,
        "message": "Contacts changed successfully",
    })
}

func parseRequestBodyGin(c *gin.Context) (url.Values, error) {
    body, err := io.ReadAll(c.Request.Body)
    if err != nil {
        return nil, fmt.Errorf("unable to read request body")
    }

    values, err := url.ParseQuery(string(body))
    if err != nil {
        return nil, fmt.Errorf("unable to parse query")
    }
    return values, nil
}

func determineChangeType(values url.Values) (string, error) {
    switch {
    case values.Get("contacts[update][0][id]") != "":
        return "update", nil
    case values.Get("contacts[add][0][id]") != "":
        return "add", nil
    case values.Get("contacts[delete][0][id]") != "":
        return "delete", nil
    default:
        return "", fmt.Errorf("invalid change type")
    }
}

func extractContactData(values url.Values, changeType string) (models.Contact, error) {
    nameParam := fmt.Sprintf("contacts[%s][0][name]", changeType)
    emailParam := fmt.Sprintf("contacts[%s][0][custom_fields][0][values][0][value]", changeType)
    contactIDParam := fmt.Sprintf("contacts[%s][0][id]", changeType)

    name := values.Get(nameParam)
    email := values.Get(emailParam)
    contactIDStr := values.Get(contactIDParam)

    accountID, err := strconv.Atoi(values.Get("account[id]"))
    if err != nil {
        return models.Contact{}, fmt.Errorf("invalid account_id")
    }

    contactID, err := strconv.Atoi(contactIDStr)
    if err != nil {
        return models.Contact{}, fmt.Errorf("invalid contact_id")
    }

    return models.Contact{
        ID:        contactID,
        Name:      name,
        AccountID: accountID,
        Email:     email,
    }, nil
}
