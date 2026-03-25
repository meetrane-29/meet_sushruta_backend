package main

import (
	"fmt"
	"log"

	"meet_sushruta/config"
	"meet_sushruta/handler"
	"meet_sushruta/middleware"
	"meet_sushruta/repository"
	"meet_sushruta/service"

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

	// Create Gin router
	router := gin.Default()

	// Initialize auth service and handler
	authService := service.NewAuthService(cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	// Initialize repositories
	patientRepo := repository.NewPatientRepository()
	doctorRepo := repository.NewDoctorRepository()
	doctorScheduleRepo := repository.NewDoctorScheduleRepository()
	appointmentRepo := repository.NewAppointmentRepository()

	// Initialize services
	patientService := service.NewPatientService(patientRepo)
	doctorService := service.NewDoctorService(doctorRepo, doctorScheduleRepo, appointmentRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo, patientRepo, doctorRepo)
	prescriptionService := service.NewPrescriptionService()
	pharmacyService := service.NewPharmacyService()
	labService := service.NewLabService()

	// Initialize handlers
	patientHandler := handler.NewPatientHandler(patientService)
	doctorHandler := handler.NewDoctorHandler(doctorService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)
	prescriptionHandler := handler.NewPrescriptionHandler(prescriptionService)
	pharmacyHandler := handler.NewPharmacyHandler(pharmacyService)
	labHandler := handler.NewLabHandler(labService)

	// Add routes
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected routes require auth middleware
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Patient routes
			patients := protected.Group("/patients")
			{
				patients.GET("", middleware.RequireRole("admin", "doctor", "nurse"), patientHandler.GetAllPatients)
				patients.POST("", middleware.RequireRole("admin", "nurse"), patientHandler.CreatePatient)
				patients.GET("/:id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), patientHandler.GetPatient)
				patients.PATCH("/:id", middleware.RequireRole("admin"), patientHandler.UpdatePatient)
				patients.DELETE("/:id", middleware.RequireRole("admin"), patientHandler.DeletePatient)
			}

			// Doctor routes
			doctors := protected.Group("/doctors")
			{
				doctors.GET("", doctorHandler.GetAllDoctors)
				doctors.POST("", middleware.RequireRole("admin"), doctorHandler.CreateDoctor)
				doctors.GET("/:id", doctorHandler.GetDoctor)
				doctors.PATCH("/:id", middleware.RequireRole("admin"), doctorHandler.UpdateDoctor)
				doctors.GET("/:id/slots", doctorHandler.GetAvailableSlots)
			}

			// Appointment routes
			appointments := protected.Group("/appointments")
			{
				appointments.GET("", middleware.RequireRole("admin", "doctor", "nurse"), appointmentHandler.GetAllAppointments)
				appointments.POST("", middleware.RequireRole("admin", "doctor", "nurse", "patient"), appointmentHandler.BookAppointment)
				appointments.GET("/:id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), appointmentHandler.GetAppointment)
				appointments.PATCH("/:id/status", middleware.RequireRole("admin", "doctor"), appointmentHandler.UpdateAppointmentStatus)
				appointments.PATCH("/:id/vitals", middleware.RequireRole("nurse"), appointmentHandler.UpdateAppointmentVitals)
			}

			// Prescription routes
			prescriptions := protected.Group("/prescriptions")
			{
				prescriptions.POST("", middleware.RequireRole("admin", "doctor"), prescriptionHandler.CreatePrescription)
				prescriptions.GET("", middleware.RequireRole("admin", "doctor", "nurse", "pharmacist"), prescriptionHandler.ListPrescriptions)
				prescriptions.GET("/:id", middleware.RequireRole("admin", "doctor", "nurse", "pharmacist", "patient"), prescriptionHandler.GetPrescription)
				prescriptions.PATCH("/:id/status", middleware.RequireRole("admin", "doctor"), prescriptionHandler.UpdatePrescriptionStatus)
				prescriptions.GET("/patient/:patient_id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), prescriptionHandler.GetPatientPrescriptions)
			}

			// Pharmacy routes
			pharmacy := protected.Group("/pharmacy")
			{
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
		}
	}

	// Start server
	port := cfg.ServerPort
	log.Printf("Starting server on port %s", port)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
