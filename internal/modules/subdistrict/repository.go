package subdistrict

import (
	"context"
	"errors"
	"fmt"
	subdistrictDto "github.com/API_Lelang_Online_Go/internal/modules/subdistrict/dto"
	"gorm.io/gorm"
	"slices"
	"strings"
)

type Repository interface {
	List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Subdistrict, int64, error)
	Create(ctx context.Context, subdistrict *Subdistrict) error
	Batch(ctx context.Context, subdistricts *[]Subdistrict) error
	GetByID(ctx context.Context, id uint) (*Subdistrict, error)
	Update(ctx context.Context, subdistrict *Subdistrict) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Subdistrict, int64, error) {
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

	var subdistricts []Subdistrict
	query := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order(orderStr)

	if search != nil {
		searchParam := "%" + strings.ToUpper(*search) + "%"
		query = query.Where("code LIKE ? OR nama LIKE ? OR full_code LIKE ?", searchParam, searchParam, searchParam)
	}

	res := query.Find(&subdistricts)
	if err := res.Error; err != nil {
		return nil, 0, err
	}

	if search != nil {
		return &subdistricts, res.RowsAffected, nil
	}

	var total int64
	if err := r.db.WithContext(ctx).Model(&Subdistrict{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &subdistricts, total, nil
}

func (r *repository) Create(ctx context.Context, subdistrict *Subdistrict) error {
	var foundSubdistrict Subdistrict
	err := r.db.WithContext(ctx).Where("code = ? OR nama = ? OR full_code = ?", subdistrict.Code, subdistrict.Nama, subdistrict.FullCode).First(&foundSubdistrict).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		return fmt.Errorf("subdistrict already exists")
	}

	return r.db.WithContext(ctx).Create(subdistrict).Error
}

func (r *repository) Batch(ctx context.Context, subdistricts *[]Subdistrict) error {
	if subdistricts == nil || len(*subdistricts) == 0 {
		return gorm.ErrInvalidData
	}

	for _, subdistrict := range *subdistricts {
		var foundSubdistrict Subdistrict
		err := r.db.WithContext(ctx).
			Where("code = ? OR nama = ? OR full_code = ?", subdistrict.Code, subdistrict.Nama, subdistrict.FullCode).
			First(&foundSubdistrict).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil {
			return fmt.Errorf("subdistrict already exists: %s", subdistrict.Nama)
		}
	}

	return r.db.WithContext(ctx).Create(&subdistricts).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Subdistrict, error) {
	var subdistrict Subdistrict
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&subdistrict).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &subdistrict, nil
}

func (r *repository) Update(ctx context.Context, subdistrict *Subdistrict) error {
	var foundSubdistrict Subdistrict
	if err := r.db.WithContext(ctx).Where("id = ?", subdistrict.ID).First(&foundSubdistrict).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	updateData := subdistrictDto.UpdateRequest{
		Code:     &subdistrict.Code,
		Nama:     &subdistrict.Nama,
		FullCode: &subdistrict.FullCode,
	}

	return r.db.WithContext(ctx).Model(&Subdistrict{}).Where("id = ?", subdistrict.ID).Updates(updateData).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	var foundSubdistrict Subdistrict
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&foundSubdistrict).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	return r.db.WithContext(ctx).Delete(&Subdistrict{}, id).Error
}
