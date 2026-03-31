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

	// DEBUG: Print Hardi doctor and his appointments
	debugHardiDoctor()

	// Note: Seed database disabled - use direct SQL file instead
	// if err := config.SeedDatabase(); err != nil {
	//	log.Fatalf("Failed to seed database: %v", err)
	// }

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
	admissionRepo := repository.NewAdmissionRepository()
	progressNoteRepo := repository.NewProgressNoteRepository()
	nurseInstructionRepo := repository.NewNurseInstructionRepository()
	dischargeSummaryRepo := repository.NewDischargeSummaryRepository()
	ratingRepo := repository.NewRatingRepository()

	// New OPD-related repositories
	uhidRepo := repository.NewUHIDRepository(config.DB)
	insurancePolicyRepo := repository.NewInsurancePolicyRepository(config.DB)
	waitingListRepo := repository.NewWaitingListRepository(config.DB)
	opdReceiptRepo := repository.NewOPDReceiptRepository(config.DB)

	// Initialize notification service FIRST - needed by other services
	notifService := service.NewNotificationService(config.DB)
	go notifService.StartWorker()

	// Initialize services
	patientService := service.NewPatientService(patientRepo)
	doctorService := service.NewDoctorService(doctorRepo, doctorScheduleRepo, appointmentRepo)
	nurseService := service.NewNurseService(nurseRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo, patientRepo, doctorRepo)
	prescriptionService := service.NewPrescriptionService(notifService)
	medicineRepo := repository.NewMedicineRepository()
	pharmacyService := service.NewPharmacyService(medicineRepo, notifService)
	labService := service.NewLabService(notifService)
	billingService := service.NewBillingService(billingRepo, patientRepo)
	adminService := service.NewAdminService(adminRepo, bedRepo, medicalEquipmentRepo, operationTheatreRepo, operationScheduleRepo)
	vitalsService := service.NewVitalsService(vitalsRepo)
	bedService := service.NewBedService(bedRepo)
	specializationService := service.NewSpecializationService(specializationRepo)
	hospitalService := service.NewHospitalService(hospitalRepo)

	// Performance optimization & analytics services
	// cacheService := service.NewCacheService(30 * time.Minute) // TODO: Inject into handlers when ready
	exportService := service.NewExportService()
	analyticsService := service.NewAnalyticsService()

	// New OPD-related services
	uhidService := service.NewUHIDService(uhidRepo)
	insurancePolicyService := service.NewInsurancePolicyService(insurancePolicyRepo)
	waitingListService := service.NewWaitingListService(waitingListRepo, appointmentRepo)
	opdReceiptService := service.NewOPDReceiptService(opdReceiptRepo)
	ratingService := service.NewRatingService(ratingRepo, patientRepo, doctorRepo)

	// Initialize handlers
	patientHandler := handler.NewPatientHandler(patientService, appointmentService)
	nurseHandler := handler.NewNurseHandler(nurseService)
	doctorHandler := handler.NewDoctorHandler(doctorService, admissionRepo, progressNoteRepo, nurseInstructionRepo, dischargeSummaryRepo)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService, patientRepo, doctorRepo)
	prescriptionHandler := handler.NewPrescriptionHandler(prescriptionService, pharmacyService, labService)
	pharmacyHandler := handler.NewPharmacyHandler(pharmacyService)
	labHandler := handler.NewLabHandler(labService)
	billingHandler := handler.NewBillingHandler(billingService)
	adminHandler := handler.NewAdminHandler(adminService)
	vitalsHandler := handler.NewVitalsHandler(vitalsService)
	bedHandler := handler.NewBedHandler(bedService)
	admissionHandler := handler.NewAdmissionHandler(admissionRepo)
	specializationHandler := handler.NewSpecializationHandler(specializationService)
	hospitalHandler := handler.NewHospitalHandler(hospitalService)

	// Export and Analytics handlers
	exportHandler := handler.NewExportHandler(exportService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// New OPD-related handlers
	uhidHandler := handler.NewUHIDHandler(uhidService)
	insuranceHandler := handler.NewInsurancePolicyHandler(insurancePolicyService)
	waitingListHandler := handler.NewWaitingListHandler(waitingListService)
	opdReceiptHandler := handler.NewOPDReceiptHandler(opdReceiptService)
	ratingHandler := handler.NewRatingHandler(ratingService)

	// Add routes
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/register", authHandler.Register)
			auth.POST("/admin-register", middleware.RequireRole("admin"), authHandler.AdminRegisterUser)
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
			protected.PATCH("/users/:id", middleware.RequireRole("admin"), adminHandler.UpdateStaffUser)
			protected.DELETE("/users/:id", middleware.RequireRole("admin"), adminHandler.DeleteStaffUser)

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
				patients.GET("", patientHandler.GetAllPatients)
				patients.POST("", middleware.RequireRole("admin", "nurse"), patientHandler.CreatePatient)
				patients.GET("/:id", patientHandler.GetPatient)
				patients.GET("/:id/appointments", patientHandler.GetPatientAppointments)
				patients.PATCH("/:id", middleware.RequireRole("admin"), patientHandler.UpdatePatient)
				patients.DELETE("/:id", middleware.RequireRole("admin"), patientHandler.DeletePatient)
			}

			// Vitals routes
			vitals := protected.Group("/vitals")
			{
				vitals.GET("/patient/:patient_id", vitalsHandler.GetPatientVitals)
				vitals.POST("", middleware.RequireRole("admin", "nurse"), vitalsHandler.CreateVitals)
				vitals.GET("/:id", vitalsHandler.GetVitalsByID)
				vitals.PATCH("/:id", middleware.RequireRole("admin", "nurse"), vitalsHandler.UpdateVitals)
				vitals.DELETE("/:id", middleware.RequireRole("admin"), vitalsHandler.DeleteVitals)
				vitals.GET("/patient/:patient_id/trend", vitalsHandler.GetVitalsTrend)
				vitals.GET("/patient/:patient_id/alerts", middleware.RequireRole("admin", "doctor", "nurse"), vitalsHandler.GetVitalsAlerts)
			}

			// Doctor management routes (protected)
			doctors := protected.Group("/doctors")
			{
				// Get current doctor info - must come BEFORE /:id routes
				doctors.GET("/me", middleware.RequireRole("doctor"), doctorHandler.GetMe)

				doctors.POST("", middleware.RequireRole("admin"), doctorHandler.CreateDoctor)
				doctors.PATCH("/:id", middleware.RequireRole("admin"), doctorHandler.UpdateDoctor)
				doctors.DELETE("/:id", middleware.RequireRole("admin"), doctorHandler.DeleteDoctor)
				doctors.GET("/:id/slots", doctorHandler.GetAvailableSlots)
				doctors.GET("/analytics/my", middleware.RequireRole("doctor"), doctorHandler.GetMyAnalytics)
				doctors.GET("/:id/analytics", middleware.RequireRole("admin", "doctor"), doctorHandler.GetDoctorAnalytics)
				doctors.GET("/:id/performance", middleware.RequireRole("admin", "doctor"), doctorHandler.GetDoctorPerformanceMetrics)
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
				appointments.GET("", appointmentHandler.GetAllAppointments)
				appointments.POST("", middleware.RequireRole("admin", "doctor", "nurse", "patient", "receptionist"), appointmentHandler.BookAppointment)
				appointments.GET("/today", middleware.RequireRole("admin", "receptionist"), appointmentHandler.GetTodayAppointments)
				appointments.GET("/my/schedule", middleware.RequireRole("doctor"), appointmentHandler.GetMyAppointments)
				appointments.GET("/next-7-days", middleware.RequireRole("admin", "nurse", "doctor"), appointmentHandler.GetNext7DaysAppointments)
				appointments.GET("/:id", appointmentHandler.GetAppointment)
				appointments.GET("/:id/vitals", middleware.RequireRole("admin", "doctor", "nurse", "patient"), appointmentHandler.GetAppointmentVitals)
				appointments.PATCH("/:id/status", middleware.RequireRole("admin", "doctor", "receptionist"), appointmentHandler.UpdateAppointmentStatus)
				appointments.PATCH("/:id/vitals", middleware.RequireRole("admin", "nurse"), appointmentHandler.UpdateAppointmentVitals)
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
				prescriptions.POST("/:id/send-to-pharmacy", middleware.RequireRole("admin", "doctor"), prescriptionHandler.SendToPharmacy)
				prescriptions.POST("/:id/send-to-lab", middleware.RequireRole("admin", "doctor"), prescriptionHandler.SendToLab)
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
				lab.POST("/orders", middleware.RequireRole("admin", "doctor", "lab"), labHandler.CreateOrder)
				lab.GET("/orders", middleware.RequireRole("admin", "doctor", "nurse", "lab"), labHandler.ListOrders)
				lab.GET("/orders/:id", middleware.RequireRole("admin", "doctor", "nurse", "patient", "lab"), labHandler.GetOrder)
				lab.PATCH("/orders/:id/status", middleware.RequireRole("admin", "nurse", "lab"), labHandler.UpdateStatus)
				lab.POST("/orders/:id/report", middleware.RequireRole("admin", "nurse", "lab"), labHandler.UploadReport)
				lab.GET("/patients/:patient_id/orders", middleware.RequireRole("admin", "doctor", "nurse", "patient", "lab"), labHandler.GetPatientLabOrders)
			}

			// Billing routes
			billing := protected.Group("/billing")
			{
				billing.GET("", middleware.RequireRole("admin", "doctor", "patient"), billingHandler.ListBills)
				billing.POST("/generate", middleware.RequireRole("admin", "doctor", "nurse", "receptionist"), billingHandler.GenerateBill)
				billing.GET("/:id", middleware.RequireRole("admin", "doctor", "patient"), billingHandler.GetBill)
				billing.PATCH("/:id/pay", middleware.RequireRole("admin", "patient"), billingHandler.PayBill)
			}

			// OPD/Reception routes
			opd := protected.Group("/opd")
			{
				// UHID routes
				uhid := opd.Group("/uhid")
				{
					uhid.POST("/generate", middleware.RequireRole("admin", "receptionist"), uhidHandler.GenerateUHID)
					uhid.GET("/patient/:patient_id", middleware.RequireRole("admin", "receptionist", "doctor", "nurse"), uhidHandler.GetUHIDByPatient)
					uhid.GET("/validate/:uhid", middleware.RequireRole("admin", "receptionist", "doctor", "nurse"), uhidHandler.ValidateUHID)
				}

				// Insurance policy routes
				insurance := opd.Group("/insurance")
				{
					insurance.POST("/policies", middleware.RequireRole("admin", "receptionist"), insuranceHandler.CreateInsurancePolicy)
					insurance.GET("/patient/:patient_id", middleware.RequireRole("admin", "receptionist", "doctor", "patient"), insuranceHandler.GetInsurancePoliciesByPatient)
					insurance.GET("/active/:patient_id", middleware.RequireRole("admin", "receptionist", "doctor", "patient"), insuranceHandler.GetActiveInsurancePolicy)
					insurance.PUT("/policies/:policy_id", middleware.RequireRole("admin", "receptionist"), insuranceHandler.UpdateInsurancePolicy)
					insurance.DELETE("/policies/:policy_id", middleware.RequireRole("admin", "receptionist"), insuranceHandler.DeactivateInsurancePolicy)
				}

				// Waiting list routes
				waitingList := opd.Group("/waiting-list")
				{
					waitingList.POST("/add", middleware.RequireRole("admin", "receptionist"), waitingListHandler.AddToWaitingList)
					waitingList.GET("/doctor/:doctor_id", middleware.RequireRole("admin", "receptionist", "doctor"), waitingListHandler.GetDoctorWaitingList)
					waitingList.POST("/call-next/:doctor_id", middleware.RequireRole("admin", "receptionist"), waitingListHandler.CallNextPatient)
					waitingList.PUT("/:id/mark-seen", middleware.RequireRole("admin", "doctor", "nurse"), waitingListHandler.MarkPatientSeen)
					waitingList.PUT("/:id/complete", middleware.RequireRole("admin", "doctor"), waitingListHandler.CompleteConsultation)
					waitingList.DELETE("/:id/cancel", middleware.RequireRole("admin", "receptionist", "patient"), waitingListHandler.CancelWaitingListEntry)
					waitingList.GET("/patient/:patient_id/position", middleware.RequireRole("admin", "patient"), waitingListHandler.GetPatientQueuePosition)
					waitingList.GET("/doctor/:doctor_id/estimated-time", middleware.RequireRole("admin", "receptionist", "patient"), waitingListHandler.GetEstimatedWaitTime)
				}

				// OPD Receipt routes
				receipt := opd.Group("/receipt")
				{
					receipt.POST("/generate", middleware.RequireRole("admin", "receptionist"), opdReceiptHandler.GenerateOPDReceipt)
					receipt.GET("/:id", middleware.RequireRole("admin", "receptionist", "doctor", "patient"), opdReceiptHandler.GetOPDReceipt)
					receipt.GET("/number/:receipt_number", middleware.RequireRole("admin", "receptionist", "doctor", "patient"), opdReceiptHandler.GetOPDReceiptByNumber)
					receipt.GET("/patient/:patient_id", middleware.RequireRole("admin", "receptionist", "patient"), opdReceiptHandler.GetPatientOPDReceipts)
					receipt.POST("/:id/print", middleware.RequireRole("admin", "receptionist"), opdReceiptHandler.PrintOPDReceipt)
					receipt.GET("/:id/summary", middleware.RequireRole("admin", "receptionist", "doctor", "patient"), opdReceiptHandler.GetReceiptSummary)
				}
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

			// Rating routes
			ratings := protected.Group("/ratings")
			{
				ratings.POST("", middleware.RequireRole("admin", "patient"), ratingHandler.CreateRating)
				ratings.GET("", middleware.RequireRole("admin", "doctor", "patient"), ratingHandler.ListRatings)
				ratings.GET("/doctor/:doctor_id", middleware.RequireRole("admin", "doctor", "patient"), ratingHandler.GetDoctorRatings)
				ratings.GET("/patient/:patient_id", middleware.RequireRole("admin", "patient"), ratingHandler.GetPatientRatings)
				ratings.GET("/:id", middleware.RequireRole("admin", "doctor", "patient"), ratingHandler.GetRating)
				ratings.PUT("/:id", middleware.RequireRole("admin", "patient"), ratingHandler.UpdateRating)
				ratings.DELETE("/:id", middleware.RequireRole("admin", "patient"), ratingHandler.DeleteRating)
			}

			// Bed management routes
			beds := protected.Group("/beds")
			{
				beds.GET("/all", bedHandler.GetAllBedsWithOccupancy)
				beds.GET("/stats", bedHandler.GetBedStats)
				beds.GET("/available", bedHandler.GetAvailableBeds)
				beds.GET("/patient/:patient_id", bedHandler.GetPatientBed)
				beds.GET("/ward/:ward", bedHandler.GetBedsByWard)
				beds.GET("/ward/:ward/stats", bedHandler.GetWardStats)
				beds.GET("/status/:status", bedHandler.GetBedsByStatus)
				beds.POST("/admit", middleware.RequireRole("admin", "nurse"), bedHandler.AdmitPatient)
				beds.POST("/discharge", middleware.RequireRole("admin", "nurse"), bedHandler.DischargePatient)
				beds.PATCH("/:id/status", middleware.RequireRole("admin", "nurse"), bedHandler.UpdateBedStatus)
			}

			// Admission records routes (complete CRUD)
			admissions := protected.Group("/admissions")
			{
				admissions.GET("", admissionHandler.GetAllAdmissions)
				admissions.POST("", middleware.RequireRole("admin", "doctor", "nurse"), admissionHandler.CreateAdmission)
				admissions.GET("/:id", admissionHandler.GetAdmissionByID)
				admissions.PATCH("/:id", middleware.RequireRole("admin", "doctor", "nurse"), admissionHandler.UpdateAdmission)
				admissions.DELETE("/:id", middleware.RequireRole("admin"), admissionHandler.DeleteAdmission)
				admissions.GET("/active", admissionHandler.GetActiveAdmissions)
				admissions.GET("/patient/:patient_id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), admissionHandler.GetPatientAdmissions)
			}

			// Progress notes routes
			progressNotes := protected.Group("/progress-notes")
			{
				progressNotes.GET("", middleware.RequireRole("admin", "doctor", "nurse"), doctorHandler.GetProgressNotes)
				progressNotes.POST("", middleware.RequireRole("admin", "doctor"), doctorHandler.CreateProgressNote)
				progressNotes.GET("/:id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), doctorHandler.GetProgressNote)
				progressNotes.PATCH("/:id", middleware.RequireRole("admin", "doctor"), doctorHandler.UpdateProgressNote)
				progressNotes.DELETE("/:id", middleware.RequireRole("admin", "doctor"), doctorHandler.DeleteProgressNote)
				progressNotes.GET("/admission/:admission_id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), doctorHandler.GetAdmissionProgressNotes)
			}

			// Nurse instructions routes
			nurseInstructions := protected.Group("/nurse-instructions")
			{
				nurseInstructions.GET("", middleware.RequireRole("admin", "doctor", "nurse"), doctorHandler.GetNurseInstructions)
				nurseInstructions.POST("", middleware.RequireRole("admin", "doctor"), doctorHandler.CreateNurseInstruction)
				nurseInstructions.GET("/:id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), doctorHandler.GetNurseInstruction)
				nurseInstructions.PATCH("/:id", middleware.RequireRole("admin", "doctor"), doctorHandler.UpdateNurseInstruction)
				nurseInstructions.DELETE("/:id", middleware.RequireRole("admin", "doctor"), doctorHandler.DeleteNurseInstruction)
				nurseInstructions.GET("/admission/:admission_id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), doctorHandler.GetAdmissionNurseInstructions)
			}

			// Discharge summaries routes
			dischargeSummaries := protected.Group("/discharge-summaries")
			{
				dischargeSummaries.GET("", middleware.RequireRole("admin", "doctor", "nurse"), doctorHandler.GetDischargeSummaries)
				dischargeSummaries.POST("", middleware.RequireRole("admin", "doctor"), doctorHandler.CreateDischargeSummary)
				dischargeSummaries.GET("/:id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), doctorHandler.GetDischargeSummary)
				dischargeSummaries.PATCH("/:id", middleware.RequireRole("admin", "doctor"), doctorHandler.UpdateDischargeSummary)
				dischargeSummaries.DELETE("/:id", middleware.RequireRole("admin", "doctor"), doctorHandler.DeleteDischargeSummary)
				dischargeSummaries.GET("/admission/:admission_id", middleware.RequireRole("admin", "doctor", "nurse", "patient"), doctorHandler.GetAdmissionDischargeSummary)
			}

			// Export routes (data exports to Excel/PDF)
			exports := protected.Group("/exports")
			{
				exports.GET("/appointments/excel", middleware.RequireRole("admin", "doctor"), exportHandler.ExportAppointmentsExcel)
				exports.GET("/prescriptions/excel", middleware.RequireRole("admin", "doctor"), exportHandler.ExportPrescriptionsExcel)
				exports.GET("/admissions/excel", middleware.RequireRole("admin", "doctor"), exportHandler.ExportAdmissionsExcel)
				exports.GET("/lab-orders/excel", middleware.RequireRole("admin", "doctor"), exportHandler.ExportLabOrdersExcel)
				exports.GET("/discharge-summary/:admission_id", middleware.RequireRole("admin", "doctor", "patient"), exportHandler.ExportDischargeSummaryText)
			}

			// Advanced Analytics routes (new comprehensive analytics endpoints)
			advancedAnalytics := protected.Group("/advanced-analytics")
			{
				// Doctor statistics
				advancedAnalytics.GET("/doctors/:doctor_id/statistics", middleware.RequireRole("admin", "doctor"), analyticsHandler.GetDoctorStatistics)
				advancedAnalytics.GET("/doctors/rankings", middleware.RequireRole("admin", "doctor"), analyticsHandler.GetDoctorRankings)

				// Appointment trends
				advancedAnalytics.GET("/appointments/trends", middleware.RequireRole("admin", "doctor"), analyticsHandler.GetAppointmentTrends)

				// Lab metrics
				advancedAnalytics.GET("/lab/completion", middleware.RequireRole("admin", "doctor"), analyticsHandler.GetLabCompletionMetrics)

				// Patient outcomes
				advancedAnalytics.GET("/patients/outcomes", middleware.RequireRole("admin", "doctor"), analyticsHandler.GetPatientOutcomeStatistics)

				// Revenue analytics
				advancedAnalytics.GET("/doctors/:doctor_id/revenue", middleware.RequireRole("admin", "doctor"), analyticsHandler.GetRevenueStatistics)

				// Performance metrics
				advancedAnalytics.GET("/performance", middleware.RequireRole("admin"), analyticsHandler.GetPerformanceMetrics)
				advancedAnalytics.GET("/cache/stats", middleware.RequireRole("admin"), analyticsHandler.GetCacheStats)
			}

			// Analytics routes (legacy)
			analytics := protected.Group("/analytics")
			{
				analytics.GET("/admissions", middleware.RequireRole("admin", "doctor", "nurse"), doctorHandler.GetAdmissionStats)
				analytics.GET("/doctor/:doctor_id/performance", middleware.RequireRole("admin", "doctor"), doctorHandler.GetDoctorPerformanceMetrics)
			}

			// Debug route - check current user's role
			protected.GET("/debug/whoami", func(c *gin.Context) {
				userID, _ := c.Get("userID")
				role, _ := c.Get("role")
				email, _ := c.Get("email")
				c.JSON(200, gin.H{
					"user_id":   userID,
					"email":     email,
					"role":      role,
					"role_type": fmt.Sprintf("%T", role),
				})
			})

			// Debug route - check database counts
			protected.GET("/debug/counts", func(c *gin.Context) {
				var appointmentCount int64
				var patientCount int64
				var doctorCount int64
				var userCount int64

				config.DB.Model(&model.Appointment{}).Count(&appointmentCount)
				config.DB.Model(&model.Patient{}).Count(&patientCount)
				config.DB.Model(&model.Doctor{}).Count(&doctorCount)
				config.DB.Model(&model.User{}).Count(&userCount)

				c.JSON(200, gin.H{
					"appointments": appointmentCount,
					"patients":     patientCount,
					"doctors":      doctorCount,
					"users":        userCount,
				})
			})

			// Debug route - raw appointments data
			protected.GET("/debug/appointments-raw", func(c *gin.Context) {
				var appointments []model.Appointment
				result := config.DB.
					Preload("Patient").
					Preload("Patient.User").
					Preload("Doctor").
					Preload("Doctor.User").
					Limit(50).
					Order("created_at DESC").
					Find(&appointments)

				if result.Error != nil {
					c.JSON(500, gin.H{"error": result.Error.Error()})
					return
				}

				c.JSON(200, gin.H{
					"count":        len(appointments),
					"appointments": appointments,
				})
			})
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

// debugHardiDoctor prints Doctor Hardi details and his appointments
func debugHardiDoctor() {
	var doctor model.Doctor
	result := config.DB.Joins("JOIN users ON users.id = doctors.user_id").
		Where("users.first_name ILIKE ? OR users.last_name ILIKE ?", "%Hardi%", "%Hardi%").
		Preload("User").
		First(&doctor)

	if result.Error != nil {
		log.Printf("[DEBUG] No Hardi doctor found in database: %v", result.Error)
		return
	}

	log.Printf("\n===== HARDI DOCTOR DEBUG INFO =====")
	log.Printf("Doctor Name: %s %s", doctor.User.FirstName, doctor.User.LastName)
	log.Printf("Doctor ID: %s", doctor.ID)
	log.Printf("Specialization: %s", doctor.Specialization)
	log.Printf("====================================\n")

	// Find appointments for this doctor
	var appointments []model.Appointment
	result = config.DB.Where("doctor_id = ?", doctor.ID).
		Preload("Patient").
		Preload("Patient.User").
		Order("appointment_date ASC, appointment_time ASC").
		Find(&appointments)

	if result.Error != nil {
		log.Printf("[DEBUG] Error finding appointments: %v", result.Error)
		return
	}

	log.Printf("\n===== %s's APPOINTMENTS (%d total) =====", doctor.User.FirstName, len(appointments))
	for i, apt := range appointments {
		patientName := "Unknown"
		if apt.Patient != nil && apt.Patient.User != nil {
			patientName = apt.Patient.User.FirstName + " " + apt.Patient.User.LastName
		}
		log.Printf("[%d] %s | %s %s | Status: %s | ID: %s",
			i+1, apt.AppointmentDate, apt.AppointmentTime, patientName, apt.Status, apt.ID)
	}
	log.Printf("==========================================\n")
}
