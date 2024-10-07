package pages

import (
	"context"
	"fmt"
	"log"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PagesPostgresStorage struct {
	db *pgxpool.Pool
}

func NewPagesStorage(db *pgxpool.Pool) *PagesPostgresStorage {
	return &PagesPostgresStorage{db: db}
}

const (
	createAbstractPageQuery = `
	INSERT INTO pages_abstractpages(lesson_id, created_by, last_modified_by, created_at, modified, content_type)
	VALUES ($1, $2, $3, now(), now(), $4)
	RETURNING id`
	updateAbstractPageQuery = `
	UPDATE pages_abstractpages
	SET
		last_modified_by = $2,
		modified = now()
	WHERE id = $1`
	// createImagePageQuery = `
	// INSERT INTO image_imagepage(abstractpage_id, image_file_url, image_name)
	// VALUES ($1, $2, $3)`
	// createVideoPageQuery = `
	// INSERT INTO video_videopage(abstractpage_id, video_file_url, video_name)
	// VALUES ($1, $2, $3)`
	// createPDFPageQuery = `
	// INSERT INTO pdf_pdfpage(abstractpage_id, pdf_file_url, pdf_name)
	// VALUES ($1, $2, $3)`
)

func (p *PagesPostgresStorage) insertPageSpecific(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) error {
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("insertPageSpecific: %w", err)
	}
	return nil
}

func (p *PagesPostgresStorage) insertAbstractPage(ctx context.Context, page CreatePage) (int64, pgx.Tx, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("create: %w", storage.ErrFailedTransaction)
	}

	commonFields := page.GetCommonFields()
	var pageID int64
	err = tx.QueryRow(
		ctx,
		createAbstractPageQuery,
		commonFields.LessonID,
		commonFields.CreatedBy,
		commonFields.LastModifiedBy,
		commonFields.ContentType,
	).Scan(&pageID)
	if err != nil {
		return 0, nil, fmt.Errorf("create abstract: %w", err)
	}

	return pageID, tx, err
}

func (p *PagesPostgresStorage) updateAbstractPage(ctx context.Context, page UpdatePage) (int64, pgx.Tx, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("update: %w", storage.ErrFailedTransaction)
	}

	commonFields := page.GetCommonFields()
	_, err = tx.Exec(
		ctx,
		updateAbstractPageQuery,
		commonFields.ID,
		commonFields.LastModifiedBy,
	)
	if err != nil {
		return 0, nil, fmt.Errorf("update abstract: %w", err)
	}

	return commonFields.ID, tx, err
}

func (p *PagesPostgresStorage) CreatePage(ctx context.Context, page CreatePage) (int64, error) {
	const op = "storage.postgresql.pages.pages.CreatePage"

	pageID, tx, err := p.insertAbstractPage(ctx, page)
	if err != nil {
		return 0, fmt.Errorf("CreatePage abstract: %w", err)
	}
	defer func(err error) {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Printf("%s: %v", op, storage.ErrRollBack)
			}
		}
	}(err)

	err = p.insertPageSpecific(
		ctx,
		tx,
		page.GetInsertQuery(),
		append([]interface{}{pageID},
			page.GetContentTypeSpecificFields()...)...,
	)
	if err != nil {
		return 0, fmt.Errorf("create page specific: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("create page commit: %w", storage.ErrCommitTransaction)
	}
	return pageID, nil
}

