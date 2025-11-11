package user

type iRepository interface {
}

type Service struct {
	repository iRepository
}

func New(repository iRepository) *Service {
	return &Service{
		repository: repository,
	}
}
