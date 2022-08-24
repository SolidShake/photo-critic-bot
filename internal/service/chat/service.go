package service

import (
	chat "github.com/SolidShake/photo-critic-bot/internal/repository/chat"
)

type Service struct {
	repository chat.Repository
}

func NewService(repository chat.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) SaveInstaLink(chatID int64, link string) error {
	return s.repository.SaveInstaLink(chatID, link)
}

func (s *Service) GetInstaLink(chatID int64) (string, error) {
	return s.repository.GetInstaLink(chatID)
}

func (s *Service) GetInstaLinkForReview(chatID int64) (chat.Chat, error) {
	return s.repository.GetInstaForReview(chatID)
}
