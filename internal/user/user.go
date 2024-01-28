package user

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	OrganizationID int    `json:"organization_id"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Role           string `json:"role"`
}

func New(name string, organizationID int, email string, password string, role string) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	return User{
		Name:           name,
		OrganizationID: organizationID,
		Email:          email,
		Password:       string(hashedPassword),
		Role:           role,
	}, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) HahPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
