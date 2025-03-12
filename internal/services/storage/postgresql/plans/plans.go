package plans

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlansPostgresStorage struct {
	db *pgxpool.Pool
}

func NewPlansStorage(db *pgxpool.Pool) *PlansPostgresStorage {
	return &PlansPostgresStorage{db: db}
}

const (
	createPlanQuery = `
	INSERT INTO plans(name, description, created_by, last_modified_by, created_at, modified)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`
	createChannelsPlansQuery = `
	INSERT INTO channels_plans(channel_id, plan_id)
	VALUES ($1, $2)`
)

func (p *PlansPostgresStorage) CreatePlan(ctx context.Context, plan *CreatePlan) (int64, error) {
	const op = "storage.postgresql.plans.plans.CreatePlan"

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Printf("%s: %v", op, storage.ErrRollBack)
			}
		}
	}()

	var planID int64
	err = tx.QueryRow(ctx, createPlanQuery,
		plan.Name,
		plan.Description,
		plan.CreatedBy,
		plan.LastModifiedBy,
		plan.CreatedAt,
		plan.Modified,
	).Scan(&planID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique violation code
				return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx,
		createChannelsPlansQuery,
		plan.ChannelID,
		planID,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
	}

	return planID, nil
}

// TODO: if not group admin you didn't see is_published = false
const getPlanByIDQuery = `
	SELECT
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
		plans p
	INNER JOIN 
		channels_plans cp ON p.id = cp.plan_id
	INNER JOIN 
		channels c ON cp.channel_id = c.id
	WHERE 
		plan_id = $1
		AND cp.channel_id = $2`

func (p *PlansPostgresStorage) GetPlanByID(ctx context.Context, planCh *GetPlan) (Plan, error) {
	const op = "storage.postgresql.plans.plans.GetPlanByID"

	var plan DBPlan

	err := p.db.QueryRow(ctx, getPlanByIDQuery, planCh.PlanID, planCh.ChannelID).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return (Plan)(plan), fmt.Errorf("%s: %w", op, storage.ErrPlanNotFound)
		}
		return (Plan)(plan), fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}

	return (Plan)(plan), nil
}

const getPlansAll = `
	SELECT
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
		plans p
	INNER JOIN 
		channels_plans cp ON p.id = cp.plan_id
	INNER JOIN 
		channels c ON cp.channel_id = c.id
	INNER JOIN
		shared_plans_users spu ON p.id = spu.plan_id
	WHERE 
		cp.channel_id = $1
		AND spu.user_id = $2
	ORDER BY p.id
	LIMIT $3 OFFSET $4;`

func (p *PlansPostgresStorage) GetPlansAll(ctx context.Context, inputParams *GetPlans) ([]Plan, error) {
	const op = "storage.postgresql.plans.plans.GetPlans"

	var plans []DBPlan

	rows, err := p.db.Query(
		ctx,
		getPlansAll,
		inputParams.ChannelID,
		inputParams.UserID,
		inputParams.Limit,
		inputParams.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var plan DBPlan
		if err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Description,
			&plan.CreatedBy,
			&plan.LastModifiedBy,
			&plan.IsPublished,
			&plan.Public,
			&plan.CreatedAt,
			&plan.Modified,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrScanFailed)
		}
		plans = append(plans, plan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var mappedPlans []Plan
	for _, plan := range plans {
		mappedPlans = append(mappedPlans, Plan(plan))
	}

	return mappedPlans, nil
}

const getPlansQuery = `
	SELECT
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
		plans p
	INNER JOIN 
		channels_plans cp ON p.id = cp.plan_id
	INNER JOIN 
		channels c ON cp.channel_id = c.id
	INNER JOIN
		shared_plans_users spu ON p.id = spu.plan_id
	WHERE 
		cp.channel_id = $1
		AND spu.user_id = $2
		AND public 
		AND is_published
	ORDER BY p.id
	LIMIT $3 OFFSET $4;`

func (p *PlansPostgresStorage) GetPlans(ctx context.Context, inputParams *GetPlans) ([]Plan, error) {
	const op = "storage.postgresql.plans.plans.GetPlans"

	var plans []DBPlan

	rows, err := p.db.Query(
		ctx,
		getPlansQuery,
		inputParams.ChannelID,
		inputParams.UserID,
		inputParams.Limit,
		inputParams.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var plan DBPlan
		if err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Description,
			&plan.CreatedBy,
			&plan.LastModifiedBy,
			&plan.IsPublished,
			&plan.Public,
			&plan.CreatedAt,
			&plan.Modified,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrScanFailed)
		}
		plans = append(plans, plan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var mappedPlans []Plan
	for _, plan := range plans {
		mappedPlans = append(mappedPlans, Plan(plan))
	}

	return mappedPlans, nil
}

const updatePlanQuery = `
	UPDATE plans p
	SET name = COALESCE($3, p.name), 
    	description = COALESCE($4, p.description), 
    	last_modified_by = $5, 
    	is_published = COALESCE($6, p.is_published), 
    	public = COALESCE($7, p.public), 
    	modified = now() 
	FROM 
		channels_plans cp
	WHERE 
		p.id = $1
		AND cp.channel_id = $2
	RETURNING 
		p.id;`

