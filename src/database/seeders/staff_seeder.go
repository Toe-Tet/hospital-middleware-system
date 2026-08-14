//go:build seeder
// +build seeder

package main

import (
	"context"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type StaffSeeder struct {
	db *sql.DB
}

func NewStaffSeeder(db *sql.DB) *StaffSeeder {
	return &StaffSeeder{db: db}
}

func (s *StaffSeeder) Run(ctx context.Context) error {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	staffs := []struct {
		HospitalID int
		Username   string
		Name       string
	}{
		{1, "admin", "Alice Admin"},
		{1, "dr-smith", "Dr. Bob Smith"},
		{1, "carol", "Nurse Carol"},
		{1, "dave", "Dave Reception"},
		{2, "dr-johnson", "Dr. Eve Johnson"},
		{2, "nurse-frank", "Nurse Frank"},
	}

	for _, st := range staffs {
		var id int
		err := s.db.QueryRowContext(ctx, `
			INSERT INTO staffs (hospital_id, username, name, password, status)
			VALUES ($1, $2, $3, $4, 'active')
			ON CONFLICT (username) DO NOTHING
			RETURNING id
		`, st.HospitalID, st.Username, st.Name,  string(pwdHash)).Scan(&id)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("seed staff %s: %w", st.Username, err)
		}
		if err == sql.ErrNoRows {
			fmt.Printf("  [SKIP] Staff: %s (exists)\n", st.Username)
		} else {
			fmt.Printf("  [ OK ] Staff: %s (id=%d)\n", st.Username, id)
		}
	}
	return nil
}
