package channels

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimTur/lp_learning_platform/internal/domain/models"
	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChannelPostgresStorage struct {
	db *pgxpool.Pool
}

func NewChannelStorage(db *pgxpool.Pool) *ChannelPostgresStorage {
	return &ChannelPostgresStorage{db: db}
}

const createChannelQuery = `
	INSERT INTO channels(name, description, created_by, last_modified_by, created_at, modified)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`

func (c *ChannelPostgresStorage) CreateChannel(ctx context.Context, channel models.CreateChannel) (int64, error) {
	const op = "storage.postgresql.channels.channels.CreateChannel"

	var id int64

	err := c.db.QueryRow(ctx, createChannelQuery,
		channel.Name,
		channel.Description,
		channel.CreatedBy,
		channel.LastModifiedBy,
		channel.CreatedAt,
		channel.Modified,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique violation code
				return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

const getChannelByIDQuery = `
	SELECT id, name, description, created_by, last_modified_by, created_at, modified 
	FROM channels 
	WHERE id = $1`

func (c *ChannelPostgresStorage) GetChannelByID(ctx context.Context, channelID int64) (models.Channel, error) {
	const op = "storage.postgresql.channels.channels.GetChannelByID"

	var channel models.DBChannel

	err := c.db.QueryRow(ctx, getChannelByIDQuery, channelID).Scan(
		&channel.ID,
		&channel.Name,
		&channel.Description,
		&channel.CreatedBy,
		&channel.LastModifiedBy,
		&channel.CreatedAt,
		&channel.Modified,
	)
	if err != nil {
		return (models.Channel)(channel), fmt.Errorf("%s: %w", op, storage.ErrChannelNotFound)
	}

	return (models.Channel)(channel), nil
}

const getChannelsQuery = `
	SELECT *
	FROM channels
	ORDER BY id
	LIMIT $1 OFFSET $2`

func (c *ChannelPostgresStorage) GetChannels(ctx context.Context, limit, offset int64) ([]models.Channel, error) {
	const op = "storage.postgresql.channels.channels.GetChannels"

	var channels []models.DBChannel

	rows, err := c.db.Query(ctx, getChannelsQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var channel models.DBChannel
		if err := rows.Scan(
			&channel.ID,
			&channel.Name,
			&channel.Description,
			&channel.CreatedBy,
			&channel.LastModifiedBy,
			&channel.CreatedAt,
			&channel.Modified,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrScanFailed)
		}
		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var mappedChannels []models.Channel
	for _, channel := range channels {
		mappedChannels = append(mappedChannels, models.Channel(channel))
	}

	return mappedChannels, nil
}

const updateChannelQuery = `
	UPDATE channels 
	SET name = COALESCE($2, name), 
	    description = COALESCE($3, description), 
	    last_modified_by = $4, 
	    modified = now() 
	WHERE id = $1
	RETURNING id`

func (c *ChannelPostgresStorage) UpdateChannel(ctx context.Context, updChannel models.UpdateChannelRequest) (int64, error) {
	const op = "storage.postgresql.channels.channels.UpdateChannel"

	var id int64

	err := c.db.QueryRow(ctx, updateChannelQuery,
		updChannel.ID,
		updChannel.Name,
		updChannel.Description,
		updChannel.LastModifiedBy,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}
	return id, nil
}

const deleteChannelQuery = `
	DELETE FROM channels
	WHERE id = $1`

func (c *ChannelPostgresStorage) DeleteChannel(ctx context.Context, id int64) error {
	const op = "storage.postgresql.channels.channels.DeleteChannel"

	res, err := c.db.Exec(ctx, deleteChannelQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrChannelNotFound)
	}

	return nil
}
