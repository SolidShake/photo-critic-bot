package service

import (
	action "github.com/SolidShake/photo-critic-bot/internal/repository/action"
)

type Service struct {
	repository action.Repository
}

func NewService(repository action.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetLastAction(chatID int64) (string, error) {
	return s.repository.GetLastAction(chatID)
}

func (s *Service) SaveAction(chatID int64, action string) error {
	return s.repository.SaveAction(chatID, action)
}
