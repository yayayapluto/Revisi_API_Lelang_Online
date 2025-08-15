package country

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"slices"
	"strings"
)

type Repository interface {
	List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Country, int64, error)
	Create(ctx context.Context, country *Country) error
	Batch(ctx context.Context, countries *[]Country) error
	GetByID(ctx context.Context, id uint) (*Country, error)
	Update(ctx context.Context, country *Country) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Country, int64, error) {
	validSortColumns := []string{"id", "created_at", "nama", "kode"}
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

	var countries []Country
	query := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order(orderStr)

	if search != nil {
		searchParam := "%" + strings.ToUpper(*search) + "%"
		query = query.Where("kode LIKE ? OR nama LIKE ? OR nomor LIKE ?", searchParam, searchParam, searchParam)
	}

	res := query.Find(&countries)
	if err := res.Error; err != nil {
		return nil, 0, err
	}

	if search != nil {
		return &countries, res.RowsAffected, nil
	}

	var total int64
	if err := r.db.WithContext(ctx).Model(&Country{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &countries, total, nil
}

func (r *repository) Create(ctx context.Context, country *Country) error {
	var foundCountry Country
	err := r.db.WithContext(ctx).Where("kode = ? OR nama = ? OR nomor = ?", country.Kode, country.Nama, country.Nomor).First(&foundCountry).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		return fmt.Errorf("country already exists")
	}

	return r.db.WithContext(ctx).Create(country).Error
}

func (r *repository) Batch(ctx context.Context, countries *[]Country) error {
	if countries == nil || len(*countries) == 0 {
		return gorm.ErrInvalidData
	}

	for _, country := range *countries {
		var foundCountry Country
		err := r.db.WithContext(ctx).
			Where("kode = ? OR nama = ? OR nomor = ?", country.Kode, country.Nama, country.Nomor).
			First(&foundCountry).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil {
			return fmt.Errorf("country already exists: %s", country.Nama)
		}
	}

	return r.db.WithContext(ctx).Create(&countries).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Country, error) {
	var country Country
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&country).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &country, nil
}

func (r *repository) Update(ctx context.Context, country *Country) error {
	var foundCountry Country
	if err := r.db.WithContext(ctx).Where("id = ?", country.ID).First(&foundCountry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	updateData := map[string]interface{}{
		"kode":  country.Kode,
		"nama":  country.Nama,
		"nomor": country.Nomor,
	}

	return r.db.WithContext(ctx).Model(&Country{}).Where("id = ?", country.ID).Updates(updateData).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	var foundCountry Country
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&foundCountry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	return r.db.WithContext(ctx).Delete(&Country{}, id).Error
}
