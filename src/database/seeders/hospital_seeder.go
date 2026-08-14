//go:build seeder
// +build seeder

package main

import (
	"context"
	"database/sql"
	"fmt"
)

type HospitalSeeder struct {
	db *sql.DB
}

func NewHospitalSeeder(db *sql.DB) *HospitalSeeder {
	return &HospitalSeeder{db: db}
}

func (s *HospitalSeeder) Run(ctx context.Context) error {
	hospitals := []struct {
		Name    string
		Address string
		Phone   string
	}{
		{
			Name:    "City General Hospital",
			Address: "123 Main Street, Downtown",
			Phone:   "+1-555-0101",
		},
		{
			Name:    "Sunrise Medical Center",
			Address: "456 Oak Avenue, Midtown",
			Phone:   "+1-555-0102",
		},
		{
			Name:    "Green Valley Clinic",
			Address: "789 Pine Road, Suburbs",
			Phone:   "+1-555-0103",
		},
	}

	for _, h := range hospitals {
		var id int
		err := s.db.QueryRowContext(ctx, `
			INSERT INTO hospitals (name, address, phone, status)
			VALUES ($1, $2, $3, 'active')
			ON CONFLICT (name) DO NOTHING
			RETURNING id
		`, h.Name, h.Address, h.Phone).Scan(&id)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("seed hospital %s: %w", h.Name, err)
		}
		if err == sql.ErrNoRows {
			fmt.Printf("  [SKIP] Hospital: %s (exists)\n", h.Name)
		} else {
			fmt.Printf("  [ OK ] Hospital: %s (id=%d)\n", h.Name, id)
		}
	}
	return nil
}