func (p *PlansPostgresStorage) UpdatePlan(ctx context.Context, updPlan *UpdatePlanRequest) (int64, error) {
	const op = "storage.postgresql.plans.plans.UpdatePlan"

	var id int64

	err := p.db.QueryRow(ctx, updatePlanQuery,
		updPlan.PlanID,
		updPlan.ChannelID,
		updPlan.Name,
		updPlan.Description,
		updPlan.LastModifiedBy,
		updPlan.IsPublished,
		updPlan.Public,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrPlanNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}

	return id, nil
}

const deletePlanQuery = `
	DELETE FROM plans p
	USING channels_plans cp, channels c
	WHERE p.id = cp.plan_id
	  AND cp.channel_id = c.id
	  AND p.id = $1
	  AND cp.channel_id = $2;`

func (p *PlansPostgresStorage) DeletePlan(ctx context.Context, planCh *DeletePlan) error {
	const op = "storage.postgresql.plans.plans.DeletePlan"

	res, err := p.db.Exec(
		ctx,
		deletePlanQuery,
		planCh.PlanID,
		planCh.ChannelID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrPlanNotFound)
	}

	return nil
}

const canShareQuery = `
	SELECT EXISTS (
	SELECT 1 
	FROM
		channels_plans cp
	WHERE
		cp.channel_id = $1
		AND cp.plan_id = $2
	);`

func (c *PlansPostgresStorage) CanShare(ctx context.Context, cs *DBCanShare) (bool, error) {
	const op = "storage.postgresql.plans.plans.CanShare"

	var exists bool
	err := c.db.QueryRow(
		ctx,
		canShareQuery,
		cs.ChannelID,
		cs.PlanID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

const sharedPlanQuery = `
	INSERT INTO shared_plans_users(plan_id, user_id, created_by, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

func (c *PlansPostgresStorage) SharePlanWithUser(ctx context.Context, s *DBSharePlanForUser) error {
	const op = "storage.postgresql.plans.plans.SharePlanWithUser"

	var id int64

	err := c.db.QueryRow(ctx, sharedPlanQuery,
		s.PlanID,
		s.UserID,
		s.CreatedBy,
		s.CreatedAt,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, storage.ErrPlanAlreadySharedWithUser)
			}
			if pgErr.Code == "23503" {
				return fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Forming the basic part of the request
const baseQuery = `
	INSERT INTO shared_plans_users (plan_id, user_id, created_by, created_at)
	VALUES %s
	ON CONFLICT (plan_id, user_id) DO NOTHING` // Ignore duplicates

func (c *PlansPostgresStorage) BatchSharePlansWithUsers(ctx context.Context, bs *SharePlanForUsers) error {
	const op = "storage.postgresql.plans.BatchSharePlansWithUsers"

	if len(bs.UserIDs) == 0 {
		return nil
	}

	// Prepare placeholders and arguments
	valueStrings := make([]string, 0, len(bs.UserIDs))
	valueArgs := make([]interface{}, 0, len(bs.UserIDs)*4)

	for i, userID := range bs.UserIDs {
		// $1, $2, $3, $4 -> dynamically increase placeholders
		placeholders := fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		valueStrings = append(valueStrings, placeholders)
		valueArgs = append(valueArgs, bs.PlanID, userID, bs.CreatedBy, bs.CreatedAt)
	}

	// Combining placeholders into one query
	query := fmt.Sprintf(baseQuery, strings.Join(valueStrings, ", "))

	// Do query
	_, err := c.db.Exec(ctx, query, valueArgs...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// TODO: create index for plan_id and user_id
const isUserShareWithPlanQuery = `
	SELECT EXISTS (
    SELECT 1
    FROM shared_plans_users spu
    WHERE spu.plan_id = $1 AND spu.user_id = $2
);`

func (c *PlansPostgresStorage) IsUserShareWithPlan(ctx context.Context, userPlan *IsUserShareWithPlan) (bool, error) {
	const op = "storage.postgresql.plans.plans.IsUserShareWithPlan"

	var exists bool
	err := c.db.QueryRow(
		ctx,
		isUserShareWithPlanQuery,
		userPlan.PlanID,
		userPlan.UserID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

const getPlansForSharingQuery = `
	SELECT
		p.id AS plan_id,
		sclg.channel_id AS channel_id
	FROM
		plans p
	INNER JOIN
		channels_plans cp ON p.id = cp.plan_id
	INNER JOIN
		shared_channels_learninggroups sclg ON cp.channel_id = sclg.channel_id
	WHERE
		sclg.learning_group_id = $1;`

func (c *PlansPostgresStorage) GetPlansForSharing(ctx context.Context, lgPlan *LearningGroup) (map[int64][]int64, error) {
	const op = "storage.postgresql.plans.plans.GetPlansForSharing"

	rows, err := c.db.Query(
		ctx,
		getPlansForSharingQuery,
		lgPlan.LgID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	channelPlansMap := make(map[int64][]int64)

	for rows.Next() {
		var channelID, planID int64
		if err := rows.Scan(&planID, &channelID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrScanFailed)
		}
		channelPlansMap[channelID] = append(channelPlansMap[channelID], planID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return channelPlansMap, nil
}
