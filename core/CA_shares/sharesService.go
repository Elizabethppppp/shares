package CA_shares

import (
	"CA_shares/filtres"
	"CA_shares/model"
	"CA_shares/serviceErrors"
	"context"
)

type NewsRepository interface {
	GetAll(ctx context.Context) ([]model.NewsListPublic, error)
	GetById(ctx context.Context, id int64) (model.NewsDetailPublic, error)
	GetAllList(ctx context.Context, f filtres.ListFilter) ([]model.NewsList, int64, error)
	GetByIdAdminRole(ctx context.Context, id int64) (model.NewsDetailAdmin, error)
	UpdateActive(ctx context.Context, id int64) (model.NewsUpdateActive, error)
	Delete(ctx context.Context, id int64) (string, string, error)
	Create(ctx context.Context, p model.NewsPostNew) (model.NewsList, error)
	ToEdit(ctx context.Context, e model.NewsEditNotification) (model.NewsList, string, string, error)
}

type SharesService struct {
	repo NewsRepository
}

func NewSharesService(repo NewsRepository) *SharesService {
	return &SharesService{repo: repo}
}

func (s *SharesService) List(ctx context.Context) (model.NewsListPublicResponse, error) {
	items, err := s.repo.GetAll(ctx)
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

func (s *SharesService) NotificationById(ctx context.Context, id int64) (model.NewsDetailPublic, error) {
	items, err := s.repo.GetById(ctx, id)
	if err != nil {
		return model.NewsDetailPublic{}, err
	}

	return items, nil
}

func (s *SharesService) GetAllNews(ctx context.Context, f filtres.ListFilter) (model.NewsListResponse, error) {
	items, total, err := s.repo.GetAllList(ctx, f)
	if err != nil {
		return model.NewsListResponse{}, err
	}
	if items == nil {
		items = []model.NewsList{}
	}
	return model.NewsListResponse{
		TotalCount: total,
		List:       items,
	}, nil
}

func (s *SharesService) NotificationByIdAdminRole(ctx context.Context, id int64) (model.NewsDetailAdmin, error) {
	items, err := s.repo.GetByIdAdminRole(ctx, id)
	if err != nil {
		return model.NewsDetailAdmin{}, err
	}
	return items, nil
}

func (s *SharesService) UpdateStatus(ctx context.Context, id int64) (model.NewsUpdateActive, error) {
	news, err := s.repo.UpdateActive(ctx, id)
	if err != nil {
		return model.NewsUpdateActive{}, err
	}
	return news, nil
}

func (s *SharesService) DeleteNotification(ctx context.Context, id int64) (string, string, error) {
	return s.repo.Delete(ctx, id)
}

func (s *SharesService) CreateNotification(ctx context.Context, p model.NewsPostNew) (model.NewsList, error) {
	return s.repo.Create(ctx, p)
}

func (s *SharesService) ToEditNotification(ctx context.Context, e model.NewsEditNotification) (model.NewsList, string, string, error) {
	return s.repo.ToEdit(ctx, e)
}
