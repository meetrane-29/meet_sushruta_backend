package main

import (
	"fmt"
	"log"

	"meet_sushruta/config"
	"meet_sushruta/handler"
	"meet_sushruta/middleware"
	"meet_sushruta/model"
	"meet_sushruta/repository"
	"meet_sushruta/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	cfg := config.GetConfig()

	// Initialize database
	if err := config.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDatabase()

	// Seed demo data
	if err := config.SeedDatabase(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

	// Initialize auth service and handler
	authService := service.NewAuthService(cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	// Initialize repositories
	patientRepo := repository.NewPatientRepository()
	doctorRepo := repository.NewDoctorRepository()
	nurseRepo := repository.NewNurseRepository()
	doctorScheduleRepo := repository.NewDoctorScheduleRepository()
	appointmentRepo := repository.NewAppointmentRepository()
	billingRepo := repository.NewBillingRepository()
	adminRepo := repository.NewAdminRepository()
	vitalsRepo := repository.NewVitalsRepository()
	specializationRepo := repository.NewSpecializationRepository()
	hospitalRepo := repository.NewHospitalRepository()
	bedRepo := repository.NewBedRepository(config.DB)
	medicalEquipmentRepo := repository.NewMedicalEquipmentRepository(config.DB)
	operationTheatreRepo := repository.NewOperationTheatreRepository(config.DB)
	operationScheduleRepo := repository.NewOperationScheduleRepository(config.DB)

	// Initialize services
	patientService := service.NewPatientService(patientRepo)
	doctorService := service.NewDoctorService(doctorRepo, doctorScheduleRepo, appointmentRepo)
	nurseService := service.NewNurseService(nurseRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo, patientRepo, doctorRepo)
	prescriptionService := service.NewPrescriptionService()
	medicineRepo := repository.NewMedicineRepository()
	pharmacyService := service.NewPharmacyService(medicineRepo)
	labService := service.NewLabService()
	billingService := service.NewBillingService(billingRepo, patientRepo)
	adminService := service.NewAdminService(adminRepo, bedRepo, medicalEquipmentRepo, operationTheatreRepo, operationScheduleRepo)
	vitalsService := service.NewVitalsService(vitalsRepo)
	specializationService := service.NewSpecializationService(specializationRepo)
	hospitalService := service.NewHospitalService(hospitalRepo)

	// Initialize notification service and start worker
	notifService := service.NewNotificationService(config.DB)
	go notifService.StartWorker()

	// Initialize handlers
	patientHandler := handler.NewPatientHandler(patientService)
	nurseHandler := handler.NewNurseHandler(nurseService)
	doctorHandler := handler.NewDoctorHandler(doctorService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService, patientRepo, doctorRepo)
	prescriptionHandler := handler.NewPrescriptionHandler(prescriptionService)
	pharmacyHandler := handler.NewPharmacyHandler(pharmacyService)
	labHandler := handler.NewLabHandler(labService)
	billingHandler := handler.NewBillingHandler(billingService)
	adminHandler := handler.NewAdminHandler(adminService)
	vitalsHandler := handler.NewVitalsHandler(vitalsService)
	specializationHandler := handler.NewSpecializationHandler(specializationService)
	hospitalHandler := handler.NewHospitalHandler(hospitalService)

	// Add routes
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/register", authHandler.Register)
		}

		// Public routes (no authentication required)
		public := v1.Group("")
		{
			// Public doctor routes for home page and doctor directory
			public.GET("/doctors", doctorHandler.GetAllDoctors)
			public.GET("/doctors/:id", doctorHandler.GetDoctor)
			// Public specialization routes
			public.GET("/specializations", specializationHandler.GetAllSpecializations)
			public.GET("/specializations/main", specializationHandler.GetMainSpecializations)
			public.GET("/specializations/:id", specializationHandler.GetSpecialization)
			public.GET("/specializations/:id/sub", specializationHandler.GetSubSpecializations)
		}

		// Protected routes require auth middleware and audit logging
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		protected.Use(middleware.AuditMiddleware(config.DB))
		{
			// Get users by role (for staff management pages)
			protected.GET("/users", middleware.RequireRole("admin"), adminHandler.GetStaffByRole)

			// Admin routes
			admin := protected.Group("/admin")
			{
				admin.GET("/dashboard", middleware.RequireRole("admin"), adminHandler.GetDashboardStats)
				admin.GET("/appointments/today", middleware.RequireRole("admin"), adminHandler.GetTodayAppointmentsByDoctor)
				admin.POST("/doctors/register", middleware.RequireRole("admin"), adminHandler.RegisterDoctor)

				// Inventory routes
				inventory := admin.Group("/inventory")
				{
					// Beds
					inventory.GET("/beds/stats", middleware.RequireRole("admin"), adminHandler.GetBedStats)
					inventory.GET("/beds", middleware.RequireRole("admin"), adminHandler.GetAllBeds)
					inventory.GET("/beds/type/:bedType", middleware.RequireRole("admin"), adminHandler.GetPatientsByBedType)

					// Medical Equipment
					inventory.GET("/equipment/stats", middleware.RequireRole("admin"), adminHandler.GetMedicalEquipmentStats)
					inventory.GET("/equipment", middleware.RequireRole("admin"), adminHandler.GetAllMedicalEquipment)

					// Operation Theatre
					inventory.GET("/ot/stats", middleware.RequireRole("admin"), adminHandler.GetOperationTheatreStats)
					inventory.GET("/operation-schedules", middleware.RequireRole("admin"), adminHandler.GetOperationSchedules)
				}
			}

			// Patient routes
			patients := protected.Group("/patients")
			{
				patients.GET("", middleware.RequireRole("admin", "doctor", "nurse"), patientHandler.GetAllPatients)
				patients.POST("", middleware.RequireRole("admin", "nurse"), patientHandler.CreatePatient)
				patients.GET("/:id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), patientHandler.GetPatient)
				patients.PATCH("/:id", middleware.RequireRole("admin"), patientHandler.UpdatePatient)
				patients.DELETE("/:id", middleware.RequireRole("admin"), patientHandler.DeletePatient)
			}

			// Vitals routes
			vitals := protected.Group("/vitals")
			{
				vitals.GET("/patient/:patient_id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), vitalsHandler.GetPatientVitals)
				vitals.POST("", middleware.RequireRole("admin", "nurse"), vitalsHandler.CreateVitals)
				vitals.GET("/:id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), vitalsHandler.GetVitalsByID)
				vitals.PATCH("/:id", middleware.RequireRole("admin", "nurse"), vitalsHandler.UpdateVitals)
				vitals.DELETE("/:id", middleware.RequireRole("admin"), vitalsHandler.DeleteVitals)
			}

			// Doctor management routes (protected)
			doctors := protected.Group("/doctors")
			{
				doctors.POST("", middleware.RequireRole("admin"), doctorHandler.CreateDoctor)
				doctors.PATCH("/:id", middleware.RequireRole("admin"), doctorHandler.UpdateDoctor)
				doctors.DELETE("/:id", middleware.RequireRole("admin"), doctorHandler.DeleteDoctor)
				doctors.GET("/:id/slots", doctorHandler.GetAvailableSlots)
			}

			// Nurse management routes (protected)
			nurses := protected.Group("/nurses")
			{
				nurses.GET("", middleware.RequireRole("admin"), nurseHandler.GetAllNurses)
				nurses.POST("", middleware.RequireRole("admin"), nurseHandler.CreateNurse)
				nurses.GET("/:id", middleware.RequireRole("admin"), nurseHandler.GetNurse)
				nurses.PATCH("/:id", middleware.RequireRole("admin"), nurseHandler.UpdateNurse)
				nurses.DELETE("/:id", middleware.RequireRole("admin"), nurseHandler.DeleteNurse)
				nurses.GET("/user/:user_id", middleware.RequireRole("admin", "nurse"), nurseHandler.GetNurseByUserID)
			}

			// Appointment routes
			appointments := protected.Group("/appointments")
			{
				appointments.GET("", middleware.RequireRole("admin", "doctor", "nurse", "patient"), appointmentHandler.GetAllAppointments)
				appointments.POST("", middleware.RequireRole("admin", "doctor", "nurse", "patient"), appointmentHandler.BookAppointment)
				appointments.GET("/:id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), appointmentHandler.GetAppointment)
				appointments.GET("/:id/vitals", middleware.RequireRole("admin", "doctor", "nurse", "patient"), appointmentHandler.GetAppointmentVitals)
				appointments.PATCH("/:id/status", middleware.RequireRole("admin", "doctor"), appointmentHandler.UpdateAppointmentStatus)
				appointments.PATCH("/:id/vitals", middleware.RequireRole("nurse"), appointmentHandler.UpdateAppointmentVitals)
			}

			// Prescription routes
			prescriptions := protected.Group("/prescriptions")
			{
				prescriptions.POST("", middleware.RequireRole("admin", "doctor"), prescriptionHandler.CreatePrescription)
				prescriptions.GET("", middleware.RequireRole("admin", "doctor", "nurse", "pharmacist"), prescriptionHandler.ListPrescriptions)
				prescriptions.GET("/patient/:patient_id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), prescriptionHandler.GetPatientPrescriptions)
				prescriptions.GET("/:id", middleware.RequireRole("admin", "doctor", "nurse", "pharmacist", "patient"), prescriptionHandler.GetPrescription)
				prescriptions.GET("/:id/medicines", middleware.RequireRole("admin", "doctor", "nurse", "pharmacist", "patient"), prescriptionHandler.GetPrescriptionMedicines)
				prescriptions.PATCH("/:id/status", middleware.RequireRole("admin", "doctor"), prescriptionHandler.UpdatePrescriptionStatus)
			}

			// Pharmacy routes
			pharmacy := protected.Group("/pharmacy")
			{
				// Medicine CRUD operations
				pharmacy.POST("/medicines", middleware.RequireRole("admin", "pharmacist"), pharmacyHandler.CreateMedicine)
				pharmacy.GET("/medicines", middleware.RequireRole("admin", "doctor", "nurse", "pharmacist"), pharmacyHandler.ListMedicines)

				// Stock management - specific routes first
				pharmacy.GET("/medicines/low-stock", middleware.RequireRole("admin", "pharmacist"), pharmacyHandler.GetLowStockMedicines)
				pharmacy.GET("/medicines/expiring", middleware.RequireRole("admin", "pharmacist"), pharmacyHandler.GetExpiringMedicines)

				// Then generic ID routes
				pharmacy.GET("/medicines/:id", middleware.RequireRole("admin", "doctor", "nurse", "pharmacist"), pharmacyHandler.GetMedicine)
				pharmacy.PATCH("/medicines/:id", middleware.RequireRole("admin", "pharmacist"), pharmacyHandler.UpdateMedicine)
				pharmacy.DELETE("/medicines/:id", middleware.RequireRole("admin", "pharmacist"), pharmacyHandler.DeleteMedicine)
				pharmacy.PATCH("/medicines/:id/stock", middleware.RequireRole("admin", "pharmacist"), pharmacyHandler.UpdateStock)

				// Dispensing
				pharmacy.POST("/dispense", middleware.RequireRole("admin", "pharmacist"), pharmacyHandler.Dispense)
				pharmacy.GET("/dispense/:prescription_id", middleware.RequireRole("admin", "doctor", "pharmacist"), pharmacyHandler.GetDispenseHistory)
			}

			// Lab routes
			lab := protected.Group("/lab")
			{
				lab.POST("/orders", middleware.RequireRole("admin", "doctor"), labHandler.CreateOrder)
				lab.GET("/orders", middleware.RequireRole("admin", "doctor", "nurse"), labHandler.ListOrders)
				lab.GET("/orders/:id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), labHandler.GetOrder)
				lab.PATCH("/orders/:id/status", middleware.RequireRole("admin", "nurse"), labHandler.UpdateStatus)
				lab.POST("/orders/:id/report", middleware.RequireRole("admin", "nurse"), labHandler.UploadReport)
				lab.GET("/patients/:patient_id/orders", middleware.RequireRole("admin", "doctor", "nurse", "patient"), labHandler.GetPatientLabOrders)
			}

			// Billing routes
			billing := protected.Group("/billing")
			{
				billing.GET("", middleware.RequireRole("admin", "doctor"), billingHandler.ListBills)
				billing.POST("/generate", middleware.RequireRole("admin", "doctor", "nurse"), billingHandler.GenerateBill)
				billing.GET("/:id", middleware.RequireRole("admin", "doctor", "patient"), billingHandler.GetBill)
				billing.PATCH("/:id/pay", middleware.RequireRole("admin", "patient"), billingHandler.PayBill)
			}

			// Specialization management routes (admin only)
			specializations := protected.Group("/specializations")
			{
				specializations.POST("", middleware.RequireRole("admin"), specializationHandler.CreateSpecialization)
				specializations.PATCH("/:id", middleware.RequireRole("admin"), specializationHandler.UpdateSpecialization)
				specializations.DELETE("/:id", middleware.RequireRole("admin"), specializationHandler.DeleteSpecialization)
			}

			// Hospital management routes
			hospitals := protected.Group("/hospitals")
			{
				hospitals.GET("", middleware.RequireRole("admin", "doctor"), hospitalHandler.GetAllHospitals)
				hospitals.POST("", middleware.RequireRole("admin"), hospitalHandler.CreateHospital)
				hospitals.GET("/search", middleware.RequireRole("admin", "doctor", "patient"), hospitalHandler.SearchHospitals)
				hospitals.GET("/city/:city", middleware.RequireRole("admin", "doctor", "patient"), hospitalHandler.GetHospitalsByCity)
				hospitals.GET("/:id", middleware.RequireRole("admin", "doctor", "patient"), hospitalHandler.GetHospital)
				hospitals.PATCH("/:id", middleware.RequireRole("admin"), hospitalHandler.UpdateHospital)
				hospitals.DELETE("/:id", middleware.RequireRole("admin"), hospitalHandler.DeleteHospital)
			}
		}
	}

	// Debug routes (development only) - to help with testing
	debug := v1.Group("/debug")
	{
		debug.GET("/users", func(c *gin.Context) {
			var users []model.User
			config.DB.Select("id, email, role, created_at").Find(&users)
			c.JSON(200, gin.H{
				"count": len(users),
				"users": users,
			})
		})

		debug.POST("/clear-users", func(c *gin.Context) {
			resetPasscode := c.Query("code")
			if resetPasscode != "dev-reset-123" {
				c.JSON(403, gin.H{"error": "unauthorized"})
				return
			}

			// Delete non-admin users
			var patients []model.Patient
			config.DB.Find(&patients)
			for _, p := range patients {
				config.DB.Delete(&p)
			}

			var doctors []model.Doctor
			config.DB.Find(&doctors)
			for _, d := range doctors {
				config.DB.Delete(&d)
			}

			config.DB.Where("role IN ?", []string{"patient", "doctor", "nurse", "pharmacy", "lab"}).Delete(&model.User{})

			c.JSON(200, gin.H{"message": "Non-admin users cleared. Restart backend to reseed test data."})
		})
	}

	// Start server
	port := cfg.ServerPort
	log.Printf("Starting server on port %s", port)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
