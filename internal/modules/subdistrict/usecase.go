package subdistrict

import (
	"context"
	"errors"
)

type UseCase interface {
	List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Subdistrict, int64, error)
	Create(ctx context.Context, subdistrict *Subdistrict) error
	Batch(ctx context.Context, subdistricts *[]Subdistrict) error
	GetByID(ctx context.Context, id uint) (*Subdistrict, error)
	Update(ctx context.Context, subdistrict *Subdistrict) error
	Delete(ctx context.Context, id uint) error
}

type usecase struct {
	repository Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repository: repo}
}

func (u usecase) List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Subdistrict, int64, error) {
	return u.repository.List(ctx, limit, offset, sortBy, sortDir, search)
}

func (u usecase) Create(ctx context.Context, subdistrict *Subdistrict) error {
	if subdistrict.Code == "" || subdistrict.FullCode == "" || subdistrict.Nama == "" {
		return errors.New("all field must be filled")
	}
	return u.repository.Create(ctx, subdistrict)
}

func (u usecase) Batch(ctx context.Context, subdistricts *[]Subdistrict) error {
	return u.repository.Batch(ctx, subdistricts)
}

func (u usecase) GetByID(ctx context.Context, id uint) (*Subdistrict, error) {
	return u.repository.GetByID(ctx, id)
}

func (u usecase) Update(ctx context.Context, subdistrict *Subdistrict) error {
	return u.repository.Update(ctx, subdistrict)
}

func (u usecase) Delete(ctx context.Context, id uint) error {
	return u.repository.Delete(ctx, id)
}
