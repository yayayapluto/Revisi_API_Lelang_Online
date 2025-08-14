package province

import (
	"context"
	"errors"
	"fmt"
	"github.com/API_Lelang_Online_Go/internal/modules/province/dto"
	"gorm.io/gorm"
	"slices"
	"strings"
)

type Repository interface {
	List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Province, int64, error)
	Create(ctx context.Context, province *Province) error
	Batch(ctx context.Context, provinces *[]Province) error
	GetByID(ctx context.Context, id uint) (*Province, error)
	Update(ctx context.Context, province *Province) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Province, int64, error) {
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

	var provinces []Province
	query := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order(orderStr)

	if search != nil {
		searchParam := "%" + strings.ToUpper(*search) + "%"
		query = query.Where("code LIKE ? OR nama LIKE ? OR full_code LIKE ?", searchParam, searchParam, searchParam)
	}

	res := query.Find(&provinces)
	if err := res.Error; err != nil {
		return nil, 0, err
	}

	if search != nil {
		return &provinces, res.RowsAffected, nil
	}

	var total int64
	if err := r.db.WithContext(ctx).Model(&Province{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &provinces, total, nil
}

func (r *repository) Create(ctx context.Context, province *Province) error {
	var foundProvince Province
	err := r.db.WithContext(ctx).Where("code = ? OR nama = ? OR full_code = ?", province.Code, province.Nama, province.FullCode).First(&foundProvince).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		return fmt.Errorf("province already exists")
	}

	return r.db.WithContext(ctx).Create(province).Error
}

func (r *repository) Batch(ctx context.Context, provinces *[]Province) error {
	if provinces == nil || len(*provinces) == 0 {
		return gorm.ErrInvalidData
	}

	for _, province := range *provinces {
		var foundProvince Province
		err := r.db.WithContext(ctx).
			Where("code = ? OR nama = ? OR full_code = ?", province.Code, province.Nama, province.FullCode).
			First(&foundProvince).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil {
			return fmt.Errorf("province already exists: %s", province.Nama)
		}
	}

	return r.db.WithContext(ctx).Create(&provinces).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Province, error) {
	var province Province
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&province).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &province, nil
}

func (r *repository) Update(ctx context.Context, province *Province) error {
	var foundProvince Province
	if err := r.db.WithContext(ctx).Where("id = ?", province.ID).First(&foundProvince).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	updateData := provinceDto.UpdateRequest{
		Code:     &province.Code,
		Nama:     &province.Nama,
		FullCode: &province.FullCode,
	}

	return r.db.WithContext(ctx).Model(&Province{}).Where("id = ?", province.ID).Updates(updateData).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	var foundProvince Province
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&foundProvince).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	return r.db.WithContext(ctx).Delete(&Province{}, id).Error
}
