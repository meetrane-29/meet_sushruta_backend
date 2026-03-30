package main

import (
	"fmt"
	"log"
	"meet_sushruta/config"
	"meet_sushruta/model"
)

// DebugDatabase prints database statistics for troubleshooting
func DebugDatabase() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := config.InitDatabase(config.GetConfig()); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	var appointmentCount int64
	var patientCount int64
	var doctorCount int64
	var userCount int64

	config.DB.Model(&model.Appointment{}).Count(&appointmentCount)
	config.DB.Model(&model.Patient{}).Count(&patientCount)
	config.DB.Model(&model.Doctor{}).Count(&doctorCount)
	config.DB.Model(&model.User{}).Count(&userCount)

	fmt.Printf("\n=== DATABASE STATUS ===\n")
	fmt.Printf("Appointments: %d\n", appointmentCount)
	fmt.Printf("Patients: %d\n", patientCount)
	fmt.Printf("Doctors: %d\n", doctorCount)
	fmt.Printf("Users: %d\n", userCount)

	// Print first 5 users
	var users []model.User
	config.DB.Limit(5).Find(&users)
	fmt.Printf("\nFirst 5 Users:\n")
	for i, user := range users {
		fmt.Printf("%d. Email: %s, Role: %s\n", i+1, user.Email, user.Role)
	}

	if appointmentCount > 0 {
		var appointments []model.Appointment
		config.DB.
			Preload("Patient").
			Preload("Patient.User").
			Preload("Doctor").
			Preload("Doctor.User").
			Limit(5).
			Order("created_at DESC").
			Find(&appointments)

		fmt.Printf("\nFirst 5 Appointments:\n")
		for i, apt := range appointments {
			patientName := "N/A"
			doctorName := "N/A"
			if apt.Patient != nil && apt.Patient.User != nil {
				patientName = apt.Patient.User.FirstName + " " + apt.Patient.User.LastName
			}
			if apt.Doctor != nil && apt.Doctor.User != nil {
				doctorName = apt.Doctor.User.FirstName + " " + apt.Doctor.User.LastName
			}
			fmt.Printf("%d. Patient: %s, Doctor: %s, Status: %s, Date: %s\n",
				i+1, patientName, doctorName, apt.Status, apt.AppointmentDate)
		}
	}
}
