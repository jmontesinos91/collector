package routerold

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/terrors"

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

// FindByIMEI Handles the find by imei of router record on old database
func (r *DatabaseRepository) FindByIMEI(ctx context.Context, imei string) (*RouterModel, error) {
	model := &RouterModel{}
	query := r.db.NewSelect().
		Model(model).
		Where("imei = ?", imei).
		Limit(1)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return model, terrors.New(terrors.ErrNotFound, "Router information not found", map[string]string{})
		}
		return model, fmt.Errorf("router_old_repository: Error while searching for routersvc -> %w", err)
	}

	return model, nil
}

func (r *DatabaseRepository) UpdateLatAndLong(ctx context.Context, routerID int, lat, long string) error {
	_, errUpdate := r.db.NewUpdate().
		Table("routers").
		Set("lat = ?", lat).
		Set("lng = ?", long).
		Set("updated = ?", time.Now().UTC()).
		Where("id = ?", routerID).
		Exec(ctx)
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}
