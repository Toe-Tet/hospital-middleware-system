package patients

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"hospital-middleware-system/src/modules/patients/dto"
	"hospital-middleware-system/src/modules/patients/model"
)

type PatientRepository interface {
	GetByID(ctx context.Context, id int) (*model.Patient, error)
	FindByAnyID(ctx context.Context, hospitalID int, identifier string) (*model.Patient, error)
	List(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) PatientRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*model.Patient, error) {
	query := `
		SELECT id, hospital_id, patient_hn,
			first_name_th, middle_name_th, last_name_th,
			first_name_en, middle_name_en, last_name_en,
			national_id, passport_id,
			date_of_birth, gender, phone_number, status,
			created_at, updated_at
		FROM patients WHERE id = $1
	`
	p := &model.Patient{}
	var (
		firstNameTh  sql.NullString
		middleNameTh sql.NullString
		lastNameTh   sql.NullString
		middleNameEn sql.NullString
		nationalID   sql.NullString
		passportID   sql.NullString
		updatedAt    sql.NullTime
	)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.HospitalID, &p.PatientHN,
		&firstNameTh, &middleNameTh, &lastNameTh,
		&p.FirstNameEn, &middleNameEn, &p.LastNameEn,
		&nationalID, &passportID,
		&p.DateOfBirth, &p.Gender, &p.PhoneNumber, &p.Status,
		&p.CreatedAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if firstNameTh.Valid {
		s := firstNameTh.String
		p.FirstNameTh = &s
	}
	if middleNameTh.Valid {
		s := middleNameTh.String
		p.MiddleNameTh = &s
	}
	if lastNameTh.Valid {
		s := lastNameTh.String
		p.LastNameTh = &s
	}
	if middleNameEn.Valid {
		s := middleNameEn.String
		p.MiddleNameEn = &s
	}
	if nationalID.Valid {
		s := nationalID.String
		p.NationalID = &s
	}
	if passportID.Valid {
		s := passportID.String
		p.PassportID = &s
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		p.UpdatedAt = &t
	}
	return p, nil
}

func (r *postgresRepository) FindByAnyID(ctx context.Context, hospitalID int, identifier string) (*model.Patient, error) {
	var (
		where  []string
		args   []interface{}
		argIdx = 1
	)

	where = append(where, "hospital_id = $"+strconv.Itoa(argIdx))
	args = append(args, hospitalID)
	argIdx++

	where = append(where, "(national_id = $"+strconv.Itoa(argIdx)+" OR passport_id = $"+strconv.Itoa(argIdx+1)+")")
	args = append(args, identifier, identifier)
	argIdx += 2

	query := `
		SELECT id, hospital_id, patient_hn,
			first_name_th, middle_name_th, last_name_th,
			first_name_en, middle_name_en, last_name_en,
			national_id, passport_id,
			date_of_birth, gender, email, phone_number, status,
			created_at, updated_at
		FROM patients
		WHERE ` + strings.Join(where, " AND ") + `
		LIMIT 1
	`
	p := &model.Patient{}
	var (
		firstNameTh  sql.NullString
		middleNameTh sql.NullString
		lastNameTh   sql.NullString
		middleNameEn sql.NullString
		nationalID   sql.NullString
		passportID   sql.NullString
		updatedAt    sql.NullTime
	)
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&p.ID, &p.HospitalID, &p.PatientHN,
		&firstNameTh, &middleNameTh, &lastNameTh,
		&p.FirstNameEn, &middleNameEn, &p.LastNameEn,
		&nationalID, &passportID,
		&p.DateOfBirth, &p.Gender, &p.Email, &p.PhoneNumber, &p.Status,
		&p.CreatedAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if firstNameTh.Valid {
		s := firstNameTh.String
		p.FirstNameTh = &s
	}
	if middleNameTh.Valid {
		s := middleNameTh.String
		p.MiddleNameTh = &s
	}
	if lastNameTh.Valid {
		s := lastNameTh.String
		p.LastNameTh = &s
	}
	if middleNameEn.Valid {
		s := middleNameEn.String
		p.MiddleNameEn = &s
	}
	if nationalID.Valid {
		s := nationalID.String
		p.NationalID = &s
	}
	if passportID.Valid {
		s := passportID.String
		p.PassportID = &s
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		p.UpdatedAt = &t
	}
	return p, nil
}

