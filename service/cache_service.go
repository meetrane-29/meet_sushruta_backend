package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"meet_sushruta/model"

	"github.com/google/uuid"
)

// ============================================================
// CACHING SERVICE - Performance Optimization
// ============================================================

type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

type CacheService struct {
	cache map[string]*CacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewCacheService(defaultTTL time.Duration) *CacheService {
	cs := &CacheService{
		cache: make(map[string]*CacheEntry),
		ttl:   defaultTTL,
	}

	// Start cleanup goroutine
	go cs.cleanupExpired()

	return cs
}

// Set stores a value in cache with expiration
func (c *CacheService) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ttl == 0 {
		ttl = c.ttl
	}

	c.cache[key] = &CacheEntry{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves a value from cache
func (c *CacheService) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Data, true
}

// GetDoctorProfile returns cached doctor profile
func (c *CacheService) GetDoctorProfile(doctorID uuid.UUID) (*model.Doctor, bool) {
	key := fmt.Sprintf("doctor_profile:%s", doctorID.String())
	data, exists := c.Get(key)
	if !exists {
		return nil, false
	}

	doctor, ok := data.(*model.Doctor)
	return doctor, ok
}

// SetDoctorProfile caches a doctor profile
func (c *CacheService) SetDoctorProfile(doctorID uuid.UUID, doctor *model.Doctor, ttl time.Duration) {
	key := fmt.Sprintf("doctor_profile:%s", doctorID.String())
	c.Set(key, doctor, ttl)
}

// GetPatientProfile returns cached patient profile
func (c *CacheService) GetPatientProfile(patientID uuid.UUID) (*model.Patient, bool) {
	key := fmt.Sprintf("patient_profile:%s", patientID.String())
	data, exists := c.Get(key)
	if !exists {
		return nil, false
	}

	patient, ok := data.(*model.Patient)
	return patient, ok
}

// SetPatientProfile caches a patient profile
func (c *CacheService) SetPatientProfile(patientID uuid.UUID, patient *model.Patient, ttl time.Duration) {
	key := fmt.Sprintf("patient_profile:%s", patientID.String())
	c.Set(key, patient, ttl)
}

// InvalidateDoctor removes doctor cache entries
func (c *CacheService) InvalidateDoctor(doctorID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("doctor_profile:%s", doctorID.String())
	delete(c.cache, key)
}

// InvalidatePatient removes patient cache entries
func (c *CacheService) InvalidatePatient(patientID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("patient_profile:%s", patientID.String())
	delete(c.cache, key)
}

// InvalidateAppointments removes appointment cache entries
func (c *CacheService) InvalidateAppointments(doctorID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("appointments:doctor:%s", doctorID.String())
	delete(c.cache, key)
}

// ClearAll removes all cache entries
func (c *CacheService) ClearAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]*CacheEntry)
}

// cleanupExpired removes expired entries periodically
func (c *CacheService) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.cache {
			if now.After(entry.ExpiresAt) {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}

// GetStats returns cache statistics
func (c *CacheService) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	totalEntries := len(c.cache)
	validEntries := 0
	now := time.Now()

	for _, entry := range c.cache {
		if now.Before(entry.ExpiresAt) {
			validEntries++
		}
	}

	return map[string]interface{}{
		"total_entries": totalEntries,
		"valid_entries": validEntries,
		"expired":       totalEntries - validEntries,
	}
}

// ToJSON converts cache entry to JSON bytes
func ToJSON(data interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}

// FromJSON converts JSON bytes to typed data
func FromJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
