package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MedicalEquipment struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	EquipmentName       string         `gorm:"uniqueIndex;not null" json:"equipment_name"`
	EquipmentType       string         `gorm:"index;not null" json:"equipment_type"` // e.g., "Ventilator", "Defibrillator", "CT Scanner", "X-Ray Machine", etc.
	Model               string         `json:"model"`
	SerialNumber        string         `json:"serial_number"`
	Manufacturer        string         `json:"manufacturer"`
	PurchaseDate        string         `json:"purchase_date"`                                                                                             // YYYY-MM-DD
	WarrantyExpiryDate  string         `json:"warranty_expiry_date"`                                                                                      // YYYY-MM-DD
	Location            string         `gorm:"index" json:"location"`                                                                                     // Ward, OT, Room number, etc.
	Status              string         `gorm:"default:'working';index:idx_status,index:idx_deleted_status,type:btree;composite:priority:2" json:"status"` // working, under_maintenance, repair, retired
	MaintenanceLastDate string         `json:"maintenance_last_date"`                                                                                     // YYYY-MM-DD
	MaintenanceNextDate string         `json:"maintenance_next_date"`                                                                                     // YYYY-MM-DD
	MaintenanceNotes    string         `gorm:"type:text" json:"maintenance_notes"`
	RepairNotes         string         `gorm:"type:text" json:"repair_notes"`
	TechniciansAssigned string         `json:"technicians_assigned"` // Comma-separated names/IDs
	EstimatedRepairCost float64        `json:"estimated_repair_cost"`
	ActualRepairCost    float64        `json:"actual_repair_cost"`
	DailyRentalRate     float64        `json:"daily_rental_rate"`
	OperationalHours    int            `json:"operational_hours"`                       // Total hours used
	MaxOperationalHours int            `json:"max_operational_hours"`                   // Equipment lifespan
	CriticalEquipment   bool           `gorm:"default:false" json:"critical_equipment"` // true for ICU/OT machines
	Notes               string         `gorm:"type:text" json:"notes"`
	CreatedAt           int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt           int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index,index:idx_deleted_status,type:btree;composite:priority:1" json:"deleted_at"`
}

func (MedicalEquipment) TableName() string {
	return "medical_equipments"
}

func (m *MedicalEquipment) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
