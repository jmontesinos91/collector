package facilitylocationsold

import (
	"context"

	"github.com/jmontesinos91/ologs/logger"
	"github.com/uptrace/bun"
)

// DatabaseRepository struct
type DatabaseRepository struct {
	log *logger.ContextLogger
	db  *bun.DB
}

// NewDatabaseRepository creates an instance of DatabaseRepository
func NewDatabaseRepository(l *logger.ContextLogger, conn *bun.DB) *DatabaseRepository {
	return &DatabaseRepository{
		log: l,
		db:  conn,
	}
}

// Create Handles the creation of a new Facility Location record on database.
func (r *DatabaseRepository) Create(ctx context.Context, model *FacilityLocationsModel) error {
	_, err := r.db.NewInsert().
		Model(model).
		Exec(ctx)

	// Handling error
	if err != nil {
		return err
	}

	return nil
}
