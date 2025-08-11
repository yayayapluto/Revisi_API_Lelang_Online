package example

import "context"

type UseCase interface {
	List(ctx context.Context, limit, offset int, sortBy, sortDir *string) (*[]Country, error)
	Create(ctx context.Context, country *Country) error
	GetByID(ctx context.Context, id int) (*Country, error)
	Update(ctx context.Context, country *Country) error
	Delete(ctx context.Context, id int) error
}

type usecase struct {
	repository Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repository: repo}
}

func (u usecase) List(ctx context.Context, limit, offset int, sortBy, sortDir *string) (*[]Country, error) {
	//TODO implement me
	panic("implement me")
}

func (u usecase) Create(ctx context.Context, country *Country) error {
	//TODO implement me
	panic("implement me")
}

func (u usecase) GetByID(ctx context.Context, id int) (*Country, error) {
	//TODO implement me
	panic("implement me")
}

func (u usecase) Update(ctx context.Context, country *Country) error {
	//TODO implement me
	panic("implement me")
}

func (u usecase) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
