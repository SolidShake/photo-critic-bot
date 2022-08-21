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

func (r Repository) GetLastAction(chatID int64) (string, error) {
	var action string
	err := r.conn.QueryRow(context.Background(), "SELECT action FROM actions WHERE chat_id = $1 ORDER BY id DESC", chatID).Scan(&action)
	if err != nil {
		return "", fmt.Errorf("select last action error: %s", err)
	}

	return action, nil
}

func (r Repository) SaveAction(chatID int64, action string) error {
	row := r.conn.QueryRow(context.Background(), "INSERT INTO actions (chat_id, action) VALUES ($1, $2) RETURNING id", chatID, action)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("save action error: %s", err)
	}

	return nil
}