const (
	getImagePageByIDQuery = `
	SELECT 
		ab.id AS abstractpage_id, 
		ab.lesson_id lesson_id, 
		ab.created_by AS created_by, 
		ab.last_modified_by AS last_modified_by, 
		ab.created_at AS created_at, 
		ab.modified AS modified, 
		ab.content_type AS content_type,
		ip.image_file_url AS image_file_url,
		ip.image_name AS image_name
	FROM
		pages_abstractpages ab
	INNER JOIN
		image_imagepage ip ON ab.id =  ip.abstractpage_id
	WHERE abstractpage_id = $1`
	getVideoPageByIDQuery = `
	SELECT 
		ab.id AS abstractpage_id, 
		ab.lesson_id lesson_id, 
		ab.created_by AS created_by, 
		ab.last_modified_by AS last_modified_by, 
		ab.created_at AS created_at, 
		ab.modified AS modified, 
		ab.content_type AS content_type,
		vp.video_file_url AS video_file_url,
		vp.video_name AS video_name
	FROM
		pages_abstractpages ab
	INNER JOIN
		video_videopage vp ON ab.id =  vp.abstractpage_id
	WHERE abstractpage_id = $1`
	getPDFPageByIDQuery = `
	SELECT 
		ab.id AS abstractpage_id, 
		ab.lesson_id lesson_id, 
		ab.created_by AS created_by, 
		ab.last_modified_by AS last_modified_by, 
		ab.created_at AS created_at, 
		ab.modified AS modified, 
		ab.content_type AS content_type,
		pdf.pdf_file_url AS pdf_file_url,
		pdf.pdf_name AS pdf_name
	FROM
		pages_abstractpages ab
	INNER JOIN
		pdf_pdfpage pdf ON ab.id =  pdf.abstractpage_id
	WHERE abstractpage_id = $1`
)

func (p *PagesPostgresStorage) GetPageByID(ctx context.Context, pageID int64, contentType string) (Page, error) {
	const op = "storage.postgresql.pages.pages.GetPageByID"

	var page Page

	switch contentType {
	case "image":
		var dbImagePage DBImagePage
		err := p.db.QueryRow(ctx, getImagePageByIDQuery, pageID).Scan(
			&dbImagePage.ID,
			&dbImagePage.LessonID,
			&dbImagePage.CreatedBy,
			&dbImagePage.LastModifiedBy,
			&dbImagePage.CreatedAt,
			&dbImagePage.Modified,
			&dbImagePage.ContentType,
			&dbImagePage.ImageFileUrl,
			&dbImagePage.ImageName,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrPageNotFound)
		}

		page = &ImagePage{
			BasePage: BasePage{
				ID:             dbImagePage.ID,
				LessonID:       dbImagePage.LessonID,
				CreatedBy:      dbImagePage.CreatedBy,
				LastModifiedBy: dbImagePage.LastModifiedBy,
				CreatedAt:      dbImagePage.CreatedAt,
				Modified:       dbImagePage.Modified,
				ContentType:    dbImagePage.ContentType,
			},
			ImageFileUrl: dbImagePage.ImageFileUrl,
			ImageName:    dbImagePage.ImageName,
		}

	case "video":
		var dbVideoPage DBVideoPage
		query := getVideoPageByIDQuery
		err := p.db.QueryRow(ctx, query, pageID).Scan(
			&dbVideoPage.ID,
			&dbVideoPage.LessonID,
			&dbVideoPage.CreatedBy,
			&dbVideoPage.LastModifiedBy,
			&dbVideoPage.CreatedAt,
			&dbVideoPage.Modified,
			&dbVideoPage.ContentType,
			&dbVideoPage.VideoFileUrl,
			&dbVideoPage.VideoName,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrPageNotFound)
		}

		page = &VideoPage{
			BasePage: BasePage{
				ID:             dbVideoPage.ID,
				LessonID:       dbVideoPage.LessonID,
				CreatedBy:      dbVideoPage.CreatedBy,
				LastModifiedBy: dbVideoPage.LastModifiedBy,
				CreatedAt:      dbVideoPage.CreatedAt,
				Modified:       dbVideoPage.Modified,
				ContentType:    dbVideoPage.ContentType,
			},
			VideoFileUrl: dbVideoPage.VideoFileUrl,
			VideoName:    dbVideoPage.VideoName,
		}

	case "pdf":
		var dbPDFPage DBPDFPage
		err := p.db.QueryRow(ctx, getPDFPageByIDQuery, pageID).Scan(
			&dbPDFPage.ID,
			&dbPDFPage.LessonID,
			&dbPDFPage.CreatedBy,
			&dbPDFPage.LastModifiedBy,
			&dbPDFPage.CreatedAt,
			&dbPDFPage.Modified,
			&dbPDFPage.ContentType,
			&dbPDFPage.PdfFileUrl,
			&dbPDFPage.PdfName,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrPageNotFound)
		}

		page = &PDFPage{
			BasePage: BasePage{
				ID:             dbPDFPage.ID,
				LessonID:       dbPDFPage.LessonID,
				CreatedBy:      dbPDFPage.CreatedBy,
				LastModifiedBy: dbPDFPage.LastModifiedBy,
				CreatedAt:      dbPDFPage.CreatedAt,
				Modified:       dbPDFPage.Modified,
				ContentType:    dbPDFPage.ContentType,
			},
			PdfFileUrl: dbPDFPage.PdfFileUrl,
			PdfName:    dbPDFPage.PdfName,
		}

	default:
		return nil, fmt.Errorf("%s: %w", op, storage.ErrUnContType)
	}

	return page, nil
}

