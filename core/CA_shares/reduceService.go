package CA_shares

type NewsRepository interface {

}

type ReduceService struct {
	repo NewsRepository
}

func NewReduceService(repo NewsRepository) *ReduceService {
	return &ReduceService{repo: repo}
}