package config

import (
	seed "meet_sushruta/seed"
)

// SeedDatabase calls the comprehensive seeding function from the seed package
func SeedDatabase() error {
	return seed.ComprehensiveSeed(GetDB())
}
