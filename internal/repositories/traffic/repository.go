package traffic

import (
	"context"
)

// IRepository interface
type IRepository interface {
	Create(ctx context.Context, model *Model) error
	FindByIMEI(ctx context.Context, imei string, isAlarm bool) (bool, error)
	FindByLastUsed(ctx context.Context) ([]Model, error)
	UpdateIsNotified(ctx context.Context, trafficID string) error
	UpdateByIMEI(ctx context.Context, imei, request string, isAlarm bool) error
	Retrieve(ctx context.Context, filter *Metadata) ([]Model, int, int, error)
	DeleteByID(ctx context.Context, trafficID string) error
	RetrieveData(ctx context.Context, filter *Metadata) ([]Model, error)
}
