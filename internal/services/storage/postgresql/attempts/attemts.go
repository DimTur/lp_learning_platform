package attempts

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

type AttemptsPostgresStorage struct {
	db *pgxpool.Pool
}

func NewAttemptsStorage(db *pgxpool.Pool) *AttemptsPostgresStorage {
	return &AttemptsPostgresStorage{db: db}
}

const (
	createLessonAttemptQuery = `
	INSERT INTO attempt_lessonattempt(
		lesson_id, 
		plan_id, 
		channel_id, 
		user_id,
		start_time
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING 
		id`
)

func (a *AttemptsPostgresStorage) CreateLessonAttempt(ctx context.Context, lAttempt *CreateLessonAttempt) (int64, error) {
	const op = "storage.postgresql.attempts.attempts.CreateLessonAttempt"

	var id int64

	err := a.db.QueryRow(
		ctx,
		createLessonAttemptQuery,
		lAttempt.LessonID,
		lAttempt.PlanID,
		lAttempt.ChannelID,
		lAttempt.UserID,
		lAttempt.StartTime,
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
	INSERT INTO 
		question_abstractquestionattempt(question_type, page_attempt_id, modified)
	VALUES ($1, $2, $3)
	RETURNING 
		id`
	createAbstractPageAttemptQuery = `
	INSERT INTO 
		pages_abstractpageattempt(lesson_attempt_id, content_type, modified)
	VALUES ($1, $2, $3)
	RETURNING 
		id`
	createQuestionAttemptQuery = `
	INSERT INTO 
		question_questionpageattempt(page_id, question_attempt_id)
	VALUES ($1, $2)
	RETURNING
		id, page_id`
)

func (a *AttemptsPostgresStorage) CreateQuestionPageAttempts(ctx context.Context, attempt CreateQuestionPageAttemptNew) (*CreateQuestionPageAttemptResp, error) {
	const op = "storage.postgresql.attempts.attempts.CreateQuestionPageAttempts"

	tx, err := a.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Printf("%s: %v", op, storage.ErrRollBack)
			}
		}
	}()

	var abPageAttID int64
	err = tx.QueryRow(
		ctx,
		createAbstractPageAttemptQuery,
		attempt.LessonAttemptID,
		attempt.ContentType,
		attempt.Modified,
	).Scan(&abPageAttID)
	if err != nil {
		return nil, a.checkPgError(err, op)
	}

	var abQAttID int64
	err = tx.QueryRow(
		ctx,
		createAbstractQuestionAttemptQuery,
		attempt.QuestionType,
		abPageAttID,
		attempt.Modified,
	).Scan(&abQAttID)
	if err != nil {
		return nil, a.checkPgError(err, op)
	}

	pageAttempt := &CreateQuestionPageAttemptResp{}
	err = tx.QueryRow(
		ctx,
		createQuestionAttemptQuery,
		attempt.PageID,
		abQAttID,
	).Scan(&pageAttempt.ID, &pageAttempt.PageID)
	if err != nil {
		return nil, a.checkPgError(err, op)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
	}

	return pageAttempt, nil
}

const checkLessonAttemptQuery = `
	SELECT 
		la.id
	FROM 
		attempt_lessonattempt la
	WHERE
		la.user_id = $1
		AND la.lesson_id = $2
		AND la.channel_id = $3
		AND la.plan_id = $4
		AND la.is_complete = false
	LIMIT 1;`

func (a *AttemptsPostgresStorage) CheckLessonAttempt(ctx context.Context, lessonAttempt *GetQuestionPageAttempts) (int64, error) {
	const op = "storage.postgresql.attempts.attempts.GetLessonPagesAttempts"

	var id int64
	err := a.db.QueryRow(
		ctx,
		checkLessonAttemptQuery,
		lessonAttempt.UserID,
		lessonAttempt.LessonID,
		lessonAttempt.ChannelID,
		lessonAttempt.PlanID,
	).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, a.checkPgError(err, op)
	}

	return id, nil
}

const getLessonPagesAttemptsQuery = `
	SELECT
		qpa.id AS id,
		qpa.page_id AS page_id,
		qpa.question_attempt_id AS qustion_attempt_id,
		qpa.user_answer AS user_answer,
		la.id AS lesson_attempt_id,
		la.user_id AS user_id,
		la.lesson_id AS lesson_id,
		la.start_time AS start_time,
		aqa.is_successful AS is_correct
	FROM 
		question_questionpageattempt qpa
	INNER JOIN
		question_questionpage qp ON qpa.page_id = qp.id
	INNER JOIN
		question_abstractquestionattempt aqa ON qpa.question_attempt_id = aqa.id 
	INNER JOIN
		pages_abstractpageattempt apa ON aqa.page_attempt_id = apa.id
	INNER JOIN
		attempt_lessonattempt la ON apa.lesson_attempt_id = la.id
	WHERE
		la.user_id = $1
		AND la.lesson_id = $2
		AND la.is_complete = false;`

func (a *AttemptsPostgresStorage) GetLessonPagesAttempts(ctx context.Context, lessonAttempt *GetQuestionPageAttempts) ([]QuestionPageAttempt, error) {
	const op = "storage.postgresql.attempts.attempts.GetLessonPagesAttempts"

	var attempts []DBQuestionPageAttempt
	rows, err := a.db.Query(
		ctx,
		getLessonPagesAttemptsQuery,
		lessonAttempt.UserID,
		lessonAttempt.LessonID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrPageAttemtsNotFound)
	}
	defer rows.Close()

	for rows.Next() {
		var attempt DBQuestionPageAttempt
		if err := rows.Scan(
			&attempt.ID,
			&attempt.PageID,
			&attempt.LessonAttemptID,
			&attempt.IsCorrect,
			&attempt.UserAnswer,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrScanFailed)
		}
		attempts = append(attempts, attempt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var mappedPageAttempts []QuestionPageAttempt
	for _, attempt := range attempts {
		mappedPageAttempts = append(mappedPageAttempts, QuestionPageAttempt(attempt))
	}

	return mappedPageAttempts, nil
}

const getCurrentAnswerForAttemptQuery = `
	SELECT
		mq.answer
	FROM 
		question_questionpage qp
	INNER JOIN
		question_abstractquestion aq ON qp.question_id = aq.id
	INNER JOIN
		question_multichoicequestion mq ON aq.id = mq.question_abstractquestion_id
	WHERE
		qp.id = $1;`

func (a *AttemptsPostgresStorage) GetCurrentAnswerForAttempt(ctx context.Context, pageID int64) (string, error) {
	const op = "storage.postgresql.attempts.attempts.GetCurrentAnswerForAttempt"

	var answer DBAnswer
	err := a.db.QueryRow(
		ctx,
		getCurrentAnswerForAttemptQuery,
		pageID,
	).Scan(
		&answer.Answer,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrAnswerNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}

	return answer.Answer, nil
}

const (
	updQPAttemptQuery = `
	UPDATE 
		question_questionpageattempt qpa
	SET
		user_answer = COALESCE($2, qpa.user_answer)
	WHERE
		qpa.id = $1
	RETURNING
		qpa.id;`

	updAQAttemptQuery = `
	UPDATE
		question_abstractquestionattempt aqa
	SET
		modified = COALESCE($2, aqa.modified),
		is_successful = COALESCE($3, aqa.is_successful)
	WHERE
		EXISTS (
			SELECT 1
			FROM question_questionpageattempt qpa
			WHERE qpa.id = $1
				AND qpa.question_attempt_id = aqa.id
		)
	RETURNING
		aqa.id;`

	updAPAttemptQuery = `
	UPDATE
		pages_abstractpageattempt apa
	SET
		modified = COALESCE($2, apa.modified)
	WHERE
		EXISTS (
			SELECT 1
			FROM question_abstractquestionattempt aqa
			WHERE aqa.id = $1 
				AND aqa.page_attempt_id = apa.id
		);`
)

func (a *AttemptsPostgresStorage) UpdatePageAttempt(ctx context.Context, updPAttempt *UpdatePageAttempt) error {
	const op = "storage.postgresql.attempts.attempts.UpdatePageAttempt"

	tx, err := a.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Printf("%s: %v", op, storage.ErrRollBack)
			}
		}
	}()

	var qPAttemptID int64
	err = tx.QueryRow(
		ctx,
		updQPAttemptQuery,
		updPAttempt.QPAttemptID,
		updPAttempt.UserAnswer,
	).Scan(&qPAttemptID)
	if err != nil {
		return a.checkPgError(err, op)
	}

	var aQAttemptID int64
	err = tx.QueryRow(
		ctx,
		updAQAttemptQuery,
		qPAttemptID,
		updPAttempt.Modified,
		updPAttempt.IsSuccessful,
	).Scan(&aQAttemptID)
	if err != nil {
		return a.checkPgError(err, op)
	}

	_, err = tx.Exec(
		ctx,
		updAPAttemptQuery,
		aQAttemptID,
		updPAttempt.Modified,
	)
	if err != nil {
		return a.checkPgError(err, op)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
	}

	return nil
}

const updLessonAttemptQuery = `
	UPDATE
		attempt_lessonattempt la
	SET
		end_time = COALESCE($2, la.end_time),
		is_complete = COALESCE($3, la.is_complete),
		is_successful = COALESCE($4, la.is_successful),
		percentage_score = COALESCE($5, la.percentage_score)
	WHERE
		la.id = $1
	RETURNING
		la.id;`

func (a *AttemptsPostgresStorage) UpdateLessonAttempt(ctx context.Context, updLAttempt *UpdateLessonAttempt) (int64, error) {
	const op = "storage.postgresql.attempts.attempts.UpdateLessonAttempt"

	var lAttemptID int64
	err := a.db.QueryRow(
		ctx,
		updLessonAttemptQuery,
		updLAttempt.LessonAttemptID,
		updLAttempt.EndTime,
		updLAttempt.IsComplete,
		updLAttempt.IsSuccessful,
		updLAttempt.PercentageScore,
	).Scan(&lAttemptID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrLessonAttemtNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}

	return lAttemptID, nil
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
		ap.lesson_id = $1 
		AND	ap.content_type = 'question' 
		AND aq.question_type = 'multichoice'`

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

const getLessonAttemptsQuery = `
	SELECT
		la.id AS id,
		la.user_id AS user_id,
		la.lesson_id AS lesson_id,
		la.plan_id AS plan_id,
		la.channel_id AS channel_id,
		la.start_time AS start_time,
		la.end_time AS end_time,
		la.is_complete AS is_complete,
		la.is_successful AS is_successful,
		la.percentage_score AS percentage_score
	FROM 
		attempt_lessonattempt la
	WHERE
		la.user_id = $1
		AND la.lesson_id = $2
	ORDER BY
		la.end_time
	LIMIT $3 OFFSET $4;`

func (a *AttemptsPostgresStorage) GetLessonAttempts(ctx context.Context, inputParams *GetLessonAttempts) (*GetLessonAttemptsResp, error) {
	const op = "storage.postgresql.attempts.attempts.GetLessonAttempts"

	var attempts GetLessonAttemptsResp

	rows, err := a.db.Query(
		ctx,
		getLessonAttemptsQuery,
		inputParams.UserID,
		inputParams.LessonID,
		inputParams.Limit,
		inputParams.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrLessonAttemtNotFound)
	}
	defer rows.Close()

	for rows.Next() {
		var attempt DBLessonAttempt
		if err := rows.Scan(
			&attempt.ID,
			&attempt.UserID,
			&attempt.LessonID,
			&attempt.PlanID,
			&attempt.ChannelID,
			&attempt.StartTime,
			&attempt.EndTime,
			&attempt.IsComplete,
			&attempt.IsSuccessful,
			&attempt.PercentageScore,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrScanFailed)
		}
		attempts.LessonAttempts = append(attempts.LessonAttempts, LessonAttempt(attempt))
	}

	return &attempts, nil
}

const checkPermissionForUserQuery = `
	SELECT EXISTS (
		SELECT 1
		FROM attempt_lessonattempt la
		WHERE la.id = $1 AND la.user_id = $2
	);`

func (a *AttemptsPostgresStorage) CheckPermissionForUser(ctx context.Context, userAtt *PermissionForUser) (bool, error) {
	const op = "storage.postgresql.attempts.attempts.CheckPermissionForUser"

	var exists bool
	err := a.db.QueryRow(
		ctx,
		checkPermissionForUserQuery,
		userAtt.LessonAttemptID,
		userAtt.UserID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (a *AttemptsPostgresStorage) checkPgError(err error, op string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}
	return fmt.Errorf("%s: %w", op, err)
}
