package repository

import (
	"errors"
	"http-server/danilkovalev/internal/models"

	"gorm.io/gorm"
)

type IntegrationSql struct {
	db *gorm.DB
}

func NewIntegrationRepository(db *gorm.DB) *IntegrationSql{
	return &IntegrationSql{db: db}
}


func (r *IntegrationSql) AddIntegration(id int, integration models.Integration) error {
	var account models.Account
	if err := r.db.First(&account, id).Error; err != nil {
		return errors.New("account with this id was not found")
	}
	integration.Account_ID = id

	if err := r.db.Create(&integration).Error; err != nil {
		return errors.New("error in creation")
	}
	return nil
}


func (r *IntegrationSql) GetIntegration(id int) (models.Integration, error) {
	var integration models.Integration

	if err := r.db.First(&integration, id).Error; err != nil {
		return models.Integration{}, err
	}

	return integration, nil
}


func (r *IntegrationSql) GetIntegrations() []models.Integration {
	var integrations []models.Integration

	if err := r.db.Find(&integrations).Error; err != nil {
		return []models.Integration{}
	}

	return integrations
}


func (r *IntegrationSql) DeleteIntegrations() error {
	if err := r.db.Exec("delete from integartion").Error; err != nil {
		return err
	}
	return nil
}


func (r *IntegrationSql) DeleteIntegration(id int) error {
		var integration models.Integration
		
		if err := r.db.Delete(&integration, id).Error; err != nil {
			return err
		}

		return nil
}


func (r *IntegrationSql) UpdateIntegration(id int, newIntegration models.Integration) error {
	newIntegration.ID =  id
	if err := r.db.Save(&newIntegration).Error; err != nil {
		return err
	}
	return nil
}