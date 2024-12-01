package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"http-server/danilkovalev/internal/models"
	"io"
	"net/http"
	"os"
	"strings"
)


func (s *AccountService) Authorization(code, id, referer string) error {
	subdomain := strings.Split(referer, ".")[0]

	url := fmt.Sprintf("https://%s.amocrm.ru/oauth2/access_token", subdomain)

	requestBody := models.AuthRequest{
		ClientID:     id,
		ClientSecret: os.Getenv("SECRET"),
		GrantType:    "authorization_code",
		Code:         code,
		RedirectURI:  os.Getenv("REDIRECTURI"),
	}


	jsonBody, err := json.Marshal(requestBody)

	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()	

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var authResponse models.AuthResponse
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	accountURL := fmt.Sprintf("https://%s.amocrm.ru/api/v4/account", subdomain)
	
	req, err := http.NewRequest("GET", accountURL, nil)

	if err != nil {
		return fmt.Errorf("failed to create account request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send account request: %w", err)
	}
	defer resp.Body.Close()

	accountBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read account response body: %w", err)
	}

	var res models.AccountResponse
	if err := json.Unmarshal(accountBody, &res); err != nil {
		return fmt.Errorf("failed to unmarshal account response: %w", err)
	}

	account := models.Account {
		AccessToken: authResponse.AccessToken,
		RefreshToken: authResponse.RefreshToken,
		Expires: authResponse.ExpiresIn,
		Subdomain: subdomain,
		AccountID: res.ID,
	}

	accountID, err := s.storage.AddAccount(account)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	contacts, err := s.GetContacts(accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	for i := range contacts {
		contacts[i].AccountID = account.AccountID
	}
	err = s.storage.SaveContacts(contacts)
	if err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}

	return nil
}