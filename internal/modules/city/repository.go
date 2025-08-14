package city

import (
	"context"
	"errors"
	"fmt"
	cityDto "github.com/API_Lelang_Online_Go/internal/modules/city/dto"
	"gorm.io/gorm"
	"slices"
	"strings"
)

type Repository interface {
	List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]City, int64, error)
	Create(ctx context.Context, city *City) error
	Batch(ctx context.Context, cities *[]City) error
	GetByID(ctx context.Context, id uint) (*City, error)
	Update(ctx context.Context, city *City) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]City, int64, error) {
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

	var cities []City
	query := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order(orderStr)

	if search != nil {
		searchParam := "%" + strings.ToUpper(*search) + "%"
		query = query.Where("code LIKE ? OR nama LIKE ? OR full_code LIKE ?", searchParam, searchParam, searchParam)
	}

	res := query.Find(&cities)
	if err := res.Error; err != nil {
		return nil, 0, err
	}

	if search != nil {
		return &cities, res.RowsAffected, nil
	}

	var total int64
	if err := r.db.WithContext(ctx).Model(&City{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &cities, total, nil
}

func (r *repository) Create(ctx context.Context, city *City) error {
	var foundCity City
	err := r.db.WithContext(ctx).Where("code = ? OR nama = ? OR full_code = ?", city.Code, city.Nama, city.FullCode).First(&foundCity).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		return fmt.Errorf("city already exists")
	}

	return r.db.WithContext(ctx).Create(city).Error
}

func (r *repository) Batch(ctx context.Context, cities *[]City) error {
	if cities == nil || len(*cities) == 0 {
		return gorm.ErrInvalidData
	}

	for _, city := range *cities {
		var foundCity City
		err := r.db.WithContext(ctx).
			Where("code = ? OR nama = ? OR full_code = ?", city.Code, city.Nama, city.FullCode).
			First(&foundCity).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil {
			return fmt.Errorf("city already exists: %s", city.Nama)
		}
	}

	return r.db.WithContext(ctx).Create(&cities).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*City, error) {
	var city City
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&city).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &city, nil
}

func (r *repository) Update(ctx context.Context, city *City) error {
	var foundCity City
	if err := r.db.WithContext(ctx).Where("id = ?", city.ID).First(&foundCity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	updateData := cityDto.UpdateRequest{
		Code:     &city.Code,
		Nama:     &city.Nama,
		FullCode: &city.FullCode,
	}

	return r.db.WithContext(ctx).Model(&City{}).Where("id = ?", city.ID).Updates(updateData).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	var foundCity City
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&foundCity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	return r.db.WithContext(ctx).Delete(&City{}, id).Error
}
