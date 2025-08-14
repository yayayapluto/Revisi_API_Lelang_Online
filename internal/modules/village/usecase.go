package village

import (
	"context"
	"errors"
)

type UseCase interface {
	List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Village, int64, error)
	Create(ctx context.Context, village *Village) error
	Batch(ctx context.Context, villages *[]Village) error
	GetByID(ctx context.Context, id uint) (*Village, error)
	Update(ctx context.Context, village *Village) error
	Delete(ctx context.Context, id uint) error
}

type usecase struct {
	repository Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repository: repo}
}

func (u usecase) List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Village, int64, error) {
	return u.repository.List(ctx, limit, offset, sortBy, sortDir, search)
}

func (u usecase) Create(ctx context.Context, village *Village) error {
	if village.Code == "" || village.FullCode == "" || village.Nama == "" {
		return errors.New("all field must be filled")
	}
	return u.repository.Create(ctx, village)
}

func (u usecase) Batch(ctx context.Context, villages *[]Village) error {
	return u.repository.Batch(ctx, villages)
}

func (u usecase) GetByID(ctx context.Context, id uint) (*Village, error) {
	return u.repository.GetByID(ctx, id)
}

func (u usecase) Update(ctx context.Context, village *Village) error {
	return u.repository.Update(ctx, village)
}

func (u usecase) Delete(ctx context.Context, id uint) error {
	return u.repository.Delete(ctx, id)
}
