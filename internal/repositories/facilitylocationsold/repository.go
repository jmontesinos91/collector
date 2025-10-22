package facilitylocationsold

import (
	"context"
)

// IRepository interface
type IRepository interface {
	Create(ctx context.Context, model *FacilityLocationsModel) error
}
