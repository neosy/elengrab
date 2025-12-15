package sldownload

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/Masterminds/squirrel"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/download/mappers"
	"github.com/neosy/elengrab/pkg/dbutils"
)

type YoutubeChannelRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	mu      *sync.RWMutex

	// options
	retryOptions retryOptions
}

// NewYoutubeChannelRepository returns a new object for the repository
func NewYoutubeChannelRepository(db *sql.DB, mu *sync.RWMutex) *YoutubeChannelRepository {
	return &YoutubeChannelRepository{
		mappers: mappers.NewMappers(),
		db:      db,
		mu:      mu,

		// options
		retryOptions: retryOptions{
			maxRetries: maxRetriesDefault,
			delay:      retryDelayDefault,
		},
	}
}

func (r *YoutubeChannelRepository) Insert(ctx context.Context, channel *ddownload.YoutubeChannel) error {
	return r.save(ctx, channel, false)
}

func (r *YoutubeChannelRepository) Update(ctx context.Context, channel *ddownload.YoutubeChannel) error {
	return r.save(ctx, channel, true)
}

func (r *YoutubeChannelRepository) save(ctx context.Context, channel *ddownload.YoutubeChannel, isUpd bool) error {
	if channel == nil {
		return errors.New("function parameter is a null pointer")
	}

	// Convert the domain model to a database entity
	eChannel, err := r.mappers.MapYoutubeChannelDomainToEntity(channel)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eChannel.Fields()
	values := eChannel.Values()

	// If this is an update — add the UpdatedAt field with the current time
	if isUpd {
		fields = append(fields, eChannel.FieldName(&eChannel.UpdatedAt))
		values = append(values, squirrel.Expr("CURRENT_TIMESTAMP"))
	}

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
	err = execContext(ctx, r.db, r.mu, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save youtubeChannel: %v", err)
	}

	return nil
}

func (r *YoutubeChannelRepository) FindByChannelId(ctx context.Context, channelID string) (*ddownload.YoutubeChannel, error) {
	var ent edownload.YoutubeChannel

	query, args, err := squirrel.Select(ent.FieldsAll()...).
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.ChannelID): channelID}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	db := dbOrTx(ctx, r.db)
	row := db.QueryRowContext(ctx, query, args...)

	// Scan result into entity
	if err := row.Scan(ent.FieldPointers()...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	// Map entity to domain model
	channel, err := r.mappers.MapYoutubeChannelEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *YoutubeChannelRepository) ExistsByChannelId(ctx context.Context, channelID string) (bool, error) {
	var ent edownload.YoutubeChannel

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
	db := dbOrTx(ctx, r.db)

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
