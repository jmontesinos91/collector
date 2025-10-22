package unitsold

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

// FindByRouterID Handles the find by imei of router record on old database
func (r *DatabaseRepository) FindByRouterID(ctx context.Context, routerID int) (*UnitsModel, error) {
	model := &UnitsModel{}
	query := r.db.NewSelect().
		Model(model).
		Where("id_router = ?", routerID).
		Limit(1)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return model, terrors.New(terrors.ErrNotFound, "Unit information not found", map[string]string{})
		}
		return model, fmt.Errorf("unit_old_repository: Error while searching for routersvc -> %w", err)
	}

	return model, nil
}

// FindByID Handles the find by unitID record on old database
func (r *DatabaseRepository) FindByID(ctx context.Context, unitID int) (*UnitsModel, error) {
	model := &UnitsModel{}
	query := r.db.NewSelect().
		Model(model).
		Where("id = ?", unitID).
		Limit(1)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return model, terrors.New(terrors.ErrNotFound, "Unit information not found", map[string]string{})
		}
		return model, fmt.Errorf("unit_old_repository: Error while searching for routersvc -> %w", err)
	}

	return model, nil
}
