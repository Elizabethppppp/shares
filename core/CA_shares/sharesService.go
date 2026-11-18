package CA_shares

import (
	"CA_shares/filtres"
	"CA_shares/logger"
	"CA_shares/model"
	"CA_shares/serviceErrors"
	"CA_shares/storage"
	"CA_shares/validation"
	"context"
	"time"
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
	st   *storage.ImageManager
}

func NewSharesService(repo NewsRepository, st *storage.ImageManager) *SharesService {
	return &SharesService{repo: repo, st: st}
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

func (s *SharesService) GetAllNews(ctx context.Context, q model.ListQuery) (model.NewsListResponse, error) {

	if q.Offset < 0 || q.Limit < 1 || q.Limit > 100 || len([]rune(q.Title)) > 50 {
		return model.NewsListResponse{}, serviceErrors.ErrBadRequest
	}

	sortBy := "published_at"
	if q.SortBy != "" {
		switch q.SortBy {
		case "publishedAt":
			sortBy = "published_at"
		case "id":
			sortBy = "id"
		default:
			return model.NewsListResponse{}, serviceErrors.ErrBadRequest
		}
	}

	sortDir := "desc"
	if q.SortDir != "" {
		switch q.SortDir {
		case "asc", "desc":
			sortDir = q.SortDir
		default:
			return model.NewsListResponse{}, serviceErrors.ErrBadRequest
		}
	}

	var title *string
	if q.Title != "" {
		title = &q.Title
	}

	filter := filtres.ListFilter{
		Offset:   q.Offset,
		Limit:    q.Limit,
		Title:    title,
		IsActive: q.IsActive,
		SortBy:   sortBy,
		SortDir:  sortDir,
	}

	items, total, err := s.repo.GetAllList(ctx, filter)
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

func (s *SharesService) DeleteNotification(ctx context.Context, id int64) error {
	previewImg, img, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if err := s.st.Delete(previewImg); err != nil {
		logger.Error("delete preview", "key", previewImg, "error", err.Error())
	}
	if err := s.st.Delete(img); err != nil {
		logger.Error("delete image", "key", img, "error", err.Error())
	}
	return nil
}

func (s *SharesService) CreateNotification(ctx context.Context, in model.NewsCreateInput) (model.NewsList, error) {
	if in.Title == "" || len([]rune(in.Title)) > 50 {
		return model.NewsList{}, serviceErrors.ErrNewsValidationFailed
	}
	if in.Content == "" || len([]rune(in.Content)) > 10000 {
		return model.NewsList{}, serviceErrors.ErrNewsValidationFailed
	}
	publishedAt, err := time.Parse("02-01-2006", in.PublishedAt)
	if err != nil {
		return model.NewsList{}, serviceErrors.ErrNewsValidationFailed
	}

	var previewKey, imgKey string
	if in.PreviewImg != nil {
		reader, err := validation.ValidateImgFile(in.PreviewImg.Filename, in.PreviewImg.Size, in.PreviewImg.Reader, 5<<20)
		if err != nil {
			return model.NewsList{}, serviceErrors.ErrNewsValidationFailed
		}
		key, err := s.st.Save("previewImg", in.PreviewImg.Filename, reader)
		if err != nil {
			logger.Error("save preview", "error", err.Error())
			return model.NewsList{}, serviceErrors.ErrInternal
		}
		previewKey = key
	}
	if in.Img != nil {
		reader, err := validation.ValidateImgFile(in.Img.Filename, in.Img.Size, in.Img.Reader, 10<<20)
		if err != nil {
			if previewKey != "" {
				_ = s.st.Delete(previewKey)
			}
			return model.NewsList{}, serviceErrors.ErrNewsValidationFailed
		}
		key, err := s.st.Save("img", in.Img.Filename, reader)
		if err != nil {
			if previewKey != "" {
				_ = s.st.Delete(previewKey)
			}
			logger.Error("save img", "error", err.Error())
			return model.NewsList{}, serviceErrors.ErrInternal
		}
		imgKey = key
	}
	return s.repo.Create(ctx, model.NewsPostNew{
		Title:       in.Title,
		Content:     in.Content,
		PublishedAt: publishedAt,
		PreviewImg:  previewKey,
		Img:         imgKey,
	})
}

func (s *SharesService) ToEditNotification(ctx context.Context, up model.NewsUpdateInput) (model.NewsList, error) {
	if up.Title == "" || len([]rune(up.Title)) > 50 {
		return model.NewsList{}, serviceErrors.ErrNewsUpdateFailed
	}
	if up.Content == "" || len([]rune(up.Content)) > 10000 {
		return model.NewsList{}, serviceErrors.ErrNewsUpdateFailed
	}
	publishedAt, err := time.Parse("02-01-2006", up.PublishedAt)
	if err != nil {
		return model.NewsList{}, serviceErrors.ErrNewsUpdateFailed
	}

	var newPreviewImg, newImg *string
	if up.PreviewImg != nil {
		reader, err := validation.ValidateImgFile(up.PreviewImg.Filename, up.PreviewImg.Size, up.PreviewImg.Reader, 5<<20)
		if err != nil {
			return model.NewsList{}, serviceErrors.ErrNewsUpdateFailed
		}
		key, err := s.st.Save("previewImg", up.PreviewImg.Filename, reader)
		if err != nil {
			logger.Error("save preview", "error", err.Error())
			return model.NewsList{}, serviceErrors.ErrInternal
		}
		newPreviewImg = &key
	}
	if up.Img != nil {
		reader, err := validation.ValidateImgFile(up.Img.Filename, up.Img.Size, up.Img.Reader, 10<<20)
		if err != nil {
			if newPreviewImg != nil {
				_ = s.st.Delete(*newPreviewImg)
			}
			return model.NewsList{}, serviceErrors.ErrNewsUpdateFailed
		}
		key, err := s.st.Save("img", up.Img.Filename, reader)
		if err != nil {
			if newPreviewImg != nil {
				_ = s.st.Delete(*newPreviewImg)
			}
			logger.Error("save img", "error", err.Error())
			return model.NewsList{}, serviceErrors.ErrInternal
		}
		newImg = &key
	}

	news, oldPreviewImg, oldImg, err := s.repo.ToEdit(ctx, model.NewsEditNotification{
		Id:          up.Id,
		Title:       up.Title,
		Content:     up.Content,
		PublishedAt: publishedAt,
		PreviewImg:  newPreviewImg,
		Img:         newImg,
	})
	if err != nil {
		if newPreviewImg != nil {
			_ = s.st.Delete(*newPreviewImg)
		}
		if newImg != nil {
			_ = s.st.Delete(*newImg)
		}
		return model.NewsList{}, err
	}

	if err := s.st.Delete(oldPreviewImg); err != nil {
		logger.Error("delete old preview", "key", oldPreviewImg, "error", err.Error())
	}
	if err := s.st.Delete(oldImg); err != nil {
		logger.Error("delete old img", "key", oldImg, "error", err.Error())
	}
	return news, nil

}
