package locationsold

import (
	"time"

	"github.com/uptrace/bun"
)

// LocationsModel Database model for router
type LocationsModel struct {
	bun.BaseModel `bun:"table:locations"`

	ID        int        `bun:"id,pk"`
	AlarmID   int        `bun:"id_alarm"`
	Latitude  string     `bun:"lat"`
	Longitude string     `bun:"lng"`
	UpdatedAt *time.Time `bun:"date_update"`
}
