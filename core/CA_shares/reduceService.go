package CA_shares

import (
	"CA_shares/model"
	"CA_shares/serviceErrors"
	"context"
)

type NewsRepository interface {
	GetAll(ctx context.Context) ([]model.NewsListPublic, error)
}

type ReduceService struct {
	repo NewsRepository
}

func NewReduceService(repo NewsRepository) *ReduceService {
	return &ReduceService{repo: repo}
}

func (rs *ReduceService) List(ctx context.Context) (model.NewsListPublicResponse, error) {
	items, err := rs.repo.GetAll(ctx)
	if err != nil {
		return model.NewsListPublicResponse{}, serviceErrors.ErrInternal
	}

	if items == nil {
		items = []model.NewsListPublic{}
	}

	return model.NewsListPublicResponse{
		NextCursor: "",
		List:       items,
	}, nil
}
