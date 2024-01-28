package adapter

import (
	"admin/internal/user"
	"admin/pkg/apperror"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type SQLUserRepository struct {
	db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) user.DBRepository {
	return &SQLUserRepository{db: db}
}

func (repo *SQLUserRepository) FindById(id int) (*user.User, error) {
	query := `SELECT id, name, organization_id, email, role, password FROM users WHERE id = $1`
	var u user.User
	err := repo.db.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.OrganizationID, &u.Email, &u.Role, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFoundError(fmt.Sprintf("User with ID %d not found", id))
		}
		return nil, apperror.NewInternalError(fmt.Sprintf("Error fetching user by ID: %v", err))
	}

	return &u, nil
}

func (repo *SQLUserRepository) FindByEmail(email string) (*user.User, error) {
	query := `SELECT id, name, organization_id, email, role, password FROM users WHERE email = $1`
	var u user.User
	err := repo.db.QueryRow(query, email).Scan(&u.ID, &u.Name, &u.OrganizationID, &u.Email, &u.Role, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFoundError(fmt.Sprintf("User with email %s not found", email))
		}
		return nil, apperror.NewInternalError(fmt.Sprintf("Error fetching user by email: %v", err))
	}

	return &u, nil
}

func (repo *SQLUserRepository) FindAll() ([]*user.User, error) {
	query := `SELECT id, name, organization_id, email, role, password FROM users`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, apperror.NewInternalError(fmt.Sprintf("Error fetching all users: %v", err))
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Name, &u.OrganizationID, &u.Email, &u.Role, &u.Password); err != nil {
			return nil, apperror.NewInternalError(fmt.Sprintf("Error scanning user: %v", err))
		}
		users = append(users, &u)
	}
	if err = rows.Err(); err != nil {
		return nil, apperror.NewInternalError(fmt.Sprintf("Error iterating users: %v", err))
	}

	return users, nil
}

func (repo *SQLUserRepository) FindByOrganization(organizationID int) ([]*user.User, error) {
	query := `SELECT id, name, organization_id, email, role, password FROM users WHERE organization_id = $1`
	rows, err := repo.db.Query(query, organizationID)
	if err != nil {
		return nil, apperror.NewInternalError(fmt.Sprintf("Error fetching users by organization ID: %v", err))
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Name, &u.OrganizationID, &u.Email, &u.Role, &u.Password); err != nil {
			return nil, apperror.NewInternalError(fmt.Sprintf("Error scanning user: %v", err))
		}
		users = append(users, &u)
	}
	if err = rows.Err(); err != nil {
		return nil, apperror.NewInternalError(fmt.Sprintf("Error iterating users by organization: %v", err))
	}

	return users, nil
}

func (repo *SQLUserRepository) Save(u *user.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.NewInternalError(fmt.Sprintf("Error hashing password: %v", err))
	}
	u.Password = string(hashedPassword)

	query := `INSERT INTO users (name, organization_id, email, role, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = repo.db.QueryRow(query, u.Name, u.OrganizationID, u.Email, u.Role, u.Password).Scan(&u.ID)
	if err != nil {
		return apperror.NewInternalError(fmt.Sprintf("Error saving user: %v", err))
	}

	return nil
}
