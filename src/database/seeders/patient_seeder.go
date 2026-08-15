//go:build seeder
// +build seeder

package main

import (
	"context"
	"database/sql"
	"fmt"
)

type PatientSeeder struct {
	db *sql.DB
}

func NewPatientSeeder(db *sql.DB) *PatientSeeder {
	return &PatientSeeder{db: db}
}

func (s *PatientSeeder) Run(ctx context.Context) error {
	patients := []struct {
		HospitalID   int
		FirstNameTH  *string
		MiddleNameTH *string
		LastNameTH   *string
		FirstNameEN  string
		MiddleNameEN *string
		LastNameEN   string
		PatientHN    string
		PassportID   *string
		NationalID   string
		Gender       string
		Email        string
		PhoneNumber  string
		DateOfBirth  string
	}{
		{
			HospitalID:  1,
			FirstNameEN: "John",
			LastNameEN:  "Doe",
			PatientHN:   "HN000001",
			NationalID:  "1100000000001",
			Gender:      "M",
			Email:       "john@example.com",
			PhoneNumber: "+1-555-1001",
			DateOfBirth: "1985-03-15",
		},
		{
			HospitalID:  1,
			FirstNameEN: "Jane",
			LastNameEN:  "Smith",
			PatientHN:   "HN000002",
			NationalID:  "1100000000002",
			Gender:      "F",
			Email:       "jane@example.com",
			PhoneNumber: "+1-555-1002",
			DateOfBirth: "1990-07-22",
		},
		{
			HospitalID:  1,
			FirstNameEN: "Robert",
			LastNameEN:  "Brown",
			PatientHN:   "HN000003",
			NationalID:  "1100000000003",
			Gender:      "M",
			Email:       "robert@example.com",
			PhoneNumber: "+1-555-1003",
			DateOfBirth: "1978-11-05",
		},
		{
			HospitalID:  1,
			FirstNameEN: "Emily",
			LastNameEN:  "Davis",
			PatientHN:   "HN000004",
			NationalID:  "1100000000004",
			Gender:      "F",
			Email:       "emily@example.com",
			PhoneNumber: "+1-555-1004",
			DateOfBirth: "1995-01-30",
		},
		{
			HospitalID:  2,
			FirstNameEN: "Michael",
			LastNameEN:  "Wilson",
			PatientHN:   "HN000005",
			NationalID:  "1100000000005",
			Gender:      "M",
			Email:       "michael@example.com",
			PhoneNumber: "+1-555-2001",
			DateOfBirth: "1982-09-10",
		},
		{
			HospitalID:  2,
			FirstNameEN: "Sarah",
			LastNameEN:  "Miller",
			PatientHN:   "HN000006",
			NationalID:  "1100000000006",
			Gender:      "F",
			Email:       "sarah@example.com",
			PhoneNumber: "+1-555-2002",
			DateOfBirth: "1998-04-18",
		},
		{
			HospitalID:  3,
			FirstNameEN: "David",
			LastNameEN:  "Taylor",
			PatientHN:   "HN000007",
			NationalID:  "1100000000007",
			Gender:      "M",
			Email:       "david@example.com",
			PhoneNumber: "+1-555-3001",
			DateOfBirth: "1970-12-25",
		},
		{
			HospitalID:  3,
			FirstNameEN: "Lisa",
			LastNameEN:  "Anderson",
			PatientHN:   "HN000008",
			NationalID:  "1100000000008",
			Gender:      "F",
			Email:       "lisa@example.com",
			PhoneNumber: "+1-555-3002",
			DateOfBirth: "1992-06-08",
		},
	}

	for _, p := range patients {
		var id int

		err := s.db.QueryRowContext(ctx, `
			INSERT INTO patients (
				hospital_id,
				first_name_th,
				middle_name_th,
				last_name_th,
				first_name_en,
				middle_name_en,
				last_name_en,
				patient_hn,
				passport_id,
				national_id,
				date_of_birth,
				gender,
        email,
				phone_number,
				status
			)
			VALUES (
				$1, $2, $3, $4,
				$5, $6, $7,
				$8, $9, $10,
				$11::date,
				$12,
        $13,
				$14,
				'active'
			)
			RETURNING id
		`,
			p.HospitalID,
			p.FirstNameTH,
			p.MiddleNameTH,
			p.LastNameTH,
			p.FirstNameEN,
			p.MiddleNameEN,
			p.LastNameEN,
			p.PatientHN,
			p.PassportID,
			p.NationalID,
			p.DateOfBirth,
			p.Gender,
        p.Email,
			p.PhoneNumber,
		).Scan(&id)

		if err != nil {
			return fmt.Errorf(
				"seed patient %s %s: %w",
				p.FirstNameEN,
				p.LastNameEN,
				err,
			)
		}

		fmt.Printf(
			"  [ OK ] Patient: %s %s (id=%d, hn=%s, hospital=%d)\n",
			p.FirstNameEN,
			p.LastNameEN,
			id,
			p.PatientHN,
			p.HospitalID,
		)
	}

	return nil
}
