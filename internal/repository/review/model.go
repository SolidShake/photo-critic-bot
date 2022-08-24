package review

import "time"

type Review struct {
	ID         int64     `db:"id"`
	ChatID     int64     `db:"chat_id"`
	ReviewedID int64     `db:"reviewed_id"`
	Review     string    `db:"review"`
	CreatedAt  time.Time `db:"created_at"`
}

func (r Review) GetFormatedTime() string {
	return r.CreatedAt.Format("15:04:05 02-01-2006")
}
