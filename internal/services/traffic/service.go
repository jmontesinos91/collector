package traffic

import (
	"context"
	"github.com/jmontesinos91/collector/domains/pagination"
)

type IService interface {
	HandleRetrieve(ctx context.Context, filter *FilterRequest) (pagination.PaginatedRes, error)
	HandleDelete(ctx context.Context, trafficID string) error
	HandleResetCounter(ctx context.Context, trafficID string) error
}
