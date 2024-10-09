package questions

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionsPostgresStorage struct {
	db *pgxpool.Pool
}

func NewQuestionsStorage(db *pgxpool.Pool) *QuestionsPostgresStorage {
	return &QuestionsPostgresStorage{db: db}
}

const (
	createAbstractPageQuery = `
	INSERT INTO pages_abstractpages(lesson_id, created_by, last_modified_by, created_at, modified, content_type)
	VALUES ($1, $2, $3, now(), now(), $4)
	RETURNING id`
	createAbstractQuestion = `
	INSERT INTO question_abstractquestion(question_type)
	VALUES	($1)
	RETURNING id`
	createQuestionPage = `
	INSERT INTO question_questionpage(abstractpage_id, question_id)
	VALUES ($1, $2)`
	createMultichoiceQuestion = `
	INSERT INTO question_multichoicequestion(
		question_abstractquestion_id,
		question,
		option_a,
		option_b,
		option_c,
		option_d,
		option_e,
		answer
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
)

func (q *QuestionsPostgresStorage) CreateQuestionPage(ctx context.Context, questionPage CreateQuestionPage) (int64, error) {
	const op = "storage.postgresql.pages.pages.CreateQuestionPage"

	tx, err := q.db.Begin(ctx)
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

	var abstrPageID int64
	err = tx.QueryRow(
		ctx,
		createAbstractPageQuery,
		questionPage.LessonID,
		questionPage.CreatedBy,
		questionPage.LastModifiedBy,
		questionPage.ContentType,
	).Scan(&abstrPageID)
	if err != nil {
		return q.checkPgError(err, op)
	}

	var quePageID int64
	err = tx.QueryRow(
		ctx,
		createAbstractQuestion,
		questionPage.QuestionType,
	).Scan(&quePageID)
	if err != nil {
		return q.checkPgError(err, op)
	}

	_, err = tx.Exec(
		ctx,
		createQuestionPage,
		abstrPageID,
		quePageID,
	)
	if err != nil {
		return q.checkPgError(err, op)
	}

	_, err = tx.Exec(
		ctx,
		createMultichoiceQuestion,
		quePageID,
		questionPage.Question,
		questionPage.OptionA,
		questionPage.OptionB,
		questionPage.OptionC,
		questionPage.OptionD,
		questionPage.OptionE,
		questionPage.Answer,
	)
	if err != nil {
		return q.checkPgError(err, op)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
	}

	return abstrPageID, nil
}

const getQuestionPageByIDQuery = `
	SELECT 
		ab.id AS abstractpage_id, 
		ab.lesson_id lesson_id, 
		ab.created_by AS created_by, 
		ab.last_modified_by AS last_modified_by, 
		ab.created_at AS created_at, 
		ab.modified AS modified, 
		ab.content_type AS content_type,
		aq.question_type AS question_type,
		mq.question AS question,
		mq.option_a AS option_a,
		mq.option_b AS option_b,
		mq.option_c AS option_c,
		mq.option_d AS option_d,
		mq.option_e AS option_e,
		mq.answer AS answer
	FROM
		pages_abstractpages ab
	INNER JOIN
		question_questionpage qp ON ab.id = qp.abstractpage_id
	INNER JOIN
		question_abstractquestion aq ON qp.question_id = aq.id
	INNER JOIN
		question_multichoicequestion mq ON aq.id = mq.question_abstractquestion_id
	WHERE abstractpage_id = $1`

func (q *QuestionsPostgresStorage) GetQuestionPageByID(ctx context.Context, pageID int64) (QuestionPage, error) {
	const op = "storage.postgresql.pages.pages.GetQuestionPageByID"

	var questionPage DBQuestionPage

	err := q.db.QueryRow(ctx, getQuestionPageByIDQuery, pageID).Scan(
		&questionPage.ID,
		&questionPage.LessonID,
		&questionPage.CreatedBy,
		&questionPage.LastModifiedBy,
		&questionPage.CreatedAt,
		&questionPage.Modified,
		&questionPage.ContentType,
		&questionPage.QuestionType,
		&questionPage.Question,
		&questionPage.OptionA,
		&questionPage.OptionB,
		&questionPage.OptionC,
		&questionPage.OptionD,
		&questionPage.OptionE,
		&questionPage.Answer,
	)
	if err != nil {
		return (QuestionPage)(questionPage), fmt.Errorf("%s: %w", op, storage.ErrPageNotFound)
	}

	return (QuestionPage)(questionPage), nil
}

const (
	updateAbstractPageQuery = `
	UPDATE pages_abstractpages
	SET
		last_modified_by = $2,
		modified = now()
	WHERE id = $1`
	updateQuestionPageQuery = `
	UPDATE
		question_multichoicequestion mq
	SET
		question = COALESCE($2, question),
		option_a = COALESCE($3, option_a),
		option_b = COALESCE($4, option_b),
		option_c = COALESCE($5, option_c),
		option_d = COALESCE($6, option_d),
		option_e = COALESCE($7, option_e),
		answer = COALESCE($8, answer)
	FROM
    	question_abstractquestion aq
	INNER JOIN
		question_questionpage qp ON aq.id = qp.question_id
	WHERE
		mq.question_abstractquestion_id = aq.id
	AND qp.abstractpage_id = $1`
)

func (q *QuestionsPostgresStorage) UpdateQuestionPage(ctx context.Context, updPage UpdateQuestionPage) (int64, error) {
	const op = "storage.postgresql.pages.pages.UpdateQuestionPageByID"

	tx, err := q.db.Begin(ctx)
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

	_, err = tx.Exec(
		ctx,
		updateAbstractPageQuery,
		updPage.ID,
		updPage.LastModifiedBy,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(
		ctx,
		updateQuestionPageQuery,
		updPage.ID,
		updPage.Question,
		updPage.OptionA,
		updPage.OptionB,
		updPage.OptionC,
		updPage.OptionD,
		updPage.OptionE,
		updPage.Answer,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
	}

	return updPage.ID, nil
}

func (q *QuestionsPostgresStorage) checkPgError(err error, op string) (int64, error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}
	return 0, fmt.Errorf("%s: %w", op, err)
}
