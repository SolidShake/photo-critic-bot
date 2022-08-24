package repository

import "time"

type Chat struct {
	ID        int64     `db:"id"`
	ChatID    int64     `db:"chat_id"`
	Link      string    `db:"link"`
	CreatedAt time.Time `db:"created_at"`
}
