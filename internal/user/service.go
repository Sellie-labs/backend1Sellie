package user

type Service struct {
	repo DBRepository
}

// NewService creates a new user service with the given repository.
func NewService(repo DBRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetUserByID retrieves a user by their ID.
func (s *Service) GetUserByID(id int) (*User, error) {
	return s.repo.FindById(id)
}

// GetUserByEmail retrieves a user by their email.
func (s *Service) GetUserByEmail(email string) (*User, error) {
	return s.repo.FindByEmail(email)
}

// GetAllUsers retrieves all users.
func (s *Service) GetAllUsers() ([]*User, error) {
	return s.repo.FindAll()
}

// GetUsersByOrganization retrieves all users belonging to a specific organization.
func (s *Service) GetUsersByOrganization(organizationID int) ([]*User, error) {
	return s.repo.FindByOrganization(organizationID)
}

// CreateUser saves a new user to the database.
func (s *Service) CreateUser(user *User) error {
	// TODO: add any pre-save logic, like hashing the user's password
	return s.repo.Save(user)
}
