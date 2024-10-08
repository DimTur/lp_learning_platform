package attempts

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AttemptsPostgresStorage struct {
	db *pgxpool.Pool
}

func NewAttemptsStorage(db *pgxpool.Pool) *AttemptsPostgresStorage {
	return &AttemptsPostgresStorage{db: db}
}

const (
	createLessonAttemptQuery = `
	INSERT INTO attempt_lessonattempt(lesson_id, plan_id, channel_id, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
)

func (a *AttemptsPostgresStorage) CreateLessonAttempt(ctx context.Context, lAttempt CreateLessonAttempt) (int64, error) {
	const op = "storage.postgresql.attempts.attempts.CreateLessonAttempt"

	var id int64

	err := a.db.QueryRow(
		ctx,
		createLessonAttemptQuery,
		lAttempt.LessonID,
		lAttempt.PlanId,
		lAttempt.ChannelID,
		lAttempt.UserID,
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

const (
	createAbstractPageAttemptQuery = `
	INSERT INTO pages_abstractpageattempt(lesson_attempt_id, content_type)
	VALUES ($1, $2)
	RETURNING id`
)

func (a *AttemptsPostgresStorage) CreateAbstractPageAttempt(ctx context.Context, pAttempt CreateAbstractPageAttempt) (int64, error) {
	const op = "storage.postgresql.attempts.attempts.CreateAbstractPageAttempt"

	var id int64

	err := a.db.QueryRow(
		ctx,
		createAbstractPageAttemptQuery,
		pAttempt.LessonAttemptID,
		pAttempt.ContentType,
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

const (
	createAbstractQuestionAttemptQuery = `
	INSERT INTO question_abstractquestionattempt(question_type, page_attempt_id)
	VALUES ($1, $2)
	RETURNING id`
)

func (a *AttemptsPostgresStorage) CreateAbstractQuestionAttempt(ctx context.Context, qAttempt CreateAbstractQuestionAttempt) (int64, error) {
	const op = "storage.postgresql.attempts.attempts.CreateAbstractQuestionAttempt"

	var id int64

	err := a.db.QueryRow(
		ctx,
		createAbstractQuestionAttemptQuery,
		qAttempt.QuestionType,
		qAttempt.PageAttemptID,
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

const (
	createQuestionAttemptQuery = `
	INSERT INTO question_questionpageattempt(page_id, page_attempt_id)
	VALUES ($1, $2)
	RETURNING id`
)

func (a *AttemptsPostgresStorage) CreateQuestionAttempt(ctx context.Context, qPageAttempt CreateQuestionPageAttempt) error {
	const op = "storage.postgresql.attempts.attempts.CreateAbstractQuestionAttempt"

	var id int64

	err := a.db.QueryRow(
		ctx,
		createQuestionAttemptQuery,
		qPageAttempt.PageID,
		qPageAttempt.PageAttemptID,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique violation code
				return fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

const getQuestionPagesQuery = `
	SELECT 
		ap.content_type AS content_type,
		aq.question_type AS question_type,
		qp.id AS question_questionpage_id
	FROM 
		pages_abstractpages ap
	INNER JOIN
		question_questionpage qp ON ap.id = qp.abstractpage_id
	INNER JOIN
		question_abstractquestion aq ON qp.question_id = aq.id
	WHERE 
		ap.lesson_id = $1 AND
		ap.content_type = 'question' AND
		aq.question_type = 'multichoice'`

func (a *AttemptsPostgresStorage) GetQuestionPages(ctx context.Context, lessonID int64) ([]QuestionPage, error) {
	const op = "storage.postgresql.attempts.attempts.GetQuestionPages"

	var qPages []DBQuestionPage

	rows, err := a.db.Query(ctx, getQuestionPagesQuery, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}
	defer rows.Close()

	for rows.Next() {
		var qPage DBQuestionPage
		if err := rows.Scan(
			&qPage.ContentType,
			&qPage.QuestionType,
			&qPage.QuestionPageID,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrScanFailed)
		}
		qPages = append(qPages, qPage)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var mappedQPages []QuestionPage
	for _, qPage := range qPages {
		mappedQPages = append(mappedQPages, QuestionPage(qPage))
	}

	return mappedQPages, nil
}
