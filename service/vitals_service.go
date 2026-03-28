package service

import (
	"fmt"
	"strconv"
	"strings"

	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
)

type VitalsService interface {
	CreateVitals(vitals *model.Vitals) error
	GetVitalsByID(id uuid.UUID) (*model.Vitals, error)
	GetLatestVitalsByPatientID(patientID uuid.UUID) (*model.Vitals, error)
	GetVitalsHistoryByPatientID(patientID uuid.UUID, limit int, offset int) ([]model.Vitals, error)
	GetVitalsByPatientID(patientID uuid.UUID) ([]model.Vitals, error)
	UpdateVitals(vitals *model.Vitals) error
	DeleteVitals(id uuid.UUID) error
	GetVitalsTrend(patientID uuid.UUID, days int) (map[string]interface{}, error)
	GetVitalsAlerts(patientID uuid.UUID) ([]map[string]interface{}, error)
	CheckAbnormalReadings(vitals *model.Vitals) []map[string]interface{}
}

type vitalsService struct {
	vitalsRepo repository.VitalsRepository
}

func NewVitalsService(vitalsRepo repository.VitalsRepository) VitalsService {
	return &vitalsService{
		vitalsRepo: vitalsRepo,
	}
}

// CreateVitals creates a new vitals record
func (s *vitalsService) CreateVitals(vitals *model.Vitals) error {
	return s.vitalsRepo.CreateVitals(vitals)
}

// GetVitalsByID retrieves vitals by ID
func (s *vitalsService) GetVitalsByID(id uuid.UUID) (*model.Vitals, error) {
	return s.vitalsRepo.GetVitalsByID(id)
}

// GetLatestVitalsByPatientID retrieves the most recent vitals for a patient
func (s *vitalsService) GetLatestVitalsByPatientID(patientID uuid.UUID) (*model.Vitals, error) {
	return s.vitalsRepo.GetLatestVitalsByPatientID(patientID)
}

// GetVitalsHistoryByPatientID retrieves vitals history with pagination
func (s *vitalsService) GetVitalsHistoryByPatientID(patientID uuid.UUID, limit int, offset int) ([]model.Vitals, error) {
	return s.vitalsRepo.GetVitalsHistoryByPatientID(patientID, limit, offset)
}

// GetVitalsByPatientID retrieves all vitals for a patient
func (s *vitalsService) GetVitalsByPatientID(patientID uuid.UUID) ([]model.Vitals, error) {
	// Get a large number of vitals (e.g., last 1000 records)
	// In practice, you might want to limit this or add date range filtering
	return s.vitalsRepo.GetVitalsHistoryByPatientID(patientID, 1000, 0)
}

// UpdateVitals updates an existing vitals record
func (s *vitalsService) UpdateVitals(vitals *model.Vitals) error {
	return s.vitalsRepo.UpdateVitals(vitals)
}

// DeleteVitals deletes a vitals record
func (s *vitalsService) DeleteVitals(id uuid.UUID) error {
	return s.vitalsRepo.DeleteVitals(id)
}

// GetVitalsTrend returns vitals trend data for the last N days
func (s *vitalsService) GetVitalsTrend(patientID uuid.UUID, days int) (map[string]interface{}, error) {
	vitals, err := s.vitalsRepo.GetVitalsHistoryByPatientID(patientID, 1000, 0)
	if err != nil {
		return nil, err
	}

	if len(vitals) == 0 {
		return map[string]interface{}{
			"message": "No vitals data available",
			"trend":   []interface{}{},
		}, nil
	}

	// Get latest vitals
	latest := vitals[len(vitals)-1]

	// Get previous vitals (before latest)
	var previous *model.Vitals
	if len(vitals) > 1 {
		previous = &vitals[len(vitals)-2]
	}

	trend := map[string]interface{}{
		"patient_id": patientID.String(),
		"latest": map[string]interface{}{
			"id":               latest.ID,
			"temperature":      latest.Temperature,
			"blood_pressure":   latest.BloodPressure,
			"heart_rate":       latest.HeartRate,
			"respiratory_rate": latest.RespiratoryRate,
			"weight":           latest.Weight,
			"height":           latest.Height,
			"spo2":             latest.Oxygen,
			"blood_sugar":      latest.BloodSugar,
			"recorded_at":      latest.RecordedAt,
		},
		"changes": map[string]interface{}{
			"temperature": calculateChange(latest.Temperature, previous.Temperature),
			"heart_rate":  calculateIntChange(latest.HeartRate, previous.HeartRate),
			"weight":      calculateChange(latest.Weight, previous.Weight),
			"blood_sugar": calculateChange(latest.BloodSugar, previous.BloodSugar),
			"spo2":        calculateIntChange(latest.Oxygen, previous.Oxygen),
		},
		"total_records": len(vitals),
	}

	return trend, nil
}

// GetVitalsAlerts returns abnormal readings for the patient
func (s *vitalsService) GetVitalsAlerts(patientID uuid.UUID) ([]map[string]interface{}, error) {
	vitals, err := s.vitalsRepo.GetLatestVitalsByPatientID(patientID)
	if err != nil {
		return nil, err
	}

	if vitals == nil {
		return []map[string]interface{}{}, nil
	}

	alerts := s.CheckAbnormalReadings(vitals)
	return alerts, nil
}

