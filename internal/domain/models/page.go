package models

import "time"

type ImagePage struct {
	ID             int64
	LessonID       int64
	CreatedBy      int64
	LastModifiedBy int64
	CreatedAt      time.Time
	Modified       time.Time
	ContentType    string

	ImageFileUrl string
	ImageName    string
}

type CreateImagePage struct {
	ID             int64     `json:"id"`
	LessonID       int64     `json:"lesson_id"`
	CreatedBy      int64     `json:"created_by"`
	LastModifiedBy int64     `json:"last_modified_by"`
	CreatedAt      time.Time `json:"created_at"`
	Modified       time.Time `json:"modified"`
	ContentType    string    `json:"content_type"`

	ImageFileUrl string `json:"image_file_url"`
	ImageName    string `json:"image_name"`
}

type VideoPage struct {
	ID             int64
	LessonID       int64
	CreatedBy      int64
	LastModifiedBy int64
	CreatedAt      time.Time
	Modified       time.Time
	ContentType    string

	VideoFileUrl string
	VideoName    string
}

type CreateVideoPage struct {
	ID             int64     `json:"id"`
	LessonID       int64     `json:"lesson_id"`
	CreatedBy      int64     `json:"created_by"`
	LastModifiedBy int64     `json:"last_modified_by"`
	CreatedAt      time.Time `json:"created_at"`
	Modified       time.Time `json:"modified"`
	ContentType    string    `json:"content_type"`

	VideoFileUrl string `json:"video_file_url"`
	VideoName    string `json:"video_name"`
}

type PDFPage struct {
	ID             int64
	LessonID       int64
	CreatedBy      int64
	LastModifiedBy int64
	CreatedAt      time.Time
	Modified       time.Time
	ContentType    string

	PdfFileUrl string
	PdfName    string
}

type CreatePDFPage struct {
	ID             int64     `json:"id"`
	LessonID       int64     `json:"lesson_id"`
	CreatedBy      int64     `json:"created_by"`
	LastModifiedBy int64     `json:"last_modified_by"`
	CreatedAt      time.Time `json:"created_at"`
	Modified       time.Time `json:"modified"`
	ContentType    string    `json:"content_type"`

	PdfFileUrl string `json:"pdf_file_url"`
	PdfName    string `json:"pdf_name"`
}

type Page struct {
	ID             int64
	LessonID       int64
	CreatedBy      int64
	LastModifiedBy int64
	CreatedAt      time.Time
	Modified       time.Time
	ContentType    string
}

type UpdateImagePage struct {
	ID             int64 `json:"id" validate:"required"`
	LastModifiedBy int64 `json:"last_modified_by" validate:"required"`

	ImageFileUrl string `json:"image_file_url,omitempty"`
	ImageName    string `json:"image_name,omitempty"`
}

type UpdateVideoPage struct {
	ID             int64 `json:"id" validate:"required"`
	LastModifiedBy int64 `json:"last_modified_by" validate:"required"`

	VideoFileUrl string `json:"video_file_url,omitempty"`
	VideoName    string `json:"video_name,omitempty"`
}

type UpdatePDFPage struct {
	ID             int64 `json:"id" validate:"required"`
	LastModifiedBy int64 `json:"last_modified_by" validate:"required"`

	PDFFileUrl string `json:"pdf_file_url,omitempty"`
	PDFName    string `json:"pdf_name,omitempty"`
}
