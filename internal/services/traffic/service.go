package traffic

import (
	"context"
)

type IService interface {
	HandleRetrieve(ctx context.Context, filter *FilterRequest) ([]Traffic, error)
	HandleDelete(ctx context.Context, trafficID string) error
}
