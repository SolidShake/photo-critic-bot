package repository

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) Repository {
	return Repository{conn: conn}
}

func (r Repository) SaveInstaLink(chatID int64, link string) error {
	row := r.conn.QueryRow(context.Background(), "INSERT INTO chats (chat_id, link) VALUES ($1, $2) RETURNING id", chatID, link)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("save instagram link error: %s", err)
	}

	return nil
}

func (r Repository) GetInstaLink(chatID int64) (string, error) {
	var link string
	err := r.conn.QueryRow(context.Background(), "SELECT link FROM chats WHERE chat_id = $1 ORDER BY id DESC", chatID).Scan(&link)
	if err != nil {
		return "", fmt.Errorf("get instagram link error: %s", err)
	}

	return link, nil
}

func (r Repository) GetInstaForReview(chatID int64) (Chat, error) {
	var chat Chat
	err := pgxscan.Get(context.Background(), r.conn, &chat, `SELECT * FROM chats WHERE chat_id != $1 ORDER BY id DESC limit 1`, chatID)
	if err != nil {
		return Chat{}, fmt.Errorf("get insta link for review error: %s", err)
	}

	return chat, nil
}
