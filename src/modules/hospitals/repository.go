package hospitals

import (
	"context"
	"database/sql"
	"strings"

	"hospital-middleware-system/src/modules/hospitals/model"
)

type HospitalRepository interface {
	List(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) HospitalRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) List(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM hospitals`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `
		SELECT id, name, address, phone, status, created_at, updated_at
		FROM hospitals
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*model.Hospital
	for rows.Next() {
		h := &model.Hospital{}
		if err := rows.Scan(&h.ID, &h.Name, &h.Address, &h.Phone, &h.Status, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, h)
	}
	return list, total, rows.Err()
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "23505")
}
