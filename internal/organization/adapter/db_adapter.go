package adapter

import (
	"admin/internal/organization"
	"admin/pkg/apperror"
	"database/sql"
	"fmt"
)

type SQLOrganizationRepository struct {
	db *sql.DB
}

func NewSQLOrganizationRepository(db *sql.DB) organization.DBRepository {
	return &SQLOrganizationRepository{db: db}
}

func (repo *SQLOrganizationRepository) FindById(id int) (*organization.Organization, error) {
	query := `SELECT id, name, address, contact_information FROM organizations WHERE id = $1`
	var org organization.Organization
	err := repo.db.QueryRow(query, id).Scan(&org.ID, &org.Name, &org.Address, &org.ContactInformation)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFoundError(fmt.Sprintf("Organization with ID %d not found", id))
		}
		return nil, apperror.NewInternalError(fmt.Sprintf("Error fetching organization by ID: %v", err))
	}

	return &org, nil
}

func (repo *SQLOrganizationRepository) FindAll() ([]*organization.Organization, error) {
	query := `SELECT id, name, address, contact_information FROM organizations`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, apperror.NewInternalError(fmt.Sprintf("Error fetching all organizations: %v", err))
	}
	defer rows.Close()

	var organizations []*organization.Organization
	for rows.Next() {
		var org organization.Organization
		if err := rows.Scan(&org.ID, &org.Name, &org.Address, &org.ContactInformation); err != nil {
			return nil, apperror.NewInternalError(fmt.Sprintf("Error scanning organization: %v", err))
		}
		organizations = append(organizations, &org)
	}
	if err = rows.Err(); err != nil {
		return nil, apperror.NewInternalError(fmt.Sprintf("Error iterating organizations: %v", err))
	}

	return organizations, nil
}
