package plans

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
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

func (p *PlansPostgresStorage) CreatePlan(ctx context.Context, plan CreatePlan) (int64, error) {
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

const getPlanByIDQuery = `
	SELECT id, name, description, created_by, last_modified_by, is_published, public, created_at, modified 
	FROM plans 
	WHERE id = $1`

func (p *PlansPostgresStorage) GetPlanByID(ctx context.Context, planID int64) (Plan, error) {
	const op = "storage.postgresql.plans.plans.GetPlanByID"

	var plan DBPlan

	err := p.db.QueryRow(ctx, getPlanByIDQuery, planID).Scan(
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
		return (Plan)(plan), fmt.Errorf("%s: %w", op, storage.ErrPlanNotFound)
	}

	return (Plan)(plan), nil
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
	WHERE cp.channel_id = $1
	ORDER BY p.id
	LIMIT $2 OFFSET $3;`

func (p *PlansPostgresStorage) GetPlans(ctx context.Context, channel_id int64, limit, offset int64) ([]Plan, error) {
	const op = "storage.postgresql.plans.plans.GetPlans"

	var plans []DBPlan

	rows, err := p.db.Query(ctx, getPlansQuery, channel_id, limit, offset)
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
	UPDATE plans 
	SET name = COALESCE($2, name), 
	    description = COALESCE($3, description), 
	    last_modified_by = $4, 
	    is_published = COALESCE($5, is_published), 
	    public = COALESCE($6, public), 
	    modified = now() 
	WHERE id = $1
	RETURNING id`

func (p *PlansPostgresStorage) UpdatePlan(ctx context.Context, updPlan UpdatePlanRequest) (int64, error) {
	const op = "storage.postgresql.plans.plans.UpdatePlan"

	var id int64

	err := p.db.QueryRow(ctx, updatePlanQuery,
		updPlan.ID,
		updPlan.Name,
		updPlan.Description,
		updPlan.LastModifiedBy,
		updPlan.IsPublished,
		updPlan.Public,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}
	return id, nil
}

const deletePlanQuery = `
	DELETE FROM plans
	WHERE id = $1`

func (p *PlansPostgresStorage) DeletePlan(ctx context.Context, id int64) error {
	const op = "storage.postgresql.plans.plans.DeletePlan"

	res, err := p.db.Exec(ctx, deletePlanQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, storage.ErrPlanNotFound)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrPlanNotFound)
	}

	return nil
}
