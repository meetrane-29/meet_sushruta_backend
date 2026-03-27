package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillItemType string

const (
	BillItemConsultation BillItemType = "consultation"
	BillItemMedicine     BillItemType = "medicine"
	BillItemLab          BillItemType = "lab"
	BillItemOther        BillItemType = "other"
)

type BillItem struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BillID      uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bill_id"`
	Bill        *Bill          `gorm:"foreignKey:BillID;references:ID" json:"bill,omitempty"`
	ItemType    BillItemType   `gorm:"type:varchar(50);not null" json:"item_type"` // consultation, medicine, lab, other
	Description string         `gorm:"type:text;not null" json:"description"`
	Quantity    int            `gorm:"default:1" json:"quantity"`
	UnitPrice   float64        `gorm:"not null" json:"unit_price"`
	Total       float64        `gorm:"not null" json:"total"` // quantity * unit_price
	CreatedAt   int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (BillItem) TableName() string {
	return "bill_items"
}

func (bi *BillItem) BeforeCreate(tx *gorm.DB) error {
	if bi.ID == uuid.Nil {
		bi.ID = uuid.New()
	}
	// Calculate total
	bi.Total = float64(bi.Quantity) * bi.UnitPrice
	return nil
}

func (bi *BillItem) BeforeSave(tx *gorm.DB) error {
	// Recalculate total on save
	bi.Total = float64(bi.Quantity) * bi.UnitPrice
	return nil
}
