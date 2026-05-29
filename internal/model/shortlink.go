package model

import "time"

const ShortLinkStatusActive = 1

type ShortLink struct {
	ID          uint64     `db:"id"`
	ShortCode   string     `db:"short_code"`
	OriginalURL string     `db:"original_url"`
	URLHash     string     `db:"url_hash"`
	VisitCount  uint64     `db:"visit_count"`
	Status      uint8      `db:"status"`
	ExpiresAt   *time.Time `db:"expires_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}
