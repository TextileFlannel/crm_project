package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"http-server/danilkovalev/internal/models"
	"http-server/danilkovalev/internal/repository"

	"log"
	"github.com/beanstalkd/go-beanstalk"
)

type AccountService struct {
	storage repository.Account
	beanstalkClient *BeanstalkClient
}

func NewAccountService(storage repository.Account, beanstalkClient *BeanstalkClient) *AccountService {
    return &AccountService{
        storage:     storage,
        beanstalkClient: beanstalkClient,
    }
}


func (s *AccountService) GetAllAccounts() []models.Account {
	accounts := s.storage.GetAllAccounts()
	
	return accounts
}


func (s *AccountService) CreateAccount(account models.Account) (int, error) {
	id, err := s.storage.AddAccount(account)
	return id, err
}


func (s *AccountService) GetAccount(id int) (models.Account, error) {
	return s.storage.GetAccount(id)
}


func (s *AccountService) DeleteAccounts() error {
	err := s.storage.DeleteAccounts()
	return err
}


func (s *AccountService) DeleteAccount(id int) error {
	err := s.storage.DeleteAccount(id)
	return err
}


func (s *AccountService) UpdateAccount(id int, account models.Account) error {
	err := s.storage.UpdateAccount(id, account)
	return err
}


func (s *AccountService) GetContacts(accountID int) ([]models.Contact, error) {
	account, err := s.storage.GetAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	url := fmt.Sprintf("https://%s.amocrm.ru/api/v4/contacts", account.Subdomain)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+account.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response models.AmoContactsResponse

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch contacts: %s", body)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	var contacts []models.Contact
	for _, c := range response.Embedded.Contacts {
		var email string
		for _, field := range c.CustomFieldsValues {
			if field.FieldName == "Email" {
				for _, value := range field.Values {
					email = value.Value
				}
			}
		}

		contacts = append(contacts, models.Contact{
			Name:      c.Name,
			Email:     email,
			AccountID: accountID,
			ID:        c.ID,
		})
	}

	return contacts, nil
}

func (s *AccountService) AddUnisenderKey(accountID int, key string) error {
	task := models.Task{
		AccountID:    accountID,
		UnisenderKey: key,
		TaskType:     "Account",
	}

	if err := addTaskToQueue(s.beanstalkClient.Conn, task); err != nil {
		log.Printf("Failed to add task to queue: %v", err)
		return fmt.Errorf("failed to add task to queue: %w", err)
	}

	return nil
}

func addTaskToQueue(conn *beanstalk.Conn, task models.Task) error {
	taskData, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task data: %w", err)
	}

	id, err := conn.Put(taskData, 1, 0, 120)
	if err != nil {
		return fmt.Errorf("failed to put task in queue: %w", err)
	}

	log.Printf("Task added with ID: %d", id)
	return nil
}