const getPagesQuery = `
	SELECT
		ab.id AS abstractpage_id,
		ab.lesson_id lesson_id,
		ab.created_by AS created_by,
		ab.last_modified_by AS last_modified_by,
		ab.created_at AS created_at,
		ab.modified AS modified,
		ab.content_type AS content_type
	FROM
		pages_abstractpages ab
	INNER JOIN
		lessons l ON ab.lesson_id = l.id
	WHERE l.id = $1
	ORDER BY abstractpage_id
	LIMIT $2 OFFSET $3`

func (p *PagesPostgresStorage) GetPages(ctx context.Context, lessonID int64, limit, offset int64) ([]BasePage, error) {
	const op = "storage.postgresql.pages.pages.GetPages"

	var pages []DBBasePage

	rows, err := p.db.Query(ctx, getPagesQuery, lessonID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var page DBBasePage
		if err := rows.Scan(
			&page.ID,
			&page.LessonID,
			&page.CreatedBy,
			&page.LastModifiedBy,
			&page.CreatedAt,
			&page.Modified,
			&page.ContentType,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrScanFailed)
		}
		pages = append(pages, page)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var mappedPages []BasePage
	for _, page := range pages {
		mappedPages = append(mappedPages, BasePage(page))
	}

	return mappedPages, nil
}

func (p *PagesPostgresStorage) UpdatePage(ctx context.Context, updPage UpdatePage) (int64, error) {
	const op = "storage.postgresql.pages.pages.UpdatePage"

	pageID, tx, err := p.updateAbstractPage(ctx, updPage)
	if err != nil {
		return 0, fmt.Errorf("update page abstract: %w", err)
	}
	defer func(err error) {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Printf("%s: %v", op, storage.ErrRollBack)
			}
		}
	}(err)

	err = p.insertPageSpecific(
		ctx,
		tx,
		updPage.GetUpdateQuery(),
		append([]interface{}{pageID},
			updPage.GetContentTypeSpecificFields()...)...,
	)
	if err != nil {
		return 0, fmt.Errorf("update page specific: %w", storage.ErrInvalidCredentials)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("update page commit: %w", storage.ErrCommitTransaction)
	}

	return pageID, nil
}

const deletePageQuery = `
	DELETE FROM pages_abstractpages
	WHERE id = $1`

func (p *PagesPostgresStorage) DeletePage(ctx context.Context, id int64) error {
	const op = "storage.postgresql.pages.pages.DeletePage"

	res, err := p.db.Exec(ctx, deletePageQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, storage.ErrPageNotFound)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrPageNotFound)
	}

	return nil
}

// func (p *PagesPostgresStorage) UpdateImagePage(ctx context.Context, updPage models.UpdateImagePage) (int64, error) {
// 	const op = "storage.postgresql.pages.pages.UpdateImagePage"

