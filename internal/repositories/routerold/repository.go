package routerold

import (
	"context"
)

// IRepository interface
type IRepository interface {
	FindByIMEI(ctx context.Context, imei string) (*RouterModel, error)
	UpdateLatAndLong(ctx context.Context, routerID int, lat, long string) error
}
