package alarmold

import (
	"context"
	"database/sql"
	"fmt"

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

// FindByRouterID Handles the creation of a new alarm record on database
func (r *DatabaseRepository) FindByRouterID(ctx context.Context, routerID int) (bool, string, error) {
	model := &AlarmModel{}
	query := r.db.NewSelect().
		Model(model).
		Where("id_router = ?", routerID).
		Where("waiting = ?", 1).
		WhereOr("attending = ?", 1).
		Order("created DESC").
		Limit(1)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return false, "", terrors.New(terrors.ErrNotFound, "Router information not found", map[string]string{})
		}
		return false, "", fmt.Errorf("router_old_repository: Error while searching for routersvc -> %w", err)
	}

	return true, model.ID, nil
}
