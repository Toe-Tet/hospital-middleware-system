//go:build seeder
// +build seeder

package main

import (
	"context"
	"fmt"
	"os"

	"hospital-middleware-system/src/config"
	"hospital-middleware-system/src/database"
)

func main() {
	if err := config.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}

	if err := database.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "DB error: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	ctx := context.Background()

	fmt.Println("==> Seeding hospitals...")
	if err := NewHospitalSeeder(database.DB).Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Hospital seeder failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("==> Seeding staffs...")
	if err := NewStaffSeeder(database.DB).Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Staff seeder failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("==> Seeding patients...")
	if err := NewPatientSeeder(database.DB).Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Patient seeder failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("==> All seeders completed successfully!")
	fmt.Println("\nDefault login credentials:")
	fmt.Println("  Email:    admin@hospital.com")
	fmt.Println("  Password: password")
}
