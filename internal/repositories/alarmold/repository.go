package alarmold

import (
	"context"
)

// IRepository interface
type IRepository interface {
	FindByRouterID(ctx context.Context, routerID int) (bool, string, error)
}
