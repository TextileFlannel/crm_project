package service

import (
	"http-server/danilkovalev/internal/models"
	"http-server/danilkovalev/internal/repository"
)


type Account interface {
	CreateAccount(account models.Account) (int, error)
	UpdateAccount(id int, account models.Account) error
	GetAllAccounts() []models.Account
	GetAccount(id int) (models.Account, error)
	DeleteAccounts() error
	DeleteAccount(id int) error
	Authorization(code, clientID, referer string) error
	AddUnisenderKey(accountID int, key string) error
	GetContacts(accountId int) ([]models.Contact, error)
}

type Integration interface {
	CreateIntegration(accountId int, integration models.Integration) error
	UpdateIntegration(id int, integration models.Integration) error
	GetIntegrations() []models.Integration
	GetIntegration(id int) (models.Integration, error)
	DeleteIntegrations() error
	DeleteIntegration(id int) error
}

type Contact interface {
	AddContact(contact models.Contact) error
	UpdateContact(contact models.Contact) error
	DeleteContact(contact models.Contact) error
	GetUnisenderKey(accountID int) (string, error)
	ToUnisenderContact(accountID int, key string) error
	ChangeContact(contact models.Contact, changeType string) error
}

type Service struct {
	Account
	Integration
	Contact
}

func NewService(repo *repository.Repository, client *BeanstalkClient) *Service {
	return &Service{
		Account:      NewAccountService(repo, client),
		Integration:  NewIntegrationService(repo),
		Contact:      NewContactService(repo, client),
	}
}
