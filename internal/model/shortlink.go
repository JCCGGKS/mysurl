package model

import "time"

const ShortLinkStatusActive = 1

type ShortLink struct {
	ID          uint64     `db:"id"`
	UserID      *uint64    `db:"user_id"`
	ShortCode   string     `db:"short_code"`
	OriginalURL string     `db:"original_url"`
	URLHash     string     `db:"url_hash"`
	ExpiresAt   *time.Time `db:"expires_at"`
	Status      uint8      `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}
