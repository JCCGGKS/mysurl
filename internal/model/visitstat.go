package model

import "time"

type VisitStat struct {
	ID          uint64    `db:"id"`
	ShortLinkID uint64    `db:"short_link_id"`
	VisitCount  uint64    `db:"visit_count"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
