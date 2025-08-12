package country

import "context"

type UseCase interface {
	List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Country, int64, error)
	Create(ctx context.Context, country *Country) error
	Batch(ctx context.Context, countries *[]Country) error
	GetByID(ctx context.Context, id uint) (*Country, error)
	Update(ctx context.Context, country *Country) error
	Delete(ctx context.Context, id uint) error
}

type usecase struct {
	repository Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repository: repo}
}

func (u usecase) List(ctx context.Context, limit, offset int, sortBy, sortDir *string, search *string) (*[]Country, int64, error) {
	return u.repository.List(ctx, limit, offset, sortBy, sortDir, search)
}

func (u usecase) Create(ctx context.Context, country *Country) error {
	return u.repository.Create(ctx, country)
}

func (u usecase) Batch(ctx context.Context, countries *[]Country) error {
	return u.repository.Batch(ctx, countries)
}

func (u usecase) GetByID(ctx context.Context, id uint) (*Country, error) {
	return u.repository.GetByID(ctx, id)
}

func (u usecase) Update(ctx context.Context, country *Country) error {
	return u.repository.Update(ctx, country)
}

func (u usecase) Delete(ctx context.Context, id uint) error {
	return u.repository.Delete(ctx, id)
}