// CheckAbnormalReadings checks vitals against normal ranges and returns alerts
func (s *vitalsService) CheckAbnormalReadings(vitals *model.Vitals) []map[string]interface{} {
	var alerts []map[string]interface{}

	// Temperature check: Normal 36.1-37.2°C (assumed input is in Celsius)
	if vitals.Temperature < 36.1 || vitals.Temperature > 37.2 {
		severity := "medium"
		if vitals.Temperature < 35.0 || vitals.Temperature > 38.5 {
			severity = "critical"
		}
		alerts = append(alerts, map[string]interface{}{
			"parameter":    "temperature",
			"value":        vitals.Temperature,
			"normal_range": "36.1-37.2°C",
			"severity":     severity,
			"message":      fmt.Sprintf("Temperature %.1f°C is outside normal range", vitals.Temperature),
		})
	}

	// Heart Rate check: Normal 60-100 bpm
	if vitals.HeartRate < 60 || vitals.HeartRate > 100 {
		severity := "medium"
		if vitals.HeartRate < 40 || vitals.HeartRate > 120 {
			severity = "critical"
		}
		alerts = append(alerts, map[string]interface{}{
			"parameter":    "heart_rate",
			"value":        vitals.HeartRate,
			"normal_range": "60-100 bpm",
			"severity":     severity,
			"message":      fmt.Sprintf("Heart rate %d bpm is outside normal range", vitals.HeartRate),
		})
	}

	// Respiratory Rate check: Normal 12-20 breaths/min
	if vitals.RespiratoryRate < 12 || vitals.RespiratoryRate > 20 {
		severity := "medium"
		if vitals.RespiratoryRate < 8 || vitals.RespiratoryRate > 30 {
			severity = "critical"
		}
		alerts = append(alerts, map[string]interface{}{
			"parameter":    "respiratory_rate",
			"value":        vitals.RespiratoryRate,
			"normal_range": "12-20 breaths/min",
			"severity":     severity,
			"message":      fmt.Sprintf("Respiratory rate %d breaths/min is outside normal range", vitals.RespiratoryRate),
		})
	}

	// SpO2 check: Normal >95%
	if vitals.Oxygen < 95 {
		severity := "medium"
		if vitals.Oxygen < 90 {
			severity = "critical"
		}
		alerts = append(alerts, map[string]interface{}{
			"parameter":    "spo2",
			"value":        vitals.Oxygen,
			"normal_range": ">95%",
			"severity":     severity,
			"message":      fmt.Sprintf("SpO2 %d%% is low", vitals.Oxygen),
		})
	}

	// Blood Sugar check: Fasting normal 70-100 mg/dL, non-fasting <140 mg/dL
	if vitals.BloodSugar > 0 { // Only check if recorded
		if vitals.BloodSugar < 70 || vitals.BloodSugar > 140 {
			severity := "medium"
			if vitals.BloodSugar < 50 || vitals.BloodSugar > 200 {
				severity = "critical"
			}
			alerts = append(alerts, map[string]interface{}{
				"parameter":    "blood_sugar",
				"value":        vitals.BloodSugar,
				"normal_range": "70-140 mg/dL",
				"severity":     severity,
				"message":      fmt.Sprintf("Blood sugar %.1f mg/dL is outside normal range", vitals.BloodSugar),
			})
		}
	}

	// Blood Pressure check
	if vitals.BloodPressure != "" {
		systolic, diastolic := parseBloodPressure(vitals.BloodPressure)
		if systolic > 0 && diastolic > 0 {
			// Normal: Systolic <120 and Diastolic <80
			// Elevated: Systolic 120-129 and Diastolic <80
			// High: Systolic >=130 or Diastolic >=80
			if systolic >= 130 || diastolic >= 80 {
				severity := "medium"
				if systolic >= 160 || diastolic >= 100 {
					severity = "critical"
				}
				alerts = append(alerts, map[string]interface{}{
					"parameter":    "blood_pressure",
					"value":        vitals.BloodPressure,
					"normal_range": "<120/80",
					"severity":     severity,
					"message":      fmt.Sprintf("Blood pressure %s is high", vitals.BloodPressure),
				})
			}
		}
	}

	return alerts
}

// Helper functions

func parseBloodPressure(bp string) (int, int) {
	parts := strings.Split(bp, "/")
	if len(parts) != 2 {
		return 0, 0
	}
	systolic, _ := strconv.Atoi(parts[0])
	diastolic, _ := strconv.Atoi(parts[1])
	return systolic, diastolic
}

func calculateChange(latest, previous float64) map[string]interface{} {
	if previous == 0 {
		return map[string]interface{}{
			"value":       latest,
			"change":      0.0,
			"change_type": "N/A",
		}
	}
	change := latest - previous
	changeType := "stable"
	if change > 0.5 {
		changeType = "increased"
	} else if change < -0.5 {
		changeType = "decreased"
	}
	return map[string]interface{}{
		"value":       latest,
		"change":      fmt.Sprintf("%.2f", change),
		"change_type": changeType,
	}
}

func calculateIntChange(latest, previous int) map[string]interface{} {
	if previous == 0 {
		return map[string]interface{}{
			"value":       latest,
			"change":      0,
			"change_type": "N/A",
		}
	}
	change := latest - previous
	changeType := "stable"
	if change > 2 {
		changeType = "increased"
	} else if change < -2 {
		changeType = "decreased"
	}
	return map[string]interface{}{
		"value":       latest,
		"change":      change,
		"change_type": changeType,
	}
}
