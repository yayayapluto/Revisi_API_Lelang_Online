package city

import "context"

type UseCase interface {
	List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]City, int64, error)
	Create(ctx context.Context, city *City) error
	Batch(ctx context.Context, cities *[]City) error
	GetByID(ctx context.Context, id uint) (*City, error)
	Update(ctx context.Context, city *City) error
	Delete(ctx context.Context, id uint) error
}

type usecase struct {
	repository Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repository: repo}
}

func (u usecase) List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]City, int64, error) {
	return u.repository.List(ctx, limit, offset, sortBy, sortDir, search)
}

func (u usecase) Create(ctx context.Context, city *City) error {
	return u.repository.Create(ctx, city)
}

func (u usecase) Batch(ctx context.Context, cities *[]City) error {
	return u.repository.Batch(ctx, cities)
}

func (u usecase) GetByID(ctx context.Context, id uint) (*City, error) {
	return u.repository.GetByID(ctx, id)
}

func (u usecase) Update(ctx context.Context, city *City) error {
	return u.repository.Update(ctx, city)
}

func (u usecase) Delete(ctx context.Context, id uint) error {
	return u.repository.Delete(ctx, id)
}
