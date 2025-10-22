package facilitylocationsold

import (
	"github.com/uptrace/bun"
)

// FacilityLocationsModel Database model for router
type FacilityLocationsModel struct {
	bun.BaseModel `bun:"table:facility_locations"`

	UnitID     int `bun:"id_unit"`
	LocationID int `bun:"id_location"`
}
