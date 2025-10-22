package locationsold

import (
	"context"
)

// IRepository interface
type IRepository interface {
	Create(ctx context.Context, model *LocationsModel) (int, error)
}
