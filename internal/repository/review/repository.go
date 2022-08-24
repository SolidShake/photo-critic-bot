package review

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

func (r Repository) SaveReview(chatID, reviewedID int64, review string) error {
	row := r.conn.QueryRow(context.Background(), "INSERT INTO reviews (chat_id, reviewed_id, review) VALUES ($1, $2, $3) RETURNING id", chatID, reviewedID, review)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("save review error: %s", err)
	}

	return nil
}

func (r Repository) GetReviews(chatID int64) ([]Review, error) {
	var reviews []Review
	err := pgxscan.Select(context.Background(), r.conn, &reviews, `SELECT * FROM reviews WHERE reviewed_id = $1 ORDER BY created_at DESC`, chatID)
	if err != nil {
		return nil, fmt.Errorf("get reviews error: %s", err)
	}

	return reviews, nil
}
