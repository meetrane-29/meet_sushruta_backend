package repository

import (
	"time"

	"meet_sushruta/config"
	"meet_sushruta/model"
)

type AdminRepository interface {
	GetTotalPatients() (int64, error)
	GetTotalDoctors() (int64, error)
	GetTotalNurses() (int64, error)
	GetTotalUsers() (int64, error)
	GetTodayAppointments() (int64, error)
	GetPendingBills() (int64, error)
	GetRecentPatients() ([]map[string]interface{}, error)
	GetTopDoctorsByAppointmentsToday() ([]map[string]interface{}, error)
	GetDepartmentLoad() ([]map[string]interface{}, error)
	GetTodayAppointmentsByDoctor() ([]map[string]interface{}, error)
	GetNurseRoleBreakdown() ([]map[string]interface{}, error)
}

type adminRepository struct{}

func NewAdminRepository() AdminRepository {
	return &adminRepository{}
}

// GetTotalPatients returns count of total patients
func (r *adminRepository) GetTotalPatients() (int64, error) {
	var count int64
	err := config.DB.Model(&model.Patient{}).Count(&count).Error
	return count, err
}

// GetTotalDoctors returns count of total doctors
func (r *adminRepository) GetTotalDoctors() (int64, error) {
	var count int64
	err := config.DB.Model(&model.Doctor{}).Count(&count).Error
	return count, err
}

// GetTotalNurses returns count of total nurses
func (r *adminRepository) GetTotalNurses() (int64, error) {
	var count int64
	err := config.DB.Model(&model.Nurse{}).Count(&count).Error
	return count, err
}

// GetTodayAppointments uses appointment_date string field (YYYY-MM-DD) — no timestamp casting needed
func (r *adminRepository) GetTodayAppointments() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := config.DB.Model(&model.Appointment{}).
		Where("appointment_date = ?", today).
		Count(&count).Error
	return count, err
}

// GetPendingBills returns count of bills with payment_status='pending'
func (r *adminRepository) GetPendingBills() (int64, error) {
	var count int64
	err := config.DB.Model(&model.Bill{}).
		Where("payment_status = ?", "pending").
		Count(&count).Error
	return count, err
}

// GetRecentPatients returns last 5 patients with their user info
// Fixes:
//   - EXTRACT(YEAR FROM AGE(...)) returns float8 in PostgreSQL; cast to int with ::int
//   - Guard empty/invalid date_of_birth with regex check before calling TO_DATE
func (r *adminRepository) GetRecentPatients() ([]map[string]interface{}, error) {
	query := `
		SELECT
			p.id::text,
			CONCAT(u.first_name, ' ', u.last_name) AS name,
			CASE
				WHEN p.date_of_birth ~ '^\d{4}-\d{2}-\d{2}$'
				THEN FLOOR(EXTRACT(YEAR FROM AGE(TO_DATE(p.date_of_birth, 'YYYY-MM-DD'))))::int
				ELSE 0
			END AS age,
			COALESCE(NULLIF(d.department, ''), 'General') AS department,
			'active' AS status,
			p.created_at
		FROM patients p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN doctors d ON u.id = d.user_id
		ORDER BY p.created_at DESC
		LIMIT 5
	`

	rows, err := config.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, name, department, status string
		var age int
		var createdAt int64

		if err := rows.Scan(&id, &name, &age, &department, &status, &createdAt); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id":         id,
			"name":       name,
			"age":        age,
			"department": department,
			"status":     status,
			"created_at": createdAt,
		})
	}

	return result, nil
}

// GetTopDoctorsByAppointmentsToday uses appointment_date field — no timestamp casting needed
// Rating is normalized 0-1 based on today's appointment count (max 5 = full rating)
func (r *adminRepository) GetTopDoctorsByAppointmentsToday() ([]map[string]interface{}, error) {
	today := time.Now().Format("2006-01-02")

	query := `
		SELECT
			d.id::text,
			CONCAT(u.first_name, ' ', u.last_name) AS name,
			d.specialization,
			COUNT(a.id) AS today_appointments,
			LEAST(1.0, COUNT(a.id)::float / 5.0) AS rating
		FROM doctors d
		JOIN users u ON d.user_id = u.id
		LEFT JOIN appointments a ON d.id = a.doctor_id AND a.appointment_date = ?
		GROUP BY d.id, u.first_name, u.last_name, d.specialization
		ORDER BY today_appointments DESC
		LIMIT 3
	`

	rows, err := config.DB.Raw(query, today).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, name, specialization string
		var todayAppointments int
		var rating float64

		if err := rows.Scan(&id, &name, &specialization, &todayAppointments, &rating); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id":                 id,
			"name":               name,
			"specialization":     specialization,
			"today_appointments": todayAppointments,
			"rating":             rating,
		})
	}

	return result, nil
}

