package organization

type DBRepository interface {
	FindById(id int) (*Organization, error)
	FindAll() ([]*Organization, error)
}
