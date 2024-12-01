package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	repository "http-server/danilkovalev/internal/repository"
	"http-server/danilkovalev/internal/models"
	"log"
	"net/http"
	"net/url"

	"github.com/beanstalkd/go-beanstalk"
)

type ContactService struct {
	repository 		repository.Contact
	beanstalkClient *BeanstalkClient

}

func NewContactService(repository repository.Contact, beanstalkClient *BeanstalkClient) *ContactService {
	return &ContactService{
		repository:      repository,
		beanstalkClient: beanstalkClient,
	}
}

func (s *ContactService) AddContact(contact models.Contact) error {
	return s.repository.AddContact(contact)
}

func (s *ContactService) DeleteContact(contact models.Contact) error {
	return s.repository.DeleteContact(contact)
}

func (s *ContactService) UpdateContact(contact models.Contact) error {
	return s.repository.UpdateContact(contact)
}

func (s *ContactService) GetUnisenderKey(accountID int) (string, error) {
	return s.repository.GetUnisenderKey(accountID)
}

func (s *ContactService) ToUnisenderContact(accountID int, key string) error {
	emails, err := s.repository.GetEmailsContact(accountID)
	if err != nil {
		return fmt.Errorf("failed to get emails: %w", err)
	}

	fieldNames := []string{"email"}
	data := make([][]string, len(emails))

	for i, email := range emails {
		data[i] = []string{email}
	}

	apiURL := "https://api.unisender.com/ru/api/importContacts"
	params := url.Values{}
	params.Set("format", "json")
	params.Set("api_key", key)

	for i, field := range fieldNames {
		params.Set(fmt.Sprintf("field_names[%d]", i), field)
	}

	for i, row := range data {
		for j, value := range row {
			params.Set(fmt.Sprintf("data[%d][%d]", i, j), value)
		}
	}

	resp, err := http.Post(apiURL+"?"+params.Encode(), "application/x-www-form-urlencoded", bytes.NewBuffer(nil))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}


func (s *ContactService) ChangeContact(contact models.Contact, changeType string) error {
	contactChange := models.ContactChange{
		Name:       contact.Name,
		Email:      contact.Email,
		AccountID:  contact.AccountID,
		ID:         contact.ID,
		TypeChange: changeType,
		TaskType:   "Contact",
	}

	if err := addContactToQueue(s.beanstalkClient.Conn, contactChange); err != nil {
		log.Printf("Failed to add task to queue: %v", err)
		return fmt.Errorf("failed to add task to queue: %w", err)
	}

	return nil

}

func addContactToQueue(conn *beanstalk.Conn, contact models.ContactChange) error {
	taskData, err := json.Marshal(contact)
	if err != nil {
		log.Printf("Failed to marshal task data: %v", err)
		return err
	}

	id, err := conn.Put(taskData, 1, 0, 120)
	if err != nil {
		log.Printf("Failed to put task in queue: %v", err)
		return err
	}

	log.Printf("Task added with ID: %d", id)
	return nil
}

