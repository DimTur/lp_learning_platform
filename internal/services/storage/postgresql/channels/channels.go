package channels

import (
	"context"
	"errors"
	"fmt"

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

func (c *ChannelPostgresStorage) CreateChannel(ctx context.Context, channel CreateChannel) (int64, error) {
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

const getChannelWithPlansQuery = `
	SELECT
		c.id AS channel_id,
		c.name AS channel_name,
		c.description AS channel_description,
		c.created_by AS channel_created_by,
		c.last_modified_by AS channel_last_modified_by,
		c.created_at AS channel_created_at,
		c.modified AS channel_modified,
		p.id AS plan_id,
		p.name AS plan_name,
		p.description AS plan_description,
		p.created_by AS plan_created_by,
		p.last_modified_by AS plan_last_modified_by,
		p.is_published AS plan_is_published,
		p.public AS plan_public,
		p.created_at AS plan_created_at,
		p.modified AS plan_modified
	FROM
		channels c
	LEFT JOIN
		channels_plans cp ON c.id = cp.channel_id
	LEFT JOIN
		plans p ON cp.plan_id = p.id
	WHERE
		c.id = $1`

func (c *ChannelPostgresStorage) GetChannelByID(ctx context.Context, channelID int64) (ChannelWithPlans, error) {
	const op = "storage.postgresql.channels.channels.GetChannelByID"

	rows, err := c.db.Query(ctx, getChannelWithPlansQuery, channelID)
	if err != nil {
		return ChannelWithPlans{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var dbChannel DBChannelWithPlans
	dbChannel.Plans = []DBPlanInChannels{}

	for rows.Next() {
		var plan DBPlanInChannels
		err := rows.Scan(
			&dbChannel.ID,
			&dbChannel.Name,
			&dbChannel.Description,
			&dbChannel.CreatedBy,
			&dbChannel.LastModifiedBy,
			&dbChannel.CreatedAt,
			&dbChannel.Modified,
			&plan.ID,
			&plan.Name,
			&plan.Description,
			&plan.CreatedBy,
			&plan.LastModifiedBy,
			&plan.IsPublished,
			&plan.Public,
			&plan.CreatedAt,
			&plan.Modified,
		)
		if err != nil {
			return ChannelWithPlans{}, fmt.Errorf("%s: %w", op, err)
		}

		if plan.ID.Valid {
			dbChannel.Plans = append(dbChannel.Plans, plan)
		}
	}

	if err := rows.Err(); err != nil {
		return ChannelWithPlans{}, fmt.Errorf("%s: %w", op, err)
	}

	channel := ChannelWithPlans{
		ID:             dbChannel.ID,
		Name:           dbChannel.Name,
		Description:    dbChannel.Description,
		CreatedBy:      dbChannel.CreatedBy,
		LastModifiedBy: dbChannel.LastModifiedBy,
		CreatedAt:      dbChannel.CreatedAt,
		Modified:       dbChannel.Modified,
		Plans:          make([]PlanInChannel, 0, len(dbChannel.Plans)),
	}

	for _, dbPlan := range dbChannel.Plans {
		plan := PlanInChannel{
			ID:             dbPlan.ID.Int64,
			Name:           dbPlan.Name.String,
			Description:    dbPlan.Description.String,
			CreatedBy:      dbPlan.CreatedBy.String,
			LastModifiedBy: dbPlan.LastModifiedBy.String,
			IsPublished:    dbPlan.IsPublished.Bool,
			Public:         dbPlan.Public.Bool,
			CreatedAt:      dbPlan.CreatedAt.Time,
			Modified:       dbPlan.Modified.Time,
		}
		channel.Plans = append(channel.Plans, plan)
	}

	return channel, nil
}

const getChannelsQuery = `
	SELECT *
	FROM channels
	ORDER BY id
	LIMIT $1 OFFSET $2`

func (c *ChannelPostgresStorage) GetChannels(ctx context.Context, limit, offset int64) ([]Channel, error) {
	const op = "storage.postgresql.channels.channels.GetChannels"

	var channels []DBChannel

	rows, err := c.db.Query(ctx, getChannelsQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var channel DBChannel
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

	var mappedChannels []Channel
	for _, channel := range channels {
		mappedChannels = append(mappedChannels, Channel(channel))
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

func (c *ChannelPostgresStorage) UpdateChannel(ctx context.Context, updChannel UpdateChannelRequest) (int64, error) {
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

const chareChannelQuery = `
	INSERT INTO shared_channels_learninggroups(channel_id, learning_group_id, created_by, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

func (c *ChannelPostgresStorage) ShareChannelToGroup(ctx context.Context, s DBShareChannelToGroup) error {
	const op = "storage.postgresql.channels.channels.ShareChannelToGroup"

	var id int64

	err := c.db.QueryRow(ctx, chareChannelQuery,
		s.ChannelID,
		s.LGroupID,
		s.CreatedBy,
		s.CreatedAt,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("Postgres error code: %s, message: %s\n", pgErr.Code, pgErr.Message)
			if pgErr.Code == "23505" { // unique violation code
				return fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
			} else if pgErr.Code == "23503" { // foreign key violation code
				return fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
