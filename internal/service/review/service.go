package service

import (
	review "github.com/SolidShake/photo-critic-bot/internal/repository/review"
)

type Service struct {
	repository review.Repository
}

func NewService(repository review.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) SaveReview(chatID, reviewedID int64, review string) error {
	return s.repository.SaveReview(chatID, reviewedID, review)
}

func (s *Service) GetReviews(chatID int64) ([]review.Review, error) {
	return s.repository.GetReviews(chatID)
}
