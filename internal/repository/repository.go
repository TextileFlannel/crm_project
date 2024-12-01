package repository

import (
	"http-server/danilkovalev/internal/models"

	"gorm.io/gorm"
)

type Account interface {
	AddAccount(account models.Account) (int, error)
	UpdateAccount(id int, NewAccount models.Account) error
	GetAllAccounts() []models.Account 
	GetAccount(id int) (models.Account, error)
	DeleteAccounts() error
	DeleteAccount(id int) error
	AddUnisenderKey(accountID int, key string) error
	SaveContacts(contacts []models.Contact) error
	GetEmails(accountID int) ([]string, error)
}

type Integration interface {
	AddIntegration(id int, integration models.Integration) error
	UpdateIntegration(id int, newIntegration models.Integration) error
	GetIntegrations() []models.Integration
	GetIntegration(id int) (models.Integration, error)
	DeleteIntegrations() error
	DeleteIntegration(id int) error
}

type Contact interface {
	AddContact(contact models.Contact) error
	DeleteContact(contact models.Contact) error
	UpdateContact(contact models.Contact) error
	GetUnisenderKey(accountID int) (string, error)
	GetEmailsContact(accountID int) ([]string, error)
}


type Repository struct {
    Account
	Integration
	Contact
}


func NewRepository(db *gorm.DB) *Repository {
    return &Repository{
        Account:     NewAccountRepository(db),
		Integration: NewIntegrationRepository(db),
		Contact:     NewContactMySQL(db),
    }
}