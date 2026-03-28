package repository

import (
	"meet_sushruta/model"

	"github.com/google/uuid"
)

// PatientRepository defines methods for patient data access
type PatientRepository interface {
	Create(patient *model.Patient) error
	GetByID(id uuid.UUID) (*model.Patient, error)
	GetAll(page, limit int, search string) ([]model.Patient, int64, error)
	Update(patient *model.Patient) error
	SoftDelete(id uuid.UUID) error
	GetByUserID(userID uuid.UUID) (*model.Patient, error)
}

// DoctorRepository defines methods for doctor data access
type DoctorRepository interface {
	Create(doctor *model.Doctor) error
	GetByID(id uuid.UUID) (*model.Doctor, error)
	GetAll(page, limit int, search string) ([]model.Doctor, int64, error)
	Update(doctor *model.Doctor) error
	SoftDelete(id uuid.UUID) error
	GetByUserID(userID uuid.UUID) (*model.Doctor, error)
}

// AppointmentRepository defines methods for appointment data access
type AppointmentRepository interface {
	Create(appointment *model.Appointment) error
	GetByID(id uuid.UUID) (*model.Appointment, error)
	GetAll(page, limit int) ([]model.Appointment, int64, error)
	Update(appointment *model.Appointment) error
	SoftDelete(id uuid.UUID) error
	GetByPatientID(patientID uuid.UUID, page, limit int) ([]model.Appointment, int64, error)
	GetByDoctorID(doctorID uuid.UUID, page, limit int) ([]model.Appointment, int64, error)
	GetByDateRange(doctorID uuid.UUID, startDate, endDate string) ([]model.Appointment, error)
}

// DoctorScheduleRepository defines methods for doctor schedule data access
type DoctorScheduleRepository interface {
	GetByDoctorID(doctorID uuid.UUID) ([]model.DoctorSchedule, error)
	GetByDoctorAndDay(doctorID uuid.UUID, dayOfWeek string) (*model.DoctorSchedule, error)
}

// BillingRepository defines methods for bill data access
type BillingRepository interface {
	Create(bill *model.Bill) error
	GetByID(id uuid.UUID) (*model.Bill, error)
	GetAll(page, limit int) ([]model.Bill, int64, error)
	GetByPatientID(patientID uuid.UUID, page, limit int) ([]model.Bill, int64, error)
	GetByStatus(status string, page, limit int) ([]model.Bill, int64, error)
	GetByDateRange(startDate, endDate string, page, limit int) ([]model.Bill, int64, error)
	Update(bill *model.Bill) error
	SoftDelete(id uuid.UUID) error
	CreateBillItem(billItem *model.BillItem) error
	GetBillItems(billID uuid.UUID) ([]model.BillItem, error)
}

// MedicineRepository defines methods for medicine data access
type MedicineRepository interface {
	Create(medicine *model.Medicine) error
	GetByID(id uuid.UUID) (*model.Medicine, error)
	GetAll(page, limit int, search string) ([]model.Medicine, int64, error)
	Update(medicine *model.Medicine) error
	SoftDelete(id uuid.UUID) error
	GetLowStock(reorderLevel int, page, limit int) ([]model.Medicine, int64, error)
	GetExpiring(daysThreshold int, page, limit int) ([]model.Medicine, int64, error)
	UpdateStock(medicineID uuid.UUID, quantity int) error
	CreateDispenseHistory(history *model.DispenseHistory) error
	GetDispenseHistoryByPrescriptionID(prescriptionID uuid.UUID) ([]model.DispenseHistory, error)
}

// BedRepository defines methods for bed data access
type BedRepository interface {
	Create(bed *model.Bed) error
	GetByID(id uuid.UUID) (*model.Bed, error)
	GetAll(page, limit int) ([]model.Bed, int64, error)
	Update(bed *model.Bed) error
	SoftDelete(id uuid.UUID) error
	GetByWard(ward string) ([]model.Bed, error)
	GetByStatus(status string) ([]model.Bed, error)
	GetByBedType(bedType string) ([]model.Bed, error)
	GetWardStats() (map[string]interface{}, error)
	GetBedStats() (map[string]interface{}, error)
}

// MedicalEquipmentRepository defines methods for medical equipment data access
type MedicalEquipmentRepository interface {
	Create(equipment *model.MedicalEquipment) error
	GetByID(id uuid.UUID) (*model.MedicalEquipment, error)
	GetAll(page, limit int) ([]model.MedicalEquipment, int64, error)
	Update(equipment *model.MedicalEquipment) error
	SoftDelete(id uuid.UUID) error
	GetByStatus(status string) ([]model.MedicalEquipment, error)
	GetCriticalEquipment() ([]model.MedicalEquipment, error)
	GetEquipmentStats() (map[string]interface{}, error)
}

// OperationTheatreRepository defines methods for operation theatre data access
type OperationTheatreRepository interface {
	Create(theatre *model.OperationTheatre) error
	GetByID(id uuid.UUID) (*model.OperationTheatre, error)
	GetAll(page, limit int) ([]model.OperationTheatre, int64, error)
	Update(theatre *model.OperationTheatre) error
	SoftDelete(id uuid.UUID) error
	GetByStatus(status string) ([]model.OperationTheatre, error)
}

// OperationScheduleRepository defines methods for operation schedule data access
type OperationScheduleRepository interface {
	Create(operation *model.OperationSchedule) error
	GetByID(id uuid.UUID) (*model.OperationSchedule, error)
	GetAll(page, limit int) ([]model.OperationSchedule, int64, error)
	Update(operation *model.OperationSchedule) error
	SoftDelete(id uuid.UUID) error
	GetByStatus(status string) ([]model.OperationSchedule, error)
	GetByDate(date string) ([]model.OperationSchedule, error)
	GetByTheatreID(theatreID uuid.UUID) ([]model.OperationSchedule, error)
	GetUpcomingOperations(daysAhead int) ([]model.OperationSchedule, error)
	GetOperationStats() (map[string]interface{}, error)
}
