package page

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/pages"
	"github.com/go-playground/validator/v10"
)

type PageSaver interface {
	CreatePage(ctx context.Context, page pages.CreatePage) (int64, error)
	UpdatePage(ctx context.Context, updPage pages.UpdatePage) (int64, error)
}

type PageProvider interface {
	GetImagePage(ctx context.Context, pageLesson *pages.GetPage) (pages.Page, error)
	GetVideoPage(ctx context.Context, pageLesson *pages.GetPage) (pages.Page, error)
	GetPDFPage(ctx context.Context, pageLesson *pages.GetPage) (pages.Page, error)
	GetPages(ctx context.Context, inputParams *pages.GetPages) ([]pages.BasePage, error)
}
type PageDel interface {
	DeletePage(ctx context.Context, pageLesson *pages.DeletePage) error
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPageID      = errors.New("invalid page id")
	ErrPageExitsts        = errors.New("page already exists")
	ErrPageNotFound       = errors.New("page not found")
	ErrUnContType         = errors.New("unsupported content type")
)

type PageHandlers struct {
	log          *slog.Logger
	validator    *validator.Validate
	pageSaver    PageSaver
	pageProvider PageProvider
	pageDel      PageDel
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	pageSaver PageSaver,
	pageProvider PageProvider,
	pageDel PageDel,
) *PageHandlers {
	return &PageHandlers{
		log:          log,
		validator:    validator,
		pageSaver:    pageSaver,
		pageProvider: pageProvider,
		pageDel:      pageDel,
	}
}

func (ph *PageHandlers) CreateImagePage(ctx context.Context, imagePage *pages.CreateImagePage) (int64, error) {
	const op = "page.CreateImagePage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("lesson_id", imagePage.GetCommonFields().LessonID),
		slog.String("page_type", imagePage.GetCommonFields().ContentType),
	)

	// Validation
	err := ph.validator.Struct(imagePage)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("start creating image page")

	id, err := ph.pageSaver.CreatePage(ctx, imagePage)

	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			ph.log.Warn("invalid arguments", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save image page", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (ph *PageHandlers) CreatePDFPage(ctx context.Context, pdfPage *pages.CreatePDFPage) (int64, error) {
	const op = "page.CreatePDFPage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("lesson_id", pdfPage.GetCommonFields().LessonID),
		slog.String("page_type", pdfPage.GetCommonFields().ContentType),
	)

	// Validation
	err := ph.validator.Struct(pdfPage)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("start creating pdf page")

	id, err := ph.pageSaver.CreatePage(ctx, pdfPage)

	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			ph.log.Warn("invalid arguments", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save pdf page", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (ph *PageHandlers) CreateVideoPage(ctx context.Context, videoPage *pages.CreateVideoPage) (int64, error) {
	const op = "page.CreateVideoPage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("lesson_id", videoPage.GetCommonFields().LessonID),
		slog.String("page_type", videoPage.GetCommonFields().ContentType),
	)

	// Validation
	err := ph.validator.Struct(videoPage)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("start creating video page")

	id, err := ph.pageSaver.CreatePage(ctx, videoPage)

	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			ph.log.Warn("invalid arguments", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save video page", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (ph *PageHandlers) GetImagePage(ctx context.Context, pageLesson *pages.GetPage) (*pages.ImagePage, error) {
	const op = "page.GetImagePage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("page_id", pageLesson.PageID),
		slog.Int64("lesson_id", pageLesson.LessonID),
	)

	log.Info("getting page")

	page, err := ph.pageProvider.GetImagePage(ctx, pageLesson)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrPageNotFound):
			ph.log.Warn("image page not found", slog.String("err", err.Error()))
			return nil, ErrPageNotFound
		default:
			log.Error("failed to get image page", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	imagePage := &pages.ImagePage{
		BasePage: pages.BasePage{
			ID:             page.GetCommonFields().ID,
			LessonID:       page.GetCommonFields().LessonID,
			CreatedBy:      page.GetCommonFields().CreatedBy,
			LastModifiedBy: page.GetCommonFields().LastModifiedBy,
			CreatedAt:      page.GetCommonFields().CreatedAt,
			Modified:       page.GetCommonFields().Modified,
			ContentType:    page.GetCommonFields().ContentType,
		},
		ImageFileUrl: page.GetContentTypeSpecificFields()[0].(string),
		ImageName:    page.GetContentTypeSpecificFields()[1].(string),
	}

	return imagePage, nil
}

func (ph *PageHandlers) GetVideoPage(ctx context.Context, pageLesson *pages.GetPage) (*pages.VideoPage, error) {
	const op = "page.GetVideoPage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("page_id", pageLesson.PageID),
		slog.Int64("lesson_id", pageLesson.LessonID),
	)

	log.Info("getting page")

	page, err := ph.pageProvider.GetVideoPage(ctx, pageLesson)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrPageNotFound):
			ph.log.Warn("video page not found", slog.String("err", err.Error()))
			return nil, ErrPageNotFound
		default:
			log.Error("failed to get video page", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	videoPage := &pages.VideoPage{
		BasePage: pages.BasePage{
			ID:             page.GetCommonFields().ID,
			LessonID:       page.GetCommonFields().LessonID,
			CreatedBy:      page.GetCommonFields().CreatedBy,
			LastModifiedBy: page.GetCommonFields().LastModifiedBy,
			CreatedAt:      page.GetCommonFields().CreatedAt,
			Modified:       page.GetCommonFields().Modified,
			ContentType:    page.GetCommonFields().ContentType,
		},
		VideoFileUrl: page.GetContentTypeSpecificFields()[0].(string),
		VideoName:    page.GetContentTypeSpecificFields()[1].(string),
	}

	return videoPage, nil
}

