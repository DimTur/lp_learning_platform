package lp_handlers

import (
	"context"
	"errors"
	"time"

	pageserv "github.com/DimTur/lp_learning_platform/internal/services/page"
	pagestore "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/pages"
	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) CreateImagePage(ctx context.Context, req *lpv1.CreateImagePageRequest) (*lpv1.CreateImagePageResponse, error) {
	pageID, err := s.pageHandlers.CreateImagePage(ctx, &pagestore.CreateImagePage{
		CreateBasePage: pagestore.CreateBasePage{
			LessonID:       req.Base.GetLessonId(),
			CreatedBy:      req.Base.GetCreatedBy(),
			LastModifiedBy: req.Base.GetCreatedBy(),
			ContentType:    "image",
		},
		ImageName:    req.GetImageName(),
		ImageFileUrl: req.GetImageFileUrl(),
	})
	if err != nil {
		switch {
		case errors.Is(err, pageserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.CreateImagePageResponse{
		Id: pageID,
	}, nil
}

func (s *serverAPI) CreateVideoPage(ctx context.Context, req *lpv1.CreateVideoPageRequest) (*lpv1.CreateVideoPageResponse, error) {
	pageID, err := s.pageHandlers.CreateVideoPage(ctx, &pagestore.CreateVideoPage{
		CreateBasePage: pagestore.CreateBasePage{
			LessonID:       req.Base.GetLessonId(),
			CreatedBy:      req.Base.GetCreatedBy(),
			LastModifiedBy: req.Base.GetCreatedBy(),
			ContentType:    "video",
		},
		VideoName:    req.GetVideoName(),
		VideoFileUrl: req.GetVideoFileUrl(),
	})
	if err != nil {
		switch {
		case errors.Is(err, pageserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.CreateVideoPageResponse{
		Id: pageID,
	}, nil
}

func (s *serverAPI) CreatePDFPage(ctx context.Context, req *lpv1.CreatePDFPageRequest) (*lpv1.CreatePDFPageResponse, error) {
	pageID, err := s.pageHandlers.CreatePDFPage(ctx, &pagestore.CreatePDFPage{
		CreateBasePage: pagestore.CreateBasePage{
			LessonID:       req.Base.GetLessonId(),
			CreatedBy:      req.Base.GetCreatedBy(),
			LastModifiedBy: req.Base.GetCreatedBy(),
			ContentType:    "pdf",
		},
		PdfName:    req.GetPdfName(),
		PdfFileUrl: req.GetPdfFileUrl(),
	})
	if err != nil {
		switch {
		case errors.Is(err, pageserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.CreatePDFPageResponse{
		Id: pageID,
	}, nil
}

func (s *serverAPI) GetImagePage(ctx context.Context, req *lpv1.GetImagePageRequest) (*lpv1.GetImagePageResponse, error) {
	page, err := s.pageHandlers.GetImagePage(ctx, &pagestore.GetPage{
		PageID:   req.GetPageId(),
		LessonID: req.GetLessonId(),
	})
	if err != nil {
		switch {
		case errors.Is(err, pageserv.ErrPageNotFound):
			return nil, status.Error(codes.NotFound, "image page not found")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.GetImagePageResponse{
		Base: &lpv1.BasePage{
			Id:             page.BasePage.ID,
			LessonId:       page.BasePage.LessonID,
			CreatedBy:      page.BasePage.CreatedBy,
			LastModifiedBy: page.BasePage.LastModifiedBy,
			CreatedAt:      page.BasePage.CreatedAt.Format(time.RFC3339),
			Modified:       page.BasePage.Modified.Format(time.RFC3339),
			ContentType:    convertToContentType(page.BasePage.ContentType),
		},
		ImageFileUrl: page.ImageFileUrl,
		ImageName:    page.ImageName,
	}, nil
}

func (s *serverAPI) GetVideoPage(ctx context.Context, req *lpv1.GetVideoPageRequest) (*lpv1.GetVideoPageResponse, error) {
	page, err := s.pageHandlers.GetVideoPage(ctx, &pagestore.GetPage{
		PageID:   req.GetPageId(),
		LessonID: req.GetLessonId(),
	})
	if err != nil {
		switch {
		case errors.Is(err, pageserv.ErrPageNotFound):
			return nil, status.Error(codes.NotFound, "image page not found")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.GetVideoPageResponse{
		Base: &lpv1.BasePage{
			Id:             page.BasePage.ID,
			LessonId:       page.BasePage.LessonID,
			CreatedBy:      page.BasePage.CreatedBy,
			LastModifiedBy: page.BasePage.LastModifiedBy,
			CreatedAt:      page.BasePage.CreatedAt.Format(time.RFC3339),
			Modified:       page.BasePage.Modified.Format(time.RFC3339),
			ContentType:    convertToContentType(page.BasePage.ContentType),
		},
		VideoFileUrl: page.VideoFileUrl,
		VideoName:    page.VideoName,
	}, nil
}

func (s *serverAPI) GetPDFPage(ctx context.Context, req *lpv1.GetPDFPageRequest) (*lpv1.GetPDFPageResponse, error) {
	page, err := s.pageHandlers.GetPDFPage(ctx, &pagestore.GetPage{
		PageID:   req.GetPageId(),
		LessonID: req.GetLessonId(),
	})
	if err != nil {
		switch {
		case errors.Is(err, pageserv.ErrPageNotFound):
			return nil, status.Error(codes.NotFound, "image page not found")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.GetPDFPageResponse{
		Base: &lpv1.BasePage{
			Id:             page.BasePage.ID,
			LessonId:       page.BasePage.LessonID,
			CreatedBy:      page.BasePage.CreatedBy,
			LastModifiedBy: page.BasePage.LastModifiedBy,
			CreatedAt:      page.BasePage.CreatedAt.Format(time.RFC3339),
			Modified:       page.BasePage.Modified.Format(time.RFC3339),
			ContentType:    convertToContentType(page.BasePage.ContentType),
		},
		PdfFileUrl: page.PdfFileUrl,
		PdfName:    page.PdfName,
	}, nil
}

func (s *serverAPI) GetPages(ctx context.Context, req *lpv1.GetPagesRequest) (*lpv1.GetPagesResponse, error) {
	pages, err := s.pageHandlers.GetPages(ctx, &pagestore.GetPages{
		LessonID: req.GetLessonId(),
		Limit:    req.GetLimit(),
		Offset:   req.GetOffset(),
	})
	if err != nil {
		switch {
		case errors.Is(err, pageserv.ErrPageNotFound):
			return nil, status.Error(codes.NotFound, "pages not found")
		case errors.Is(err, pageserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var responsePages []*lpv1.BasePage
	for _, page := range pages {
		responsePages = append(responsePages, &lpv1.BasePage{
			Id:             page.ID,
			LessonId:       page.LessonID,
			CreatedBy:      page.CreatedBy,
			LastModifiedBy: page.LastModifiedBy,
			CreatedAt:      page.CreatedAt.Format(time.RFC3339),
			Modified:       page.Modified.Format(time.RFC3339),
			ContentType:    convertToContentType(page.ContentType),
		})
	}

	return &lpv1.GetPagesResponse{
		Pages: responsePages,
	}, nil
}

func (s *serverAPI) UpdateImagePage(ctx context.Context, req *lpv1.UpdateImagePageRequest) (*lpv1.UpdateImagePageResponse, error) {
	id, err := s.pageHandlers.UpdateImagePage(ctx, pagestore.UpdateImagePage{
		UpdateBasePage: pagestore.UpdateBasePage{
			ID:             req.Base.GetId(),
			LastModifiedBy: req.Base.GetLastModifiedBy(),
			ContentType:    "image",
		},
		ImageFileUrl: req.GetImageFileUrl(),
		ImageName:    req.GetImageName(),
	})
	if err != nil {
		switch {
		case errors.Is(err, pageserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		case errors.Is(err, pageserv.ErrPageNotFound):
			return nil, status.Error(codes.NotFound, "image page not found")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.UpdateImagePageResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) UpdateVideoPage(ctx context.Context, req *lpv1.UpdateVideoPageRequest) (*lpv1.UpdateVideoPageResponse, error) {
	id, err := s.pageHandlers.UpdateVideoPage(ctx, pagestore.UpdateVideoPage{
		UpdateBasePage: pagestore.UpdateBasePage{
			ID:             req.Base.GetId(),
			LastModifiedBy: req.Base.GetLastModifiedBy(),
			ContentType:    "video",
		},
		VideoFileUrl: req.GetVideoFileUrl(),
		VideoName:    req.GetVideoName(),
	})
	if err != nil {
		switch {
		case errors.Is(err, pageserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		case errors.Is(err, pageserv.ErrPageNotFound):
			return nil, status.Error(codes.NotFound, "video page not found")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.UpdateVideoPageResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) UpdatePDFPage(ctx context.Context, req *lpv1.UpdatePDFPageRequest) (*lpv1.UpdatePDFPageResponse, error) {
	id, err := s.pageHandlers.UpdatePDFPage(ctx, pagestore.UpdatePDFPage{
		UpdateBasePage: pagestore.UpdateBasePage{
			ID:             req.Base.GetId(),
			LastModifiedBy: req.Base.GetLastModifiedBy(),
			ContentType:    "video",
		},
		PdfFileUrl: req.GetPdfFileUrl(),
		PdfName:    req.GetPdfName(),
	})
	if err != nil {
		switch {
		case errors.Is(err, pageserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		case errors.Is(err, pageserv.ErrPageNotFound):
			return nil, status.Error(codes.NotFound, "pdf page not found")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.UpdatePDFPageResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) DeletePage(ctx context.Context, req *lpv1.DeletePageRequest) (*lpv1.DeletePageResponse, error) {
	err := s.pageHandlers.DeletePage(ctx, &pagestore.DeletePage{
		PageID:   req.GetPageId(),
		LessonID: req.GetLessonId(),
	})
	if err != nil {
		if errors.Is(err, pageserv.ErrPageNotFound) {
			return nil, status.Error(codes.NotFound, "page not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.DeletePageResponse{
		Success: true,
	}, nil
}

func convertToContentType(contentTypeStr string) lpv1.ContentType {
	switch contentTypeStr {
	case "image":
		return lpv1.ContentType_IMAGE
	case "video":
		return lpv1.ContentType_VIDEO
	case "pdf":
		return lpv1.ContentType_PDF
	case "question":
		return lpv1.ContentType_PDF
	default:
		return lpv1.ContentType_CONTENT_TYPE_UNSPECIFIED
	}
}
