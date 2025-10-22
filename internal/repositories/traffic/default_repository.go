package traffic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/terrors"
	"github.com/sirupsen/logrus"
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

// Create Handles the creation of a new payout record on database
func (r *DatabaseRepository) Create(ctx context.Context, model *Model) error {
	_, err := r.db.NewInsert().
		Model(model).
		Exec(ctx)

	// Handling error
	if err != nil {
		return err
	}
	return nil
}

// FindByIMEI Handles to find traffic model by imei
func (r *DatabaseRepository) FindByIMEI(ctx context.Context, imei string, isAlarm bool) (bool, error) {
	var tModel []Model
	query := r.db.NewSelect().
		Model(&tModel).
		Where("imei = ?", imei).
		Where("\"isAlarm\" = ?", isAlarm).
		Order("id DESC")

	err := query.Scan(ctx)
	if err != nil {
		return !errors.Is(err, sql.ErrNoRows), err
	} else if len(tModel) > 0 {
		return true, err
	} else {
		return false, err
	}
}

// FindByLastUsed Handles to find traffic model by imei
func (r *DatabaseRepository) FindByLastUsed(ctx context.Context) ([]Model, error) {
	var tModel []Model
	query := r.db.NewSelect().
		Model(&tModel).
		Where("updated_at < current_timestamp - interval '30 minutes'").
		Where("\"isAlarm\" = ?", false).
		Where("\"isnotified\" = ?", false).
		Order("id DESC")

	err := query.Scan(ctx)

	return tModel, err
}

// UpdateIsNotified Handles update as notified
func (r *DatabaseRepository) UpdateIsNotified(ctx context.Context, trafficID string) error {
	_, errUpdate := r.db.NewUpdate().
		Table("traffic").
		Set("isnotified = ?", true).
		Where("id = ?", trafficID).
		Exec(ctx)
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}

// UpdateByIMEI Handles update the register by IMEI
func (r *DatabaseRepository) UpdateByIMEI(ctx context.Context, imei, request string, isAlarm bool) error {
	_, errUpdate := r.db.NewUpdate().
		Table("traffic").
		Set("request = ?", request).
		Set("updated_at = ?", time.Now().UTC()).
		Set("counter=counter+1").
		Set("isnotified = ?", false).
		Where("imei = ?", imei).
		Where("\"isAlarm\" = ?", isAlarm).
		Exec(ctx)
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}

// Retrieve Retrieves traffic data by filters
func (r *DatabaseRepository) Retrieve(ctx context.Context, filter *Metadata) ([]Model, error) {
	var traffics []Model

	query := r.db.NewSelect().Model(&Model{})

	query = setFilters(query, filter)

	if err := query.Scan(ctx, &traffics); err != nil {
		r.log.Error(logrus.ErrorLevel, "Retrieve", "Error scanning traffics", err)
		return nil, terrors.InternalService("count_error", "Error retrieving traffics from the database", map[string]string{})
	}

	return traffics, nil
}

// DeleteByID Handles update the register by IMEI
func (r *DatabaseRepository) DeleteByID(ctx context.Context, trafficID string) error {
	_, errUpdate := r.db.NewDelete().
		Where("id = ?", trafficID).
		Exec(ctx)
	if errUpdate != nil {
		return terrors.InternalService("delete_traffic", "Failed delete traffic from the database", map[string]string{})
	}
	return nil
}

// ResetCounter Handles update the register by IMEI
func (r *DatabaseRepository) ResetCounter(ctx context.Context, trafficID string) error {
	_, errUpdate := r.db.NewUpdate().
		Table("traffic").
		Set("counter = 0").
		Set("updated_at = ?", time.Now().UTC()).
		Where("id = ?", trafficID).
		Exec(ctx)
	if errUpdate != nil {
		return terrors.InternalService("reset_counter", "Failed reset counter traffic from the database", map[string]string{})
	}
	return nil
}

func (r *DatabaseRepository) RetrieveData(ctx context.Context, filter *Metadata) ([]Model, error) {
	var traffics []Model

	query := r.db.NewSelect().Model(&Model{})

	query = setFilters(query, filter)

	if err := query.Scan(ctx, &traffics); err != nil {
		r.log.Error(logrus.ErrorLevel, "Retrieve", "Error scanning traffics", err)
		return nil, terrors.InternalService("count_error", "Error retrieving traffics from the database", map[string]string{})
	}

	return traffics, nil
}

func setFilters(q *bun.SelectQuery, filter *Metadata) *bun.SelectQuery {

	q = q.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {

		if filter.Qparam != "" {
			q = q.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = q.Where("id::text LIKE ?", "%"+filter.Qparam+"%").
					WhereOr("request LIKE ?", "%"+filter.Qparam+"%").
					WhereOr("imei LIKE ?", "%"+filter.Qparam+"%").
					WhereOr("ip LIKE ?", "%"+filter.Qparam+"%")
				return q
			})
		}

		if filter.ID != "" {
			q = q.Where("id::text LIKE ?", "%"+filter.ID+"%")
		}
		if filter.Request != "" {
			q = q.Where("request LIKE ?", "%"+filter.Request+"%")
		}
		if filter.IMEI != "" {
			q = q.Where("imei LIKE ?", "%"+filter.IMEI+"%")
		}
		if filter.Ip != "" {
			q = q.Where("ip LIKE ?", "%"+filter.Ip+"%")
		}
		if filter.IsAlarm != nil {
			q = q.Where("\"isAlarm\" = ?", filter.IsAlarm)
		}
		if filter.Counter != nil {
			if *filter.Counter == 0 {
				q = q.Where("counter = 0")
			} else {
				q = q.Where("counter > ?", filter.Counter)
			}
		}

		return q
	})

	return q
}
