package media

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	emedia "github.com/neosy/elengrab/internal/repository/sqlite/media/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/media/mappers"
)

type YoutubeChannelRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	// options
	retryOptions dbexec.RetryOptions
}

// NewYoutubeChannelRepository returns a new object for the repository
func NewYoutubeChannelRepository(dbEntry persistence.DBEntry) persistence.YoutubeChannelRepositoryFactory {
	return func() persistence.YoutubeChannelRepository {
		return &YoutubeChannelRepository{
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

func (r *YoutubeChannelRepository) Insert(ctx context.Context, channel *dmedia.YoutubeChannel) error {
	return r.Save(ctx, channel)
}

func (r *YoutubeChannelRepository) Update(ctx context.Context, channel *dmedia.YoutubeChannel) error {
	return r.Save(ctx, channel)
}

func (r *YoutubeChannelRepository) Save(ctx context.Context, channel *dmedia.YoutubeChannel) error {
	if channel == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eChannel, err := r.mappers.MapYoutubeChannelDomainToEntity(channel)
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
		return fmt.Errorf("failed to save youtubeChannel: %v", err)
	}

	return nil
}

func (r *YoutubeChannelRepository) FindByChannelID(ctx context.Context, channelID string) (*dmedia.YoutubeChannel, error) {
	var ent emedia.YoutubeChannel

	sqlQuery, args, err := squirrel.Select(ent.QueryFields()...).
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.ChannelID): channelID}).
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
	channel, err := r.mappers.MapYoutubeChannelEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *YoutubeChannelRepository) ExistsByChannelID(ctx context.Context, channelID string) (bool, error) {
	var ent emedia.YoutubeChannel

	// Build SQL query: SELECT 1 FROM table WHERE channel_id = $1 LIMIT 1
	query, args, err := squirrel.Select("1").
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.ChannelID): channelID}).
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

func (r *YoutubeChannelRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}

func (r *YoutubeChannelRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.dbEntry, fn)
}
