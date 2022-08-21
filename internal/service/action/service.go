package service

import (
	action "github.com/SolidShake/photo-critic-bot/internal/repository/action"
)

type ActionService struct {
	actionRepository action.Repository
}

func NewActionService(actionRepository action.Repository) *ActionService {
	return &ActionService{actionRepository: actionRepository}
}

func (a *ActionService) GetLastAction(chatID int64) (string, error) {
	return a.actionRepository.GetLastAction(chatID)
}

func (a *ActionService) SaveAction(chatID int64, action string) error {
	return a.actionRepository.SaveAction(chatID, action)
}
