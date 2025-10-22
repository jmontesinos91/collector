package unitsold

import (
	"context"
)

// IRepository interface
type IRepository interface {
	FindByRouterID(ctx context.Context, routerID int) (*UnitsModel, error)
	FindByID(ctx context.Context, unitID int) (*UnitsModel, error)
}
