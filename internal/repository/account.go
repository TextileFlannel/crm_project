package repository

import (
	"errors"
	"fmt"
	"http-server/danilkovalev/internal/models"

	"gorm.io/gorm"
)


type AccountSql struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountSql {
	return &AccountSql{db: db}
}


func (r *AccountSql) AddUnisenderKey(accountID int, key string) error {
	var account models.Account

	result := r.db.Where("account_id = ?", accountID).First(&account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("account not found")
		}
		return result.Error
	}
	account.UnisenderKey = key
	result = r.db.Save(&account)
	if result.Error != nil {
		return result.Error
	}

	return nil
}


func (r *AccountSql) AddAccount(account models.Account) (int, error) {
	if account.AccessToken == "" || account.RefreshToken == "" || account.Subdomain == "" {
		return 0, errors.New("tokens or subdomain is empty")
	}else if err := r.db.Create(&account).Error; err != nil {
		return 0, err
	}
	return account.ID, nil
}


func (r *AccountSql) GetAllAccounts() []models.Account {
	var accounts []models.Account
	if err := r.db.Find(&accounts).Error; err != nil {
		return nil
	}
	return accounts
}


func (r *AccountSql) GetAccount(id int) (models.Account, error) {
	var account models.Account

	if err := r.db.First(&account, id).Error; err != nil {
		return models.Account{}, err
	}
	
	return account, nil
}


func (r *AccountSql) DeleteAccounts() error {
	if err := r.db.Exec("delete from accounts").Error; err != nil {
		return err
	}
	return nil
}

func (r *AccountSql) DeleteAccount(id int) error {
	var account models.Account

	if err := r.db.Delete(&account, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *AccountSql) UpdateAccount(id int, NewAccount models.Account) error {
	NewAccount.ID = id

	if err := r.db.Save(&NewAccount).Error; err != nil {
		return err
	}
	
	return nil
}


func (r *AccountSql) GetEmails(accountID int) ([]string, error) {
	var contacts []models.Contact

	if err := r.db.Where("account_id = ?", accountID).Find(&contacts).Error; err != nil {
		return nil, err
	}

	emails := make([]string, len(contacts))
	for i, contact := range contacts {
		emails[i] = contact.Email
	}

	return emails, nil
}



func (r *AccountSql) SaveContacts(contacts []models.Contact) error {
	fmt.Println(contacts)

	if err := r.db.Create(&contacts).Error; err != nil {
		return err
	}

	return nil
}
