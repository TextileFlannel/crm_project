package repository

import (
	"http-server/danilkovalev/internal/models"
	"log"

	"gorm.io/gorm"
)

type ContactMySQL struct {
	db *gorm.DB
}

func NewContactMySQL(db *gorm.DB) *ContactMySQL {
	return &ContactMySQL{db: db}
}

func (r *ContactMySQL) AddContact(contact models.Contact) error {
	if err := r.db.Create(&contact).Error; err != nil {
		log.Printf("Error adding contact: %v", err)
		return err
	}
	return nil
}

func (r *ContactMySQL) DeleteContact(contact models.Contact) error {
	if err := r.db.Delete(&contact).Error; err != nil {
		log.Printf("Error deleting contact: %v", err)
		return err
	}
	return nil
}

func (r *ContactMySQL) UpdateContact(contact models.Contact) error {
	if err := r.db.Save(&contact).Error; err != nil {
		log.Printf("Error updating contact: %v", err)
		return err
	}
	return nil
}

func (r *ContactMySQL) GetUnisenderKey(accountID int) (string, error) {
	var account models.Account

	if err := r.db.Where("account_id = ?", accountID).First(&account).Error; err != nil {
		return "", err
	}

	return account.UnisenderKey, nil
}

func (r *ContactMySQL) GetEmailsContact(accountID int) ([]string, error) {
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