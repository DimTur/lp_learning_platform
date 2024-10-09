package lp_handlers

import (
	"context"
	"errors"
	"fmt"

	pageserv "github.com/DimTur/lp_learning_platform/internal/services/page"
	pagestore "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/pages"
	lpv1 "github.com/DimTur/lp_learning_platform/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serverAPI) CreatePage(ctx context.Context, req *lpv1.CreatePageRequest) (*lpv1.CreatePageResponse, error) {
	var page pagestore.CreatePage

	switch pageReq := req.GetPage().(type) {
	case *lpv1.CreatePageRequest_ImagePage:
		page = &pagestore.CreateImagePage{
			CreateBasePage: pagestore.CreateBasePage{
				LessonID:       pageReq.ImagePage.Base.LessonId,
				CreatedBy:      pageReq.ImagePage.Base.CreatedBy,
				LastModifiedBy: pageReq.ImagePage.Base.LastModifiedBy,
				ContentType:    "image",
			},
			ImageFileUrl: pageReq.ImagePage.ImageFileUrl,
			ImageName:    pageReq.ImagePage.ImageName,
		}
	case *lpv1.CreatePageRequest_VideoPage:
		page = &pagestore.CreateVideoPage{
			CreateBasePage: pagestore.CreateBasePage{
				LessonID:       pageReq.VideoPage.Base.LessonId,
				CreatedBy:      pageReq.VideoPage.Base.CreatedBy,
				LastModifiedBy: pageReq.VideoPage.Base.LastModifiedBy,
				ContentType:    "video",
			},
			VideoFileUrl: pageReq.VideoPage.VideoFileUrl,
			VideoName:    pageReq.VideoPage.VideoName,
		}
	case *lpv1.CreatePageRequest_PdfPage:
		page = &pagestore.CreatePDFPage{
			CreateBasePage: pagestore.CreateBasePage{
				LessonID:       pageReq.PdfPage.Base.LessonId,
				CreatedBy:      pageReq.PdfPage.Base.CreatedBy,
				LastModifiedBy: pageReq.PdfPage.Base.LastModifiedBy,
				ContentType:    "pdf",
			},
			PdfFileUrl: pageReq.PdfPage.PdfFileUrl,
			PdfName:    pageReq.PdfPage.PdfName,
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported page type")
	}

	pageID, err := s.pageHandlers.CreatePage(ctx, page)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create page: %v", err)
	}

	return &lpv1.CreatePageResponse{
		Id: pageID,
	}, nil
}

func (s *serverAPI) GetPage(ctx context.Context, req *lpv1.GetPageRequest) (*lpv1.GetPageResponse, error) {

	contentType, err := ContentTypeToString(req.GetContentType())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	page, err := s.pageHandlers.GetPage(ctx, req.GetId(), contentType)
	if err != nil {
		switch {
		case errors.Is(err, pageserv.ErrPageNotFound):
			return nil, status.Error(codes.NotFound, "page not found")
		case errors.Is(err, pageserv.ErrUnContType):
			return nil, status.Error(codes.InvalidArgument, "unsupported content type")
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get page: %v", err))
		}
	}

	var response lpv1.GetPageResponse
	switch p := page.(type) {
	case *pagestore.ImagePage:
		response.Page = &lpv1.GetPageResponse_ImagePage{
			ImagePage: &lpv1.ImagePage{
				Base: &lpv1.BasePage{
					Id:             p.ID,
					LessonId:       p.LessonID,
					CreatedBy:      p.CreatedBy,
					LastModifiedBy: p.LastModifiedBy,
					CreatedAt:      timestamppb.New(p.CreatedAt),
					Modified:       timestamppb.New(p.Modified),
					ContentType:    lpv1.ContentType_IMAGE,
				},
				ImageFileUrl: p.ImageFileUrl,
				ImageName:    p.ImageName,
			},
		}
	case *pagestore.VideoPage:
		response.Page = &lpv1.GetPageResponse_VideoPage{
			VideoPage: &lpv1.VideoPage{
				Base: &lpv1.BasePage{
					Id:             p.ID,
					LessonId:       p.LessonID,
					CreatedBy:      p.CreatedBy,
					LastModifiedBy: p.LastModifiedBy,
					CreatedAt:      timestamppb.New(p.CreatedAt),
					Modified:       timestamppb.New(p.Modified),
					ContentType:    lpv1.ContentType_VIDEO,
				},
				VideoFileUrl: p.VideoFileUrl,
				VideoName:    p.VideoName,
			},
		}
	case *pagestore.PDFPage:
		response.Page = &lpv1.GetPageResponse_PdfPage{
			PdfPage: &lpv1.PDFPage{
				Base: &lpv1.BasePage{
					Id:             p.ID,
					LessonId:       p.LessonID,
					CreatedBy:      p.CreatedBy,
					LastModifiedBy: p.LastModifiedBy,
					CreatedAt:      timestamppb.New(p.CreatedAt),
					Modified:       timestamppb.New(p.Modified),
					ContentType:    lpv1.ContentType_PDF,
				},
				PdfFileUrl: p.PdfFileUrl,
				PdfName:    p.PdfName,
			},
		}
	default:
		return nil, status.Error(codes.Internal, "unknown page type")
	}

	return &response, nil
}