func (ph *PageHandlers) GetPDFPage(ctx context.Context, pageLesson *pages.GetPage) (*pages.PDFPage, error) {
	const op = "page.GetPDFPage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("page_id", pageLesson.PageID),
		slog.Int64("lesson_id", pageLesson.LessonID),
	)

	log.Info("getting page")

	page, err := ph.pageProvider.GetPDFPage(ctx, pageLesson)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrPageNotFound):
			ph.log.Warn("pdf page not found", slog.String("err", err.Error()))
			return nil, ErrPageNotFound
		default:
			log.Error("failed to get pdf page", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	pdfPage := &pages.PDFPage{
		BasePage: pages.BasePage{
			ID:             page.GetCommonFields().ID,
			LessonID:       page.GetCommonFields().LessonID,
			CreatedBy:      page.GetCommonFields().CreatedBy,
			LastModifiedBy: page.GetCommonFields().LastModifiedBy,
			CreatedAt:      page.GetCommonFields().CreatedAt,
			Modified:       page.GetCommonFields().Modified,
			ContentType:    page.GetCommonFields().ContentType,
		},
		PdfFileUrl: page.GetContentTypeSpecificFields()[0].(string),
		PdfName:    page.GetContentTypeSpecificFields()[1].(string),
	}

	return pdfPage, nil
}

// GetPages gets pages and returns them.
func (ph *PageHandlers) GetPages(ctx context.Context, inputParams *pages.GetPages) ([]pages.BasePage, error) {
	const op = "page.GetPages"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("lesson_id", inputParams.LessonID),
	)

	log.Info("getting pages")

	// Validation
	params := pages.GetPages{
		LessonID: inputParams.LessonID,
		Limit:    inputParams.Limit,
		Offset:   inputParams.Offset,
	}
	params.SetDefaults()

	if err := ph.validator.Struct(params); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	var pages []pages.BasePage
	pages, err := ph.pageProvider.GetPages(ctx, inputParams)
	if err != nil {
		if errors.Is(err, storage.ErrPageNotFound) {
			ph.log.Warn("pages not found", slog.String("err", err.Error()))
			return pages, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		}

		log.Error("failed to get pages", slog.String("err", err.Error()))
		return pages, fmt.Errorf("%s: %w", op, err)
	}

	return pages, nil
}

func (ph *PageHandlers) UpdateImagePage(ctx context.Context, updPage pages.UpdateImagePage) (int64, error) {
	const op = "page.UpdateImagePage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("page_id:", updPage.GetCommonFields().ID),
	)

	log.Info("updating image page")

	// Validation
	err := ph.validator.Struct(updPage)
	if err != nil {
		log.Warn("validation failed", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	id, err := ph.pageSaver.UpdatePage(ctx, &updPage)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidCredentials):
			ph.log.Warn("invalid credentials", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, storage.ErrPageNotFound):
			ph.log.Warn("image page not found", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			log.Error("failed to update image page", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("image page updated with ", slog.Int64("page_id", id))

	return id, nil
}

func (ph *PageHandlers) UpdateVideoPage(ctx context.Context, updPage pages.UpdateVideoPage) (int64, error) {
	const op = "page.UpdateVideoPage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("page_id:", updPage.GetCommonFields().ID),
	)

	log.Info("updating video page")

	// Validation
	err := ph.validator.Struct(updPage)
	if err != nil {
		log.Warn("validation failed", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	id, err := ph.pageSaver.UpdatePage(ctx, &updPage)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidCredentials):
			ph.log.Warn("invalid credentials", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, storage.ErrPageNotFound):
			ph.log.Warn("video page not found", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			log.Error("failed to update video page", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("page updated with ", slog.Int64("page", id))

	return id, nil
}

func (ph *PageHandlers) UpdatePDFPage(ctx context.Context, updPage pages.UpdatePDFPage) (int64, error) {
	const op = "page.UpdatePDFPage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("page_id:", updPage.GetCommonFields().ID),
	)

	log.Info("updating PDF page")

	// Validation
	err := ph.validator.Struct(updPage)
	if err != nil {
		log.Warn("validation failed", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	id, err := ph.pageSaver.UpdatePage(ctx, &updPage)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidCredentials):
			ph.log.Warn("invalid credentials", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, storage.ErrPageNotFound):
			ph.log.Warn("pdf page not found", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			log.Error("failed to update pdf page", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("page updated with ", slog.Int64("page", id))

	return id, nil
}

// DeletePage
func (ph *PageHandlers) DeletePage(ctx context.Context, pageLesson *pages.DeletePage) error {
	const op = "page.DeletePage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("page_id", pageLesson.PageID),
		slog.Int64("lesson_id", pageLesson.LessonID),
	)

	log.Info("start deleting page")

	err := ph.pageDel.DeletePage(ctx, pageLesson)
	if err != nil {
		if errors.Is(err, storage.ErrPageNotFound) {
			ph.log.Warn("page not found", slog.String("err", err.Error()))
			return fmt.Errorf("%s: %w", op, ErrPageNotFound)
		}

		log.Error("failed to delete page", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
