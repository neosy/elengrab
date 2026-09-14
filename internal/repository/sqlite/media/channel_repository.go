package media

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	emedia "github.com/neosy/elengrab/internal/repository/sqlite/media/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/media/mappers"
)

type ChannelRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	// options
	retryOptions dbexec.RetryOptions
}

// NewChannelRepository returns a new object for the repository
func NewChannelRepository(dbEntry persistence.DBEntry) persistence.ChannelRepositoryFactory {
	return func() persistence.ChannelRepository {
		return &ChannelRepository{
			mappers: mappers.NewMappers(),
			dbEntry: dbEntry,

			// options
			retryOptions: dbexec.RetryOptions{
				MaxRetries: maxRetriesDefault,
				Delay:      retryDelayDefault,
			},
		}
	}
}

func (r *ChannelRepository) Insert(ctx context.Context, channel *dmedia.Channel) error {
	return r.Save(ctx, channel)
}

func (r *ChannelRepository) Update(ctx context.Context, channel *dmedia.Channel) error {
	return r.Save(ctx, channel)
}

func (r *ChannelRepository) Save(ctx context.Context, channel *dmedia.Channel) error {
	if channel == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eChannel, err := r.mappers.MapChannelDomainToEntity(channel)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eChannel.InsertFields()
	values := eChannel.InsertValues()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eChannel.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eChannel.FieldName(&eChannel.ChannelID))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save Channel: %v", err)
	}

	return nil
}

func (r *ChannelRepository) UpdateChannelID(ctx context.Context, oldChannelID, newChannelID uuid.UUID) error {
	if oldChannelID == newChannelID {
		return nil
	}

	var eChannel emedia.Channel

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Update(eChannel.TableName()).
		Set(eChannel.FieldName(&eChannel.ChannelID), newChannelID.String()).
		Where(squirrel.Eq{eChannel.FieldName(&eChannel.ChannelID): oldChannelID.String()}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save Channel: %v", err)
	}

	return nil
}

func (r *ChannelRepository) FindByChannelID(ctx context.Context, channelID uuid.UUID) (*dmedia.Channel, error) {
	var ent emedia.Channel

	sqlQuery, args, err := squirrel.Select(ent.QueryFields()...).
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.ChannelID): channelID.String()}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	var notFound bool
	db := dbexec.Resolve(ctx, r.dbEntry)
	execQuery := func() error {
		row := db.QueryRowContext(ctx, sqlQuery, args...)
		// Scan result into entity
		err := row.Scan(ent.FieldPointers()...)
		if err == sql.ErrNoRows {
			notFound = true
			return nil
		}
		return err
	}
	err = dbexec.ExecRetry(ctx, r.retryOptions, execQuery)
	if err != nil {
		return nil, err
	}
	if notFound {
		return nil, nil
	}

	// Map entity to domain model
	channel, err := r.mappers.MapChannelEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *ChannelRepository) ExistsByChannelID(ctx context.Context, channelID uuid.UUID) (bool, error) {
	var ent emedia.Channel

	// Build SQL query: SELECT 1 FROM table WHERE channel_id = $1 LIMIT 1
	query, args, err := squirrel.Select("1").
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.ChannelID): channelID.String()}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	db := dbexec.Resolve(ctx, r.dbEntry)

	// Execute query and check if any row exists
	var exists int
	err = db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *ChannelRepository) FindByExternalChannelID(
	ctx context.Context,
	externalID string,
	platform string,
) (*dmedia.Channel, error) {
	var ent emedia.Channel

	sqlWhere := squirrel.Eq{
		ent.FieldName(&ent.ExternalID): externalID,
		ent.FieldName(&ent.Platform):   platform,
	}

	sqlQuery, args, err := squirrel.Select(ent.QueryFields()...).
		From(ent.TableName()).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	var notFound bool
	db := dbexec.Resolve(ctx, r.dbEntry)
	execQuery := func() error {
		row := db.QueryRowContext(ctx, sqlQuery, args...)
		// Scan result into entity
		err := row.Scan(ent.FieldPointers()...)
		if err == sql.ErrNoRows {
			notFound = true
			return nil
		}
		return err
	}
	err = dbexec.ExecRetry(ctx, r.retryOptions, execQuery)
	if err != nil {
		return nil, err
	}
	if notFound {
		return nil, nil
	}

	// Map entity to domain model
	channel, err := r.mappers.MapChannelEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *ChannelRepository) ExistsByExternalChannelID(
	ctx context.Context,
	channelID string,
	platform string,
) (bool, error) {
	var ent emedia.Channel

	sqlWhere := squirrel.Eq{
		ent.FieldName(&ent.ChannelID): channelID,
		ent.FieldName(&ent.Platform):  platform,
	}

	// Build SQL query: SELECT 1 FROM table WHERE channel_id = $1 LIMIT 1
	query, args, err := squirrel.Select("1").
		From(ent.TableName()).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	db := dbexec.Resolve(ctx, r.dbEntry)

	// Execute query and check if any row exists
	var exists int
	err = db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *ChannelRepository) IterateAll(ctx context.Context, fn func(*dmedia.Channel) error) error {
	var eChannel emedia.Channel

	orderBys := dbutils.SortBy(eChannel.FieldName(&eChannel.CreatedAt), dbutils.OrderAscending).List()

	qb := squirrel.Select(eChannel.QueryFields()...).
		From(eChannel.TableName()).
		OrderBy(orderBys.Query()).
		PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := qb.ToSql()

	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	db := dbexec.Resolve(ctx, r.dbEntry)
	rows, err := db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows != nil {
		for rows.Next() {
			err := rows.Scan(eChannel.FieldPointers()...)
			if err != nil {
				return err
			}

			channel, err := r.mappers.MapChannelEntityToDomain(&eChannel)
			if err != nil {
				return err
			}

			err = fn(channel)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *ChannelRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}

func (r *ChannelRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.dbEntry, fn)
}
