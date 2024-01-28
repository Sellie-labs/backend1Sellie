package organization

type Service struct {
	repo DBRepository
}

func NewService(repo DBRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetUserByID retrieves a user by their ID.
func (s *Service) GetOrganizationByID(id int) (*Organization, error) {
	return s.repo.FindById(id)
}

// GetAllUsers retrieves all users.
func (s *Service) GetAllOrganizations() ([]*Organization, error) {
	return s.repo.FindAll()
}
