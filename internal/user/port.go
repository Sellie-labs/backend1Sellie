package user

type DBRepository interface {
	FindById(id int) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]*User, error)
	FindByOrganization(id int) ([]*User, error)
	Save(user *User) error
}
