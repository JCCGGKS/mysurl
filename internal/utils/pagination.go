package utils

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

type Pagination struct {
	Page     int
	PageSize int
}

func NormalizePagination(page, pageSize int) Pagination {
	if page <= 0 {
		page = DefaultPage
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}
