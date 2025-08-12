package shared

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"slices"
)

type PaginationResult[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
}

func Paginate[T any](
	ctx context.Context,
	db *gorm.DB,
	limit, offset int,
	sortBy, sortDir *string,
	validSortColumns, validSortDir []string,
) (*PaginationResult[T], error) {

	sb := "id"
	sd := "asc"

	if sortBy != nil && slices.Contains(validSortColumns, *sortBy) {
		sb = *sortBy
	}
	if sortDir != nil && slices.Contains(validSortDir, *sortDir) {
		sd = *sortDir
	}

	orderStr := fmt.Sprintf("%s %s", sb, sd)

	var total int64
	if err := db.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		return nil, err
	}

	var data []T
	if err := db.WithContext(ctx).Limit(limit).Offset(offset).Order(orderStr).Find(&data).Error; err != nil {
		return nil, err
	}

	return &PaginationResult[T]{
		Data:  data,
		Total: total,
	}, nil
}
