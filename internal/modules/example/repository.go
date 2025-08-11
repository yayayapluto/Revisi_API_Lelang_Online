package example

import (
	"context"
	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context, limit, offset int, sortBy, sortDir *string) (*[]Country, error)
	Create(ctx context.Context, country *Country) error
	GetByID(ctx context.Context, id int) (*Country, error)
	Update(ctx context.Context, country *Country) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context, limit, offset int, sortBy, sortDir *string) (*[]Country, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Create(ctx context.Context, country *Country) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetByID(ctx context.Context, id int) (*Country, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Update(ctx context.Context, country *Country) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
