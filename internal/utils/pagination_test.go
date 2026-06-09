package utils

import "testing"

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
		wantOffset   int
	}{
		{
			name:         "use defaults",
			page:         0,
			pageSize:     0,
			wantPage:     DefaultPage,
			wantPageSize: DefaultPageSize,
			wantOffset:   0,
		},
		{
			name:         "keep valid values",
			page:         3,
			pageSize:     20,
			wantPage:     3,
			wantPageSize: 20,
			wantOffset:   40,
		},
		{
			name:         "cap max page size",
			page:         2,
			pageSize:     1000,
			wantPage:     2,
			wantPageSize: MaxPageSize,
			wantOffset:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizePagination(tt.page, tt.pageSize)
			if got.Page != tt.wantPage {
				t.Fatalf("Page = %d, want %d", got.Page, tt.wantPage)
			}
			if got.PageSize != tt.wantPageSize {
				t.Fatalf("PageSize = %d, want %d", got.PageSize, tt.wantPageSize)
			}
			if got.Offset() != tt.wantOffset {
				t.Fatalf("Offset() = %d, want %d", got.Offset(), tt.wantOffset)
			}
		})
	}
}