func (s *serverAPI) GetPages(ctx context.Context, req *lpv1.GetPagesRequest) (*lpv1.GetPagesResponse, error) {
	pages, err := s.pageHandlers.GetPages(ctx, req.GetLessonId(), req.GetLimit(), req.GetOffset())
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
			CreatedAt:      timestamppb.New(page.CreatedAt),
			Modified:       timestamppb.New(page.Modified),
			ContentType:    convertToContentType(page.ContentType),
		})
	}

	return &lpv1.GetPagesResponse{
		Pages: responsePages,
	}, nil
}

func (s *serverAPI) UpdatePage(ctx context.Context, req *lpv1.UpdatePageRequest) (*lpv1.UpdatePageResponse, error) {
	var page pagestore.UpdatePage

	switch pageReq := req.GetPage().(type) {
	case *lpv1.UpdatePageRequest_ImagePage:
		page = &pagestore.UpdateImagePage{
			UpdateBasePage: pagestore.UpdateBasePage{
				ID:             pageReq.ImagePage.Base.GetId(),
				LastModifiedBy: pageReq.ImagePage.Base.GetLastModifiedBy(),
				ContentType:    "image",
			},
			ImageFileUrl: pageReq.ImagePage.GetImageFileUrl(),
			ImageName:    pageReq.ImagePage.GetImageName(),
		}
	case *lpv1.UpdatePageRequest_VideoPage:
		page = &pagestore.UpdateVideoPage{
			UpdateBasePage: pagestore.UpdateBasePage{
				ID:             pageReq.VideoPage.Base.GetId(),
				LastModifiedBy: pageReq.VideoPage.Base.GetLastModifiedBy(),
				ContentType:    "video",
			},
			VideoFileUrl: pageReq.VideoPage.GetVideoFileUrl(),
			VideoName:    pageReq.VideoPage.GetVideoName(),
		}
	case *lpv1.UpdatePageRequest_PdfPage:
		page = &pagestore.UpdatePDFPage{
			UpdateBasePage: pagestore.UpdateBasePage{
				ID:             pageReq.PdfPage.Base.GetId(),
				LastModifiedBy: pageReq.PdfPage.Base.GetLastModifiedBy(),
				ContentType:    "pdf",
			},
			PdfFileUrl: pageReq.PdfPage.PdfFileUrl,
			PdfName:    pageReq.PdfPage.GetPdfName(),
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported page type")
	}

	pageID, err := s.pageHandlers.UpdatePage(ctx, page)
	if err != nil {
		switch {
		case errors.Is(err, pageserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.UpdatePageResponse{
		Id:      pageID,
		Success: true,
	}, nil
}

func (s *serverAPI) DeletePage(ctx context.Context, req *lpv1.DeletePageRequest) (*lpv1.DeletePageResponse, error) {
	pageId := req.GetId()

	err := s.pageHandlers.DeletePage(ctx, pageId)
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

func ContentTypeToString(contentType lpv1.ContentType) (string, error) {
	switch contentType {
	case lpv1.ContentType_IMAGE:
		return "image", nil
	case lpv1.ContentType_VIDEO:
		return "video", nil
	case lpv1.ContentType_PDF:
		return "pdf", nil
	default:
		return "unknown", fmt.Errorf("unsupported content type: %s", contentType)
	}
}

func convertToContentType(contentTypeStr string) lpv1.ContentType {
	switch contentTypeStr {
	case "image":
		return lpv1.ContentType_IMAGE
	case "video":
		return lpv1.ContentType_VIDEO
	case "pdf":
		return lpv1.ContentType_PDF
	default:
		return lpv1.ContentType_CONTENT_TYPE_UNSPECIFIED
	}
}

// func (s *serverAPI) CreateImagePage(ctx context.Context, req *lpv1.CreateImagePageRequest) (*lpv1.CreateImagePageResponse, error) {
// 	imagePage := models.CreateImagePage{
// 		LessonID:       req.GetLessonId(),
// 		CreatedBy:      req.GetCreatedBy(),
// 		LastModifiedBy: req.GetCreatedBy(),
// 		ContentType:    req.GetContentType().String(),
// 		ImageFileUrl:   req.GetImageFileUrl(),
// 		ImageName:      req.GetImageName(),
// 	}

// 	pageID, err := s.pageHandlers.CreateImagePage(ctx, imagePage)
// 	if err != nil {
// 		if errors.Is(err, lessonserv.ErrInvalidCredentials) {
// 			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
// 		}

// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &lpv1.CreateImagePageResponse{
// 		Id: pageID,
// 	}, nil
// }

// func (s *serverAPI) CreateVideoPage(ctx context.Context, req *lpv1.CreateVideoPageRequest) (*lpv1.CreateVideoPageResponse, error) {
// 	videoPage := models.CreateVideoPage{
// 		LessonID:       req.GetLessonId(),
// 		CreatedBy:      req.GetCreatedBy(),
// 		LastModifiedBy: req.GetCreatedBy(),
// 		ContentType:    req.GetContentType().String(),
// 		VideoFileUrl:   req.GetVideoFileUrl(),
// 		VideoName:      req.GetVideoName(),
// 	}

// 	pageID, err := s.pageHandlers.CreateVideoPage(ctx, videoPage)
// 	if err != nil {
// 		if errors.Is(err, lessonserv.ErrInvalidCredentials) {
// 			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
// 		}

// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &lpv1.CreateVideoPageResponse{
// 		Id: pageID,
// 	}, nil
// }

// func (s *serverAPI) CreatePDFPage(ctx context.Context, req *lpv1.CreatePDFPageRequest) (*lpv1.CreatePDFPageResponse, error) {
// 	imagePage := models.CreatePDFPage{
// 		LessonID:       req.GetLessonId(),
// 		CreatedBy:      req.GetCreatedBy(),
// 		LastModifiedBy: req.GetCreatedBy(),
// 		ContentType:    req.GetContentType().String(),
// 		PdfFileUrl:     req.GetPdfFileUrl(),
// 		PdfName:        req.GetPdfName(),
// 	}

// 	pageID, err := s.pageHandlers.CreatePDFPage(ctx, imagePage)
// 	if err != nil {
// 		if errors.Is(err, lessonserv.ErrInvalidCredentials) {
// 			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
// 		}

// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &lpv1.CreatePDFPageResponse{
// 		Id: pageID,
// 	}, nil
// }

// func (s *serverAPI) GetLesson(ctx context.Context, req *lpv1.GetLessonRequest) (*lpv1.GetLessonResponse, error) {
// 	lesson, err := s.lessonHandlers.GetLesson(ctx, req.GetId())
// 	if err != nil {
// 		if errors.Is(err, pageserv.ErrPageNotFound) {
// 			return nil, status.Error(codes.NotFound, "lesson not found")
// 		}

// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &lpv1.GetLessonResponse{
// 		Lesson: &lpv1.Lesson{
// 			Id:             lesson.ID,
// 			Name:           lesson.Name,
// 			CreatedBy:      lesson.CreatedBy,
// 			LastModifiedBy: lesson.LastModifiedBy,
// 			CreatedAt:      timestamppb.New(lesson.CreatedAt),
// 			Modified:       timestamppb.New(lesson.Modified),
// 		},
// 	}, nil
// }
