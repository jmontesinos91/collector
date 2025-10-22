package collector

import (
	"context"
)

// IService Manage routers interfaces
type IService interface {
	Collector(ctx context.Context, payload *Payload) error
}
