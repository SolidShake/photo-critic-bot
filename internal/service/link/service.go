package service

import (
	link "github.com/SolidShake/photo-critic-bot/internal/repository/link"
)

type LinkService struct {
	linkRepository link.Repository
}

func NewLinkService(linkRepository link.Repository) *LinkService {
	return &LinkService{linkRepository: linkRepository}
}

func (l *LinkService) SaveInstaLink(chatID int64, link string) error {
	return l.linkRepository.SaveInstaLink(chatID, link)
}

func (l *LinkService) GetInstaLink(chatID int64) (string, error) {
	return l.linkRepository.GetInstaLink(chatID)
}
