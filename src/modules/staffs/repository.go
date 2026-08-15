package staffs

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"hospital-middleware-system/src/modules/staffs/model"
)

type StaffRepository interface {
	Create(ctx context.Context, s *model.Staff) (*model.Staff, error)
	GetByID(ctx context.Context, id int) (*model.Staff, error)
	GetByUsername(ctx context.Context, username string) (*model.Staff, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) StaffRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, s *model.Staff) (*model.Staff, error) {
	query := `
		INSERT INTO staffs (hospital_id, username, name, password, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'active', NOW(), NULL)
		RETURNING id, hospital_id, username, name, password, status, created_at, updated_at
	`
  
	result := &model.Staff{}

	err := r.db.QueryRowContext(ctx, query, s.HospitalID, s.Username, s.Name, s.Password).Scan(
		&result.ID, &result.HospitalID, &result.Username, &result.Name, &result.Password, &result.Status, &result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		if isUniqueViolation(err) {
			if strings.Contains(err.Error(), "idx_staffs_email_unique") || strings.Contains(strings.ToLower(err.Error()), "email") {
				return nil, errors.New("staff with this email already exists")
			}
			return nil, errors.New("staff with this username already exists")
		}
		if isFKViolation(err) {
			return nil, errors.New("hospital does not exist")
		}
		return nil, err
	}
	return result, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*model.Staff, error) {
	query := `
		SELECT id, hospital_id, username, name, password, status, created_at, updated_at
		FROM staffs WHERE id = $1
	`
	s := &model.Staff{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.HospitalID, &s.Username, &s.Name, &s.Password, &s.Status, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *postgresRepository) GetByUsername(ctx context.Context, username string) (*model.Staff, error) {
	query := `
		SELECT id, hospital_id, username, name, password, status, created_at, updated_at
		FROM staffs WHERE LOWER(username) = LOWER($1) AND status = 'active'
	`
	s := &model.Staff{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&s.ID, &s.HospitalID, &s.Username, &s.Name, &s.Password, &s.Status, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

type nullStringScan struct {
	dest *string
}

func (n *nullStringScan) Scan(src interface{}) error {
	if src == nil {
		*n.dest = ""
		return nil
	}
	ns := sql.NullString{}
	if err := ns.Scan(src); err != nil {
		return err
	}
	if ns.Valid {
		*n.dest = ns.String
	} else {
		*n.dest = ""
	}
	return nil
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "23505")
}

func isFKViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "23503")
}
