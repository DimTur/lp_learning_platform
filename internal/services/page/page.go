package page

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/pages"
	"github.com/DimTur/lp_learning_platform/internal/utils"
	"github.com/go-playground/validator/v10"
)

type PageSaver interface {
	CreatePage(ctx context.Context, page pages.CreatePage) (int64, error)
	UpdatePage(ctx context.Context, updPage pages.UpdatePage) (int64, error)
}

type PageProvider interface {
	GetPageByID(ctx context.Context, pageID int64, contentType string) (pages.Page, error)
	GetPages(ctx context.Context, lessonID int64, limit, offset int64) ([]pages.BasePage, error)
}
type PageDel interface {
	DeletePage(ctx context.Context, id int64) error
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

func (ph *PageHandlers) CreatePage(ctx context.Context, page pages.CreatePage) (int64, error) {
	const op = "page.CreatePage"

	log := ph.log.With(
		slog.String("op", op),
		slog.String("page with type", page.GetCommonFields().ContentType),
	)

	// Validation
	err := ph.validator.Struct(page)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	commonFields := page.GetCommonFields()

	log.Info("creating page with", slog.String("content_type", commonFields.ContentType))

	id, err := ph.pageSaver.CreatePage(ctx, page)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			ph.log.Warn("invalid arguments", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save page", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (ph *PageHandlers) GetPage(ctx context.Context, pageID int64, contentType string) (pages.Page, error) {
	const op = "page.GetPage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("pageID", pageID),
	)

	log.Info("getting page")

	page, err := ph.pageProvider.GetPageByID(ctx, pageID, contentType)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrPageNotFound):
			ph.log.Warn("page not found", slog.String("err", err.Error()))
			return nil, ErrPageNotFound
		case errors.Is(err, storage.ErrUnContType):
			ph.log.Warn("page has unsupported content_type", slog.String("err", err.Error()))
			return nil, ErrUnContType
		default:
			log.Error("failed to get page", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return page, nil
}

// GetPages gets pages and returns them.
func (ph *PageHandlers) GetPages(ctx context.Context, lessonID int64, limit, offset int64) ([]pages.BasePage, error) {
	const op = "page.GetPages"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("getting pages included in lesson with id", lessonID),
	)

	log.Info("getting pages")

	// Validation
	params := utils.PaginationQueryParams{
		Limit:  limit,
		Offset: offset,
	}
	params.SetDefaults()

	if err := ph.validator.Struct(params); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	var pages []pages.BasePage
	pages, err := ph.pageProvider.GetPages(ctx, lessonID, params.Limit, params.Offset)
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

func (ph *PageHandlers) UpdatePage(ctx context.Context, updPage pages.UpdatePage) (int64, error) {
	const op = "page.UpdatePage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("updating page with id:", updPage.GetCommonFields().ID),
	)

	log.Info("updating page")

	// Validation
	err := ph.validator.Struct(updPage)
	if err != nil {
		log.Warn("validation failed", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	id, err := ph.pageSaver.UpdatePage(ctx, updPage)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			ph.log.Warn("invalid credentials", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to update image page", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("page updated with ", slog.Int64("page", id))

	return id, nil
}

// DeletePage
func (ph *PageHandlers) DeletePage(ctx context.Context, pageID int64) error {
	const op = "page.DeletePage"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("page id", pageID),
	)

	log.Info("deleting page with: ", slog.Int64("pageID", pageID))

	err := ph.pageDel.DeletePage(ctx, pageID)
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

// // UpdateImagePage performs a partial update
// func (ph *PageHandlers) UpdateImagePage(ctx context.Context, updPage models.UpdateImagePage) (int64, error) {
// 	const op = "page.UpdateImagePage"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.Int64("updating image page with id: ", updPage.ID),
// 	)

// 	log.Info("updating image page")

// 	// Validation
// 	err := ph.validator.Struct(updPage)
// 	if err != nil {
// 		log.Warn("validation failed", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
// 	}

// 	id, err := ph.pageSaver.UpdateImagePage(ctx, updPage)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrInvalidCredentials) {
// 			ph.log.Warn("invalid credentials", slog.String("err", err.Error()))
// 			return 0, fmt.Errorf("%s: %w", op, err)
// 		}

// 		log.Error("failed to update image page", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	log.Info("image page updated with ", slog.Int64("page", id))

// 	return id, nil
// }

// // GetImagePageByID gets image page by ID and returns it.
// func (ph *PageHandlers) GetImagePageByID(ctx context.Context, pageID int64) (models.ImagePage, error) {
// 	const op = "page.GetImagePageByID"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.Int64("pageID", pageID),
// 	)

// 	log.Info("getting image page")

// 	var imagePage models.ImagePage
// 	imagePage, err := ph.pageProvider.GetImagePageByID(ctx, pageID)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrPageNotFound) {
// 			ph.log.Warn("image page not found", slog.String("err", err.Error()))
// 			return imagePage, ErrPageNotFound
// 		}

// 		log.Error("failed to get image page", slog.String("err", err.Error()))
// 		return imagePage, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return imagePage, nil
// }

// // CreateImagePage creates new image page in the system and returns page ID.
// func (ph *PageHandlers) CreateImagePage(ctx context.Context, imagePage models.CreateImagePage) (int64, error) {
// 	const op = "page.CreateImagePage"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.String("page with type", imagePage.ContentType),
// 	)

// 	// Validation
// 	err := ph.validator.Struct(imagePage)
// 	if err != nil {
// 		log.Warn("invalid parameters", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
// 	}

// 	now := time.Now()
// 	imagePage.CreatedAt = now
// 	imagePage.Modified = now

// 	log.Info("creating image page")

// 	id, err := ph.pageSaver.CreateImagePage(ctx, imagePage)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrInvalidCredentials) {
// 			ph.log.Warn("invalid arguments", slog.String("err", err.Error()))
// 			return 0, fmt.Errorf("%s: %w", op, err)
// 		}

// 		log.Error("failed to save image page", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return id, nil
// }

// // CreateVideoPage creates new video page in the system and returns page ID.
// func (ph *PageHandlers) CreateVideoPage(ctx context.Context, videoPage models.CreateVideoPage) (int64, error) {
// 	const op = "page.CreateVideoPage"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.String("page with type", videoPage.ContentType),
// 	)

// 	// Validation
// 	err := ph.validator.Struct(videoPage)
// 	if err != nil {
// 		log.Warn("invalid parameters", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
// 	}

// 	now := time.Now()
// 	videoPage.CreatedAt = now
// 	videoPage.Modified = now

// 	log.Info("creating video page")

// 	id, err := ph.pageSaver.CreateVideoPage(ctx, videoPage)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrInvalidCredentials) {
// 			ph.log.Warn("invalid arguments", slog.String("err", err.Error()))
// 			return 0, fmt.Errorf("%s: %w", op, err)
// 		}

// 		log.Error("failed to save video page", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return id, nil
// }

// // CreatePDFPage creates new pdf page in the system and returns page ID.
// func (ph *PageHandlers) CreatePDFPage(ctx context.Context, pdfPage models.CreatePDFPage) (int64, error) {
// 	const op = "page.CreatePDFPage"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.String("page with type", pdfPage.ContentType),
// 	)

// 	// Validation
// 	err := ph.validator.Struct(pdfPage)
// 	if err != nil {
// 		log.Warn("invalid parameters", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
// 	}

// 	now := time.Now()
// 	pdfPage.CreatedAt = now
// 	pdfPage.Modified = now

// 	log.Info("creating pdf page")

// 	id, err := ph.pageSaver.CreatePDFPage(ctx, pdfPage)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrInvalidCredentials) {
// 			ph.log.Warn("invalid arguments", slog.String("err", err.Error()))
// 			return 0, fmt.Errorf("%s: %w", op, err)
// 		}

// 		log.Error("failed to save pdf page", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return id, nil
// }

// // GetImagePageByID gets image page by ID and returns it.
// func (ph *PageHandlers) GetImagePageByID(ctx context.Context, pageID int64) (models.ImagePage, error) {
// 	const op = "page.GetImagePageByID"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.Int64("pageID", pageID),
// 	)

// 	log.Info("getting image page")

// 	var imagePage models.ImagePage
// 	imagePage, err := ph.pageProvider.GetImagePageByID(ctx, pageID)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrPageNotFound) {
// 			ph.log.Warn("image page not found", slog.String("err", err.Error()))
// 			return imagePage, ErrPageNotFound
// 		}

// 		log.Error("failed to get image page", slog.String("err", err.Error()))
// 		return imagePage, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return imagePage, nil
// }

// // GetVideoPageByID gets image page by ID and returns it.
// func (ph *PageHandlers) GetVideoPageByID(ctx context.Context, pageID int64) (models.VideoPage, error) {
// 	const op = "page.GetVideoPageByID"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.Int64("pageID", pageID),
// 	)

// 	log.Info("getting video page")

// 	var videoPage models.VideoPage
// 	videoPage, err := ph.pageProvider.GetVideoPageByID(ctx, pageID)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrPageNotFound) {
// 			ph.log.Warn("video page not found", slog.String("err", err.Error()))
// 			return videoPage, ErrPageNotFound
// 		}

// 		log.Error("failed to get video page", slog.String("err", err.Error()))
// 		return videoPage, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return videoPage, nil
// }

// // GetPDFPageByID gets image page by ID and returns it.
// func (ph *PageHandlers) GetPDFPageByID(ctx context.Context, pageID int64) (models.PDFPage, error) {
// 	const op = "page.GetPDFPageByID"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.Int64("pageID", pageID),
// 	)

// 	log.Info("getting pdf page")

// 	var pdfPage models.PDFPage
// 	pdfPage, err := ph.pageProvider.GetPDFPageByID(ctx, pageID)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrPageNotFound) {
// 			ph.log.Warn("pdf page not found", slog.String("err", err.Error()))
// 			return pdfPage, ErrPageNotFound
// 		}

// 		log.Error("failed to get pdf page", slog.String("err", err.Error()))
// 		return pdfPage, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return pdfPage, nil
// }

// // GetPages gets pages and returns them.
// func (ph *PageHandlers) GetPages(ctx context.Context, lessonID int64, limit, offset int64) ([]models.Page, error) {
// 	const op = "page.GetPages"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.Int64("getting pages included in lesson with id", lessonID),
// 	)

// 	log.Info("getting lessons")

// 	// Validation
// 	params := utils.PaginationQueryParams{
// 		Limit:  limit,
// 		Offset: offset,
// 	}

// 	if err := ph.validator.Struct(params); err != nil {
// 		log.Warn("invalid parameters", slog.String("err", err.Error()))
// 		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
// 	}

// 	var pages []models.Page
// 	pages, err := ph.pageProvider.GetPages(ctx, lessonID, limit, offset)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrPageNotFound) {
// 			ph.log.Warn("pages not found", slog.String("err", err.Error()))
// 			return pages, fmt.Errorf("%s: %w", op, ErrPageNotFound)
// 		}

// 		log.Error("failed to get pages", slog.String("err", err.Error()))
// 		return pages, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return pages, nil
// }

// // UpdateImagePage performs a partial update
// func (ph *PageHandlers) UpdateImagePage(ctx context.Context, updPage models.UpdateImagePage) (int64, error) {
// 	const op = "page.UpdateImagePage"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.Int64("updating image page with id: ", updPage.ID),
// 	)

// 	log.Info("updating image page")

// 	// Validation
// 	err := ph.validator.Struct(updPage)
// 	if err != nil {
// 		log.Warn("validation failed", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
// 	}

// 	id, err := ph.pageSaver.UpdateImagePage(ctx, updPage)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrInvalidCredentials) {
// 			ph.log.Warn("invalid credentials", slog.String("err", err.Error()))
// 			return 0, fmt.Errorf("%s: %w", op, err)
// 		}

// 		log.Error("failed to update image page", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	log.Info("image page updated with ", slog.Int64("page", id))

// 	return id, nil
// }

// // UpdateVideoPage performs a partial update
// func (ph *PageHandlers) UpdateVideoPage(ctx context.Context, updPage models.UpdateVideoPage) (int64, error) {
// 	const op = "page.UpdateVideoPage"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.Int64("updating video page with id: ", updPage.ID),
// 	)

// 	log.Info("updating video page")

// 	// Validation
// 	err := ph.validator.Struct(updPage)
// 	if err != nil {
// 		log.Warn("validation failed", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
// 	}

// 	id, err := ph.pageSaver.UpdateVideoPage(ctx, updPage)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrInvalidCredentials) {
// 			ph.log.Warn("invalid credentials", slog.String("err", err.Error()))
// 			return 0, fmt.Errorf("%s: %w", op, err)
// 		}

// 		log.Error("failed to update video page", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	log.Info("video page updated with ", slog.Int64("page", id))

// 	return id, nil
// }

// // UpdatePDFPage performs a partial update
// func (ph *PageHandlers) UpdatePDFPage(ctx context.Context, updPage models.UpdatePDFPage) (int64, error) {
// 	const op = "page.UpdatePDFPage"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.Int64("updating pdf page with id: ", updPage.ID),
// 	)

// 	log.Info("updating pdf page")

// 	// Validation
// 	err := ph.validator.Struct(updPage)
// 	if err != nil {
// 		log.Warn("validation failed", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
// 	}

// 	id, err := ph.pageSaver.UpdatePDFPage(ctx, updPage)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrInvalidCredentials) {
// 			ph.log.Warn("invalid credentials", slog.String("err", err.Error()))
// 			return 0, fmt.Errorf("%s: %w", op, err)
// 		}

// 		log.Error("failed to update pdf page", slog.String("err", err.Error()))
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	log.Info("pdf page updated with ", slog.Int64("page", id))

// 	return id, nil
// }

// // DeletePage
// func (ph *PageHandlers) DeletePage(ctx context.Context, pageID int64) error {
// 	const op = "page.DeletePage"

// 	log := ph.log.With(
// 		slog.String("op", op),
// 		slog.Int64("page id", pageID),
// 	)

// 	log.Info("deleting page with: ", slog.Int64("pageID", pageID))

// 	err := ph.pageDel.DeletePage(ctx, pageID)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrPageNotFound) {
// 			ph.log.Warn("page not found", slog.String("err", err.Error()))
// 			return fmt.Errorf("%s: %w", op, ErrPageNotFound)
// 		}

// 		log.Error("failed to delete page", slog.String("err", err.Error()))
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }
