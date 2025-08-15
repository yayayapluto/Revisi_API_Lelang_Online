package village

import (
	"context"
	"errors"
	"fmt"
	villageDto "github.com/API_Lelang_Online_Go/internal/modules/village/dto"
	"gorm.io/gorm"
	"slices"
	"strings"
)

type Repository interface {
	List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Village, int64, error)
	Create(ctx context.Context, village *Village) error
	Batch(ctx context.Context, villages *[]Village) error
	GetByID(ctx context.Context, id uint) (*Village, error)
	Update(ctx context.Context, village *Village) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Village, int64, error) {
	validSortColumns := []string{"id", "created_at", "nama", "code", "full_code"}
	validSortDir := []string{"asc", "desc"}

	sb := "id"
	sd := "asc"

	if sortBy != nil && slices.Contains(validSortColumns, *sortBy) {
		sb = *sortBy
	}
	if sortDir != nil && slices.Contains(validSortDir, *sortDir) {
		sd = *sortDir
	}

	orderStr := fmt.Sprintf("%s %s", sb, sd)

	var villages []Village
	query := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order(orderStr)

	if search != nil {
		searchParam := "%" + strings.ToUpper(*search) + "%"
		query = query.Where("code LIKE ? OR nama LIKE ? OR full_code LIKE ?", searchParam, searchParam, searchParam)
	}

	res := query.Find(&villages)
	if err := res.Error; err != nil {
		return nil, 0, err
	}

	if search != nil {
		return &villages, res.RowsAffected, nil
	}

	var total int64
	if err := r.db.WithContext(ctx).Model(&Village{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &villages, total, nil
}

func (r *repository) Create(ctx context.Context, village *Village) error {
	var foundVillage Village
	err := r.db.WithContext(ctx).Where("code = ? OR nama = ? OR full_code = ?", village.Code, village.Nama, village.FullCode).First(&foundVillage).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		return fmt.Errorf("village already exists")
	}

	return r.db.WithContext(ctx).Create(village).Error
}

func (r *repository) Batch(ctx context.Context, villages *[]Village) error {
	if villages == nil || len(*villages) == 0 {
		return gorm.ErrInvalidData
	}

	for _, village := range *villages {
		var foundVillage Village
		err := r.db.WithContext(ctx).
			Where("code = ? OR nama = ? OR full_code = ?", village.Code, village.Nama, village.FullCode).
			First(&foundVillage).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil {
			return fmt.Errorf("village already exists: %s", village.Nama)
		}
	}

	return r.db.WithContext(ctx).Create(&villages).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Village, error) {
	var village Village
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&village).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &village, nil
}

func (r *repository) Update(ctx context.Context, village *Village) error {
	var foundVillage Village
	if err := r.db.WithContext(ctx).Where("id = ?", village.ID).First(&foundVillage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	updateData := villageDto.UpdateRequest{
		Code:     &village.Code,
		Nama:     &village.Nama,
		FullCode: &village.FullCode,
	}

	return r.db.WithContext(ctx).Model(&Village{}).Where("id = ?", village.ID).Updates(updateData).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	var foundVillage Village
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&foundVillage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	return r.db.WithContext(ctx).Delete(&Village{}, id).Error
}
