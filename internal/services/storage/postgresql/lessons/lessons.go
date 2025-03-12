package lessons

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LessonsPostgresStorage struct {
	db *pgxpool.Pool
}

func NewLessonsStorage(db *pgxpool.Pool) *LessonsPostgresStorage {
	return &LessonsPostgresStorage{db: db}
}

const (
	createLessonQuery = `
	INSERT INTO lessons(name, description, created_by, last_modified_by, created_at, modified)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`
	createPlansLessonsQuery = `
	INSERT INTO plans_lessons(plan_id, lesson_id)
	VALUES ($1, $2)`
)

func (l *LessonsPostgresStorage) CreateLesson(ctx context.Context, lesson *CreateLesson) (int64, error) {
	const op = "storage.postgresql.lessons.lessons.CreateLesson"

	tx, err := l.db.Begin(ctx)
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

	var lessonID int64
	err = tx.QueryRow(ctx, createLessonQuery,
		lesson.Name,
		lesson.Description,
		lesson.CreatedBy,
		lesson.LastModifiedBy,
		lesson.CreatedAt,
		lesson.Modified,
	).Scan(&lessonID)
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
		createPlansLessonsQuery,
		lesson.PlanID,
		lessonID,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
	}

	return lessonID, nil
}

const getLessonByIDQuery = `
	SELECT 
		l.id AS lesson_id, 
		l.name AS lesson_name, 
		l.description AS description, 
		l.created_by AS lesson_created_by, 
		l.last_modified_by AS lesson_last_modified_by, 
		l.created_at AS lesson_created_at, 
		l.modified AS lesson_modified 
	FROM 
		lessons l
	INNER JOIN
		plans_lessons pl ON l.id = pl.lesson_id
	INNER JOIN
		plans p ON pl.plan_id = p.id
	WHERE 
		lesson_id = $1
		AND pl.plan_id = $2`

func (l *LessonsPostgresStorage) GetLessonByID(ctx context.Context, lessonPlan *GetLesson) (Lesson, error) {
	const op = "storage.postgresql.lessons.lessons.GetLessonByID"

	var lesson DBLesson

	err := l.db.QueryRow(
		ctx,
		getLessonByIDQuery,
		lessonPlan.LessonID,
		lessonPlan.PlanID,
	).Scan(
		&lesson.ID,
		&lesson.Name,
		&lesson.Description,
		&lesson.CreatedBy,
		&lesson.LastModifiedBy,
		&lesson.CreatedAt,
		&lesson.Modified,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return (Lesson)(lesson), fmt.Errorf("%s: %w", op, storage.ErrLessonNotFound)
		}

		return (Lesson)(lesson), fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}

	return (Lesson)(lesson), nil
}

const getLessonsQuery = `
	SELECT
		l.id AS lesson_id,
		l.name AS lesson_name,
		l.description AS description,
		l.created_by AS lesson_created_by,
		l.last_modified_by AS lesson_last_modified_by,
		l.created_at AS lesson_created_at,
		l.modified AS lesson_modified
	FROM 
		lessons l
	INNER JOIN 
		plans_lessons pl ON l.id = pl.lesson_id
	INNER JOIN 
		plans p ON pl.plan_id = p.id
	WHERE pl.plan_id = $1
	ORDER BY l.id
	LIMIT $2 OFFSET $3`

func (l *LessonsPostgresStorage) GetLessons(ctx context.Context, inputParams *GetLessons) ([]Lesson, error) {
	const op = "storage.postgresql.lessons.lessons.GetLessons"

	var lessons []DBLesson

	rows, err := l.db.Query(
		ctx,
		getLessonsQuery,
		inputParams.PlanID,
		inputParams.Limit,
		inputParams.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var lesson DBLesson
		if err := rows.Scan(
			&lesson.ID,
			&lesson.Name,
			&lesson.Description,
			&lesson.CreatedBy,
			&lesson.LastModifiedBy,
			&lesson.CreatedAt,
			&lesson.Modified,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrScanFailed)
		}
		lessons = append(lessons, lesson)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var mappedLessons []Lesson
	for _, lesson := range lessons {
		mappedLessons = append(mappedLessons, Lesson(lesson))
	}

	return mappedLessons, nil
}

const updateLessonQuery = `
	UPDATE lessons l
	SET name = COALESCE($3, l.name),
		description = COALESCE($4, l.description),
	    last_modified_by = $5,
	    modified = now()
	FROM
		plans_lessons pl
	WHERE 
		l.id = $1
		AND pl.plan_id = $2
	RETURNING 
		l.id;`

func (l *LessonsPostgresStorage) UpdateLesson(ctx context.Context, updLesson *UpdateLessonRequest) (int64, error) {
	const op = "storage.postgresql.lesson.lesson.UpdateLesson"

	var id int64

	err := l.db.QueryRow(ctx, updateLessonQuery,
		updLesson.LessonID,
		updLesson.PlanID,
		updLesson.Name,
		updLesson.Description,
		updLesson.LastModifiedBy,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrLessonNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}

	return id, nil
}

const deleteLessonQuery = `
	DELETE FROM lessons l
	USING plans_lessons pl, plans p
	WHERE l.id = pl.lesson_id
		AND pl.plan_id = p.id
		AND l.id = $1
		AND pl.plan_id = $2;`

func (l *LessonsPostgresStorage) DeleteLesson(ctx context.Context, lessonP *DeleteLesson) error {
	const op = "storage.postgresql.lessons.lessons.DeleteLesson"

	res, err := l.db.Exec(
		ctx,
		deleteLessonQuery,
		lessonP.LessonID,
		lessonP.PlanID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrLessonNotFound)
	}

	return nil
}