// GetDepartmentLoad returns department patient counts for today using appointment_date field
func (r *adminRepository) GetDepartmentLoad() ([]map[string]interface{}, error) {
	today := time.Now().Format("2006-01-02")

	query := `
		SELECT
			COALESCE(NULLIF(d.department, ''), 'General') AS department,
			COUNT(DISTINCT a.patient_id) AS count,
			ROUND(
				COUNT(DISTINCT a.patient_id) * 100.0
				/ NULLIF(SUM(COUNT(DISTINCT a.patient_id)) OVER (), 0),
				2
			) AS percentage
		FROM appointments a
		LEFT JOIN doctors d ON a.doctor_id = d.id
		WHERE a.appointment_date = ?
		GROUP BY d.department
		ORDER BY count DESC
	`

	rows, err := config.DB.Raw(query, today).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var department string
		var count int
		var percentage float64

		if err := rows.Scan(&department, &count, &percentage); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"department": department,
			"count":      count,
			"percentage": percentage,
		})
	}

	return result, nil
}

// GetTotalUsers returns count of total users
func (r *adminRepository) GetTotalUsers() (int64, error) {
	var count int64
	err := config.DB.Model(&model.User{}).Count(&count).Error
	return count, err
}

// GetTodayAppointmentsByDoctor returns all appointments for today grouped by doctor with patient and doctor details
func (r *adminRepository) GetTodayAppointmentsByDoctor() ([]map[string]interface{}, error) {
	today := time.Now().Format("2006-01-02")

	query := `
		SELECT
			a.id::text,
			a.appointment_date,
			a.appointment_time,
			a.status,
			a.reason,
			a.notes,
			a.created_at,
			d.id::text as doctor_id,
			CONCAT(du.first_name, ' ', du.last_name) AS doctor_name,
			d.specialization,
			d.department,
			p.id::text as patient_id,
			CONCAT(pu.first_name, ' ', pu.last_name) AS patient_name,
			pu.phone as patient_phone
		FROM appointments a
		JOIN doctors d ON a.doctor_id = d.id
		JOIN users du ON d.user_id = du.id
		JOIN patients p ON a.patient_id = p.id
		JOIN users pu ON p.user_id = pu.id
		WHERE a.appointment_date = ?
		ORDER BY d.id, a.appointment_time
	`

	rows, err := config.DB.Raw(query, today).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var appointmentID, appointmentDate, appointmentTime, status, reason, notes string
		var createdAt int64
		var doctorID, doctorName, specialization, department string
		var patientID, patientName, patientPhone string

		if err := rows.Scan(&appointmentID, &appointmentDate, &appointmentTime, &status, &reason, &notes, &createdAt,
			&doctorID, &doctorName, &specialization, &department, &patientID, &patientName, &patientPhone); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"id":                    appointmentID,
			"appointment_date":      appointmentDate,
			"appointment_time":      appointmentTime,
			"status":                status,
			"reason":                reason,
			"notes":                 notes,
			"created_at":            createdAt,
			"doctor_id":             doctorID,
			"doctor_name":           doctorName,
			"doctor_specialization": specialization,
			"doctor_department":     department,
			"patient_id":            patientID,
			"patient_name":          patientName,
			"patient_phone":         patientPhone,
		})
	}

	return result, nil
}

// GetNurseRoleBreakdown returns count of nurses by each role
func (r *adminRepository) GetNurseRoleBreakdown() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			n.role,
			COUNT(n.id) as count,
			ROUND(COUNT(n.id) * 100.0 / NULLIF(SUM(COUNT(n.id)) OVER (), 0), 2) as percentage
		FROM nurses n
		WHERE n.deleted_at IS NULL
		GROUP BY n.role
		ORDER BY count DESC
	`

	rows, err := config.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var role *string // Use pointer to handle NULL values
		var count int
		var percentage float64

		if err := rows.Scan(&role, &count, &percentage); err != nil {
			return nil, err
		}

		roleStr := "Unassigned"
		if role != nil && *role != "" {
			roleStr = *role
		}

		result = append(result, map[string]interface{}{
			"role":       roleStr,
			"count":      count,
			"percentage": percentage,
		})
	}

	return result, nil
}
