package service

import (
	"http-server/danilkovalev/internal/models"
	"http-server/danilkovalev/internal/repository"
)

type IntegrationService struct {
	storage repository.Integration
}


func NewIntegrationService(storage repository.Integration) *IntegrationService {
	return &IntegrationService{storage: storage}
}


func (s *IntegrationService) CreateIntegration(accountId int, integration models.Integration) error {
	err := s.storage.AddIntegration(accountId, integration)
	return err
}


func (s *IntegrationService) UpdateIntegration(id int, integration models.Integration) error {
	err := s.storage.UpdateIntegration(id, integration)
	return err
}


func (s *IntegrationService) GetIntegrations() []models.Integration {
	integrations := s.storage.GetIntegrations()
	return integrations
}


func (s *IntegrationService) GetIntegration(id int) (models.Integration, error) {
	return s.storage.GetIntegration(id)
}


func (s *IntegrationService) DeleteIntegrations() error {
	err := s.storage.DeleteIntegrations()
	return err
}


func (s *IntegrationService) DeleteIntegration(id int) error {
	err := s.storage.DeleteIntegration(id)
	return err
}