func (r *postgresRepository) List(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error) {
	var (
		total  int64
		where  []string
		args   []interface{}
		argIdx = 1
	)

	where = append(where, "hospital_id = $"+strconv.Itoa(argIdx))
	args = append(args, hospitalID)
	argIdx++

	if filters.NationalID != "" {
		where = append(where, "national_id = $"+strconv.Itoa(argIdx))
		args = append(args, filters.NationalID)
		argIdx++
	}

	if filters.PassportID != "" {
		where = append(where, "passport_id = $"+strconv.Itoa(argIdx))
		args = append(args, filters.PassportID)
		argIdx++
	}

	if filters.FirstName != "" {
		where = append(where, "(first_name_en ILIKE $"+strconv.Itoa(argIdx)+" OR first_name_th ILIKE $"+strconv.Itoa(argIdx)+")")
		args = append(args, "%"+filters.FirstName+"%")
		argIdx++
	}

	if filters.MiddleName != "" {
		where = append(where, "(middle_name_en ILIKE $"+strconv.Itoa(argIdx)+" OR middle_name_th ILIKE $"+strconv.Itoa(argIdx)+")")
		args = append(args, "%"+filters.MiddleName+"%")
		argIdx++
	}

	if filters.LastName != "" {
		where = append(where, "(last_name_en ILIKE $"+strconv.Itoa(argIdx)+" OR last_name_th ILIKE $"+strconv.Itoa(argIdx)+")")
		args = append(args, "%"+filters.LastName+"%")
		argIdx++
	}

	if filters.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", filters.DateOfBirth)
		if err == nil {
			where = append(where, "date_of_birth = $"+strconv.Itoa(argIdx))
			args = append(args, dob)
			argIdx++
		}
	}

	if filters.Email != "" {
		where = append(where, "email ILIKE $"+strconv.Itoa(argIdx))
		args = append(args, "%"+filters.Email+"%")
		argIdx++
	}

	if filters.PhoneNumber != "" {
		where = append(where, "phone_number ILIKE $"+strconv.Itoa(argIdx))
		args = append(args, "%"+filters.PhoneNumber+"%")
		argIdx++
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	countQuery := "SELECT COUNT(*) FROM patients " + whereClause
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	args = append(args, perPage, offset)
	query := `
		SELECT id, hospital_id, patient_hn,
			first_name_th, middle_name_th, last_name_th,
			first_name_en, middle_name_en, last_name_en,
			national_id, passport_id,
			date_of_birth, gender, email, phone_number, status,
			created_at, updated_at
		FROM patients
		` + whereClause + `
		ORDER BY id DESC
		LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*model.Patient
	for rows.Next() {
		p := &model.Patient{}
		var (
			firstNameTh  sql.NullString
			middleNameTh sql.NullString
			lastNameTh   sql.NullString
			middleNameEn sql.NullString
			nationalID   sql.NullString
			passportID   sql.NullString
			updatedAt    sql.NullTime
		)
		if err := rows.Scan(
			&p.ID, &p.HospitalID, &p.PatientHN,
			&firstNameTh, &middleNameTh, &lastNameTh,
			&p.FirstNameEn, &middleNameEn, &p.LastNameEn,
			&nationalID, &passportID,
			&p.DateOfBirth, &p.Gender, &p.Email, &p.PhoneNumber, &p.Status,
			&p.CreatedAt, &updatedAt,
		); err != nil {
			return nil, 0, err
		}
		if firstNameTh.Valid {
			s := firstNameTh.String
			p.FirstNameTh = &s
		}
		if middleNameTh.Valid {
			s := middleNameTh.String
			p.MiddleNameTh = &s
		}
		if lastNameTh.Valid {
			s := lastNameTh.String
			p.LastNameTh = &s
		}
		if middleNameEn.Valid {
			s := middleNameEn.String
			p.MiddleNameEn = &s
		}
		if nationalID.Valid {
			s := nationalID.String
			p.NationalID = &s
		}
		if passportID.Valid {
			s := passportID.String
			p.PassportID = &s
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			p.UpdatedAt = &t
		}
		list = append(list, p)
	}
	return list, total, rows.Err()
}

func isFKViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "23503")
}