// 	tx, err := p.db.Begin(ctx)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
// 	}
// 	defer func() {
// 		if err != nil {
// 			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
// 				log.Printf("%s: %v", op, storage.ErrRollBack)
// 			}
// 		}
// 	}()

// 	_, err = tx.Exec(
// 		ctx,
// 		updateAbstractPageQuery,
// 		updPage.ID,
// 		updPage.LastModifiedBy,
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	_, err = tx.Exec(
// 		ctx,
// 		updateImagePage,
// 		updPage.ImageFileUrl,
// 		updPage.ImageName,
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	if err = tx.Commit(ctx); err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
// 	}

// 	return updPage.ID, nil
// }

// const getImagePageByIDQuery = `
// 	SELECT
// 		ab.id AS abstractpage_id,
// 		ab.lesson_id lesson_id,
// 		ab.created_by AS created_by,
// 		ab.last_modified_by AS last_modified_by,
// 		ab.created_at AS created_at,
// 		ab.modified AS modified,
// 		ab.content_type AS content_type,
// 		ip.image_file_url,
// 		ip.image_name
// 	FROM
// 		pages_abstractpages ab
// 	INNER JOIN
// 		lessons l ON ab.lesson_id = l.id
// 	INNER JOIN
// 		image_imagepage ip ON ab.id = ip.abstractpage_id
// 	WHERE l.id = $1`

// func (p *PagesPostgresStorage) GetImagePageByID(ctx context.Context, pageID int64) (models.ImagePage, error) {
// 	const op = "storage.postgresql.pages.pages.GetImagePageByID"

// 	var imagePage DBImagePage

// 	err := p.db.QueryRow(ctx, getImagePageByIDQuery, pageID).Scan(
// 		&imagePage.ID,
// 		&imagePage.LessonID,
// 		&imagePage.CreatedBy,
// 		&imagePage.LastModifiedBy,
// 		&imagePage.CreatedAt,
// 		&imagePage.Modified,
// 		&imagePage.ContentType,
// 		&imagePage.ImageFileUrl,
// 		&imagePage.ImageName,
// 	)
// 	if err != nil {
// 		return (models.ImagePage)(imagePage), fmt.Errorf("%s: %w", op, storage.ErrPageNotFound)
// 	}

// 	return (models.ImagePage)(imagePage), nil
// }

// const getVideoPageByIDQuery = `
// SELECT
// 	ab.id AS abstractpage_id,
// 	ab.lesson_id lesson_id,
// 	ab.created_by AS created_by,
// 	ab.last_modified_by AS last_modified_by,
// 	ab.created_at AS created_at,
// 	ab.modified AS modified,
// 	ab.content_type AS content_type,
// 	ip.video_file_url,
// 	ip.video_name
// FROM
// 	pages_abstractpages ab
// INNER JOIN
// 	lessons l ON ab.lesson_id = l.id
// INNER JOIN
// 	video_videopage ip ON ab.id = ip.abstractpage_id
// WHERE l.id = $1`

// func (p *PagesPostgresStorage) GetVideoPageByID(ctx context.Context, pageID int64) (models.VideoPage, error) {
// 	const op = "storage.postgresql.pages.pages.GetVideoPageByID"

// 	var videoPage DBVideoPage

// 	err := p.db.QueryRow(ctx, getVideoPageByIDQuery, pageID).Scan(
// 		&videoPage.ID,
// 		&videoPage.LessonID,
// 		&videoPage.CreatedBy,
// 		&videoPage.LastModifiedBy,
// 		&videoPage.CreatedAt,
// 		&videoPage.Modified,
// 		&videoPage.ContentType,
// 		&videoPage.VideoFileUrl,
// 		&videoPage.VideoName,
// 	)
// 	if err != nil {
// 		return (models.VideoPage)(videoPage), fmt.Errorf("%s: %w", op, storage.ErrPageNotFound)
// 	}

// 	return (models.VideoPage)(videoPage), nil
// }

// const getPDFPageByIDQuery = `
// 	SELECT
// 		ab.id AS abstractpage_id,
// 		ab.lesson_id lesson_id,
// 		ab.created_by AS created_by,
// 		ab.last_modified_by AS last_modified_by,
// 		ab.created_at AS created_at,
// 		ab.modified AS modified,
// 		ab.content_type AS content_type,
// 		ip.pdf_file_url,
// 		ip.pdf_name
// 	FROM
// 		pages_abstractpages ab
// 	INNER JOIN
// 		lessons l ON ab.lesson_id = l.id
// 	INNER JOIN
// 		pdf_pdfpage ip ON ab.id = ip.abstractpage_id
// 	WHERE l.id = $1`

// func (p *PagesPostgresStorage) GetPDFPageByID(ctx context.Context, pageID int64) (models.PDFPage, error) {
// 	const op = "storage.postgresql.pages.pages.GetPDFPageByID"

// 	var pdfPage DBPDFPage

// 	err := p.db.QueryRow(ctx, getPDFPageByIDQuery, pageID).Scan(
// 		&pdfPage.ID,
// 		&pdfPage.LessonID,
// 		&pdfPage.CreatedBy,
// 		&pdfPage.LastModifiedBy,
// 		&pdfPage.CreatedAt,
// 		&pdfPage.Modified,
// 		&pdfPage.ContentType,
// 		&pdfPage.PdfFileUrl,
// 		&pdfPage.PdfName,
// 	)
// 	if err != nil {
// 		return (models.PDFPage)(pdfPage), fmt.Errorf("%s: %w", op, storage.ErrPageNotFound)
// 	}

// 	return (models.PDFPage)(pdfPage), nil
// }

// func (p *PagesPostgresStorage) CreateImagePage(ctx context.Context, imagePage models.CreateImagePage) (int64, error) {
// 	const op = "storage.postgresql.pages.pages.CreateImagePage"

// 	tx, err := p.db.Begin(ctx)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
// 	}
// 	defer func() {
// 		if err != nil {
// 			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
// 				log.Printf("%s: %v", op, storage.ErrRollBack)
// 			}
// 		}
// 	}()

// 	var pageID int64
// 	err = tx.QueryRow(ctx, createAbstractPageQuery,
// 		imagePage.ID,
// 		imagePage.LessonID,
// 		imagePage.CreatedBy,
// 		imagePage.LastModifiedBy,
// 		imagePage.CreatedAt,
// 		imagePage.Modified,
// 		imagePage.ContentType,
// 	).Scan(&pageID)
// 	if err != nil {
// 		var pgErr *pgconn.PgError
// 		if errors.As(err, &pgErr) {
// 			if pgErr.Code == "23505" { // unique violation code
// 				return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
// 			}
// 		}
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	_, err = tx.Exec(
// 		ctx,
// 		createImagePageQuery,
// 		pageID,
// 		imagePage.ImageFileUrl,
// 		imagePage.ImageName,
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
// 	}

// 	if err = tx.Commit(ctx); err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
// 	}

// 	return pageID, nil
// }

// func (p *PagesPostgresStorage) CreateVideoPage(ctx context.Context, videoPage models.CreateVideoPage) (int64, error) {
// 	const op = "storage.postgresql.pages.pages.CreateVideoPage"

// 	tx, err := p.db.Begin(ctx)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
// 	}
// 	defer func() {
// 		if err != nil {
// 			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
// 				log.Printf("%s: %v", op, storage.ErrRollBack)
// 			}
// 		}
// 	}()

// 	var pageID int64
// 	err = tx.QueryRow(ctx, createAbstractPageQuery,
// 		videoPage.ID,
// 		videoPage.LessonID,
// 		videoPage.CreatedBy,
// 		videoPage.LastModifiedBy,
// 		videoPage.CreatedAt,
// 		videoPage.Modified,
// 		videoPage.ContentType,
// 	).Scan(&pageID)
// 	if err != nil {
// 		var pgErr *pgconn.PgError
// 		if errors.As(err, &pgErr) {
// 			if pgErr.Code == "23505" { // unique violation code
// 				return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
// 			}
// 		}
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	_, err = tx.Exec(
// 		ctx,
// 		createVideoPageQuery,
// 		pageID,
// 		videoPage.VideoFileUrl,
// 		videoPage.VideoName,
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
// 	}

// 	if err = tx.Commit(ctx); err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
// 	}

// 	return pageID, nil
// }

// func (p *PagesPostgresStorage) CreatePDFPage(ctx context.Context, pdfPage models.CreatePDFPage) (int64, error) {
// 	const op = "storage.postgresql.pages.pages.CreatePDFPage"

// 	tx, err := p.db.Begin(ctx)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
// 	}
// 	defer func() {
// 		if err != nil {
// 			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
// 				log.Printf("%s: %v", op, storage.ErrRollBack)
// 			}
// 		}
// 	}()

// 	var pageID int64
// 	err = tx.QueryRow(ctx, createAbstractPageQuery,
// 		pdfPage.ID,
// 		pdfPage.LessonID,
// 		pdfPage.CreatedBy,
// 		pdfPage.LastModifiedBy,
// 		pdfPage.CreatedAt,
// 		pdfPage.Modified,
// 		pdfPage.ContentType,
// 	).Scan(&pageID)
// 	if err != nil {
// 		var pgErr *pgconn.PgError
// 		if errors.As(err, &pgErr) {
// 			if pgErr.Code == "23505" { // unique violation code
// 				return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
// 			}
// 		}
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	_, err = tx.Exec(
// 		ctx,
// 		createPDFPageQuery,
// 		pageID,
// 		pdfPage.PdfFileUrl,
// 		pdfPage.PdfName,
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
// 	}

// 	if err = tx.Commit(ctx); err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
// 	}

// 	return pageID, nil
// }

// const getPagesQuery = `
// 	SELECT
// 		ab.id AS abstractpage_id,
// 		ab.lesson_id lesson_id,
// 		ab.created_by AS created_by,
// 		ab.last_modified_by AS last_modified_by,
// 		ab.created_at AS created_at,
// 		ab.modified AS modified,
// 		ab.content_type AS content_type
// 	FROM
// 		pages_abstractpages ab
// 	INNER JOIN
// 		lessons l ON ab.lesson_id = l.id
// 	WHERE l.id = $1
// 	ORDER BY abstractpage_id
// 	LIMIT $2 OFFSET $3`

// func (p *PagesPostgresStorage) GetPages(ctx context.Context, lessonID int64, limit, offset int64) ([]models.Page, error) {
// 	const op = "storage.postgresql.pages.pages.GetPages"

// 	var pages []DBPage

// 	rows, err := p.db.Query(ctx, getPagesQuery, lessonID, limit, offset)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var page DBPage
// 		if err := rows.Scan(
// 			&page.ID,
// 			&page.LessonID,
// 			&page.CreatedBy,
// 			&page.LastModifiedBy,
// 			&page.CreatedAt,
// 			&page.Modified,
// 			&page.ContentType,
// 		); err != nil {
// 			return nil, fmt.Errorf("%s: %w", op, storage.ErrScanFailed)
// 		}
// 		pages = append(pages, page)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	var mappedPages []models.Page
// 	for _, page := range pages {
// 		mappedPages = append(mappedPages, models.Page(page))
// 	}

// 	return mappedPages, nil
// }

// const (
// 	updateAbstractPageQuery = `
// 	UPDATE pages_abstractpages
// 	SET
// 		last_modified_by = $2,
// 		modified = now()
// 	WHERE id = $1`
// 	updateImagePage = `
// 	UPDATE
// 		image_imagepage
// 	SET
// 		image_file_url = COALESCE($2, image_file_url),
// 		image_name = COALESCE($3, image_name)
// 	WHERE abstractpage_id = $1`
// )

// func (p *PagesPostgresStorage) UpdateImagePage(ctx context.Context, updPage models.UpdateImagePage) (int64, error) {
// 	const op = "storage.postgresql.pages.pages.UpdateImagePage"

// 	tx, err := p.db.Begin(ctx)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
// 	}
// 	defer func() {
// 		if err != nil {
// 			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
// 				log.Printf("%s: %v", op, storage.ErrRollBack)
// 			}
// 		}
// 	}()

// 	_, err = tx.Exec(
// 		ctx,
// 		updateAbstractPageQuery,
// 		updPage.ID,
// 		updPage.LastModifiedBy,
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	_, err = tx.Exec(
// 		ctx,
// 		updateImagePage,
// 		updPage.ImageFileUrl,
// 		updPage.ImageName,
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	if err = tx.Commit(ctx); err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
// 	}

// 	return updPage.ID, nil
// }

// const updateVideoPage = `
// 	UPDATE video_videopage
// 	SET
// 		video_file_url = COALESCE($2, video_file_url),
// 		video_name = COALESCE($3, video_name)
// 	WHERE abstractpage_id = $1`

// func (p *PagesPostgresStorage) UpdateVideoPage(ctx context.Context, updPage models.UpdateVideoPage) (int64, error) {
// 	const op = "storage.postgresql.pages.pages.UpdateVideoPage"

// 	tx, err := p.db.Begin(ctx)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
// 	}
// 	defer func() {
// 		if err != nil {
// 			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
// 				log.Printf("%s: %v", op, storage.ErrRollBack)
// 			}
// 		}
// 	}()

// 	_, err = tx.Exec(
// 		ctx,
// 		updateAbstractPageQuery,
// 		updPage.ID,
// 		updPage.LastModifiedBy,
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	_, err = tx.Exec(
// 		ctx,
// 		updateVideoPage,
// 		updPage.VideoFileUrl,
// 		updPage.VideoName,
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	if err = tx.Commit(ctx); err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
// 	}

// 	return updPage.ID, nil
// }

// const updatePDFPage = `
// 	UPDATE pdf_pdfpage
// 	SET
// 		pdf_file_url = COALESCE($2, pdf_file_url),
// 		pdf_name = COALESCE($3, pdf_name)
// 	WHERE abstractpage_id = $1`

// func (p *PagesPostgresStorage) UpdatePDFPage(ctx context.Context, updPage models.UpdatePDFPage) (int64, error) {
// 	const op = "storage.postgresql.pages.pages.UpdatePDFPage"

// 	tx, err := p.db.Begin(ctx)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrFailedTransaction)
// 	}
// 	defer func() {
// 		if err != nil {
// 			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
// 				log.Printf("%s: %v", op, storage.ErrRollBack)
// 			}
// 		}
// 	}()

// 	_, err = tx.Exec(
// 		ctx,
// 		updateAbstractPageQuery,
// 		updPage.ID,
// 		updPage.LastModifiedBy,
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	_, err = tx.Exec(
// 		ctx,
// 		updatePDFPage,
// 		updPage.PDFFileUrl,
// 		updPage.PDFName,
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	if err = tx.Commit(ctx); err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, storage.ErrCommitTransaction)
// 	}

// 	return updPage.ID, nil
// }

// const deletePageQuery = `
// 	DELETE FROM pages_abstractpages
// 	WHERE id = $1`

// func (p *PagesPostgresStorage) DeletePage(ctx context.Context, id int64) error {
// 	const op = "storage.postgresql.pages.pages.DeletePage"

// 	res, err := p.db.Exec(ctx, deletePageQuery, id)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, storage.ErrPageNotFound)
// 	}

// 	if res.RowsAffected() == 0 {
// 		return fmt.Errorf("%s: %w", op, storage.ErrPageNotFound)
// 	}

// 	return nil
// }
