package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) Repository {
	return Repository{conn: conn}
}

func (r Repository) SaveInstaLink(chatID int64, link string) error {
	row := r.conn.QueryRow(context.Background(), "INSERT INTO links (chat_id, link) VALUES ($1, $2) RETURNING id", chatID, link)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("insert db error: %s", err)
	}

	return nil
}

func (r Repository) GetInstaLink(chatID int64) (string, error) {
	var link string
	err := r.conn.QueryRow(context.Background(), "SELECT link FROM links WHERE chat_id = $1 ORDER BY id DESC", chatID).Scan(&link)
	if err != nil {
		return "", fmt.Errorf("insert db error: %s", err)
	}

	return link, nil
}
