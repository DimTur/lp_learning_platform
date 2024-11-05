package pages

import "time"

type Page interface {
	GetCommonFields() *BasePage
	GetContentTypeSpecificFields() []interface{}
}

type CreatePage interface {
	GetCommonFields() *CreateBasePage
	GetContentTypeSpecificFields() []interface{}
	GetInsertQuery() string
}

type UpdatePage interface {
	GetCommonFields() *UpdateBasePage
	GetContentTypeSpecificFields() []interface{}
	GetUpdateQuery() string
}

type BasePage struct {
	ID             int64     `json:"id"`
	LessonID       int64     `json:"lesson_id"`
	CreatedBy      string    `json:"created_by"`
	LastModifiedBy string    `json:"last_modified_by"`
	CreatedAt      time.Time `json:"created_at"`
	Modified       time.Time `json:"modified"`
	ContentType    string    `json:"content_type"`
}

type ImagePage struct {
	BasePage
	ImageFileUrl string
	ImageName    string
}

type VideoPage struct {
	BasePage
	VideoFileUrl string
	VideoName    string
}

type PDFPage struct {
	BasePage
	PdfFileUrl string
	PdfName    string
}

type CreateBasePage struct {
	LessonID       int64  `json:"lesson_id"`
	CreatedBy      string `json:"created_by"`
	LastModifiedBy string `json:"last_modified_by"`
	ContentType    string `json:"content_type"`
}

type CreateImagePage struct {
	CreateBasePage
	ImageFileUrl string `json:"image_file_url"`
	ImageName    string `json:"image_name"`
}

type CreateVideoPage struct {
	CreateBasePage
	VideoFileUrl string `json:"video_file_url"`
	VideoName    string `json:"video_name"`
}
type CreatePDFPage struct {
	CreateBasePage
	PdfFileUrl string `json:"pdf_file_url"`
	PdfName    string `json:"pdf_name"`
}

type UpdateBasePage struct {
	ID             int64  `json:"id" validate:"required"`
	LastModifiedBy string `json:"last_modified_by" validate:"required"`
	ContentType    string `json:"content_type" validate:"required"`
}

type UpdateImagePage struct {
	UpdateBasePage
	ImageFileUrl string `json:"image_file_url,omitempty"`
	ImageName    string `json:"image_name,omitempty"`
}

type UpdateVideoPage struct {
	UpdateBasePage
	VideoFileUrl string `json:"video_file_url,omitempty"`
	VideoName    string `json:"video_name,omitempty"`
}

type UpdatePDFPage struct {
	UpdateBasePage
	PdfFileUrl string `json:"pdf_file_url,omitempty"`
	PdfName    string `json:"pdf_name,omitempty"`
}

type DBBasePage struct {
	ID             int64     `db:"id"`
	LessonID       int64     `db:"lesson_id"`
	CreatedBy      string    `db:"created_by"`
	LastModifiedBy string    `db:"last_modified_by"`
	CreatedAt      time.Time `db:"created_at"`
	Modified       time.Time `db:"modified"`
	ContentType    string    `db:"content_type"`
}

type DBImagePage struct {
	DBBasePage
	ImageFileUrl string `db:"image_file_url"`
	ImageName    string `db:"image_name"`
}

type DBVideoPage struct {
	DBBasePage
	VideoFileUrl string `db:"video_file_url"`
	VideoName    string `db:"video_name"`
}

type DBPDFPage struct {
	DBBasePage
	PdfFileUrl string `db:"pdf_file_url"`
	PdfName    string `db:"pdf_name"`
}

func (p *ImagePage) GetCommonFields() *BasePage {
	return &p.BasePage
}

func (p *CreateImagePage) GetCommonFields() *CreateBasePage {
	return &p.CreateBasePage
}

func (p *UpdateImagePage) GetCommonFields() *UpdateBasePage {
	return &p.UpdateBasePage
}

func (p ImagePage) GetContentTypeSpecificFields() []interface{} {
	return []interface{}{p.ImageFileUrl, p.ImageName}
}

func (p CreateImagePage) GetContentTypeSpecificFields() []interface{} {
	return []interface{}{p.ImageFileUrl, p.ImageName}
}

func (p UpdateImagePage) GetContentTypeSpecificFields() []interface{} {
	return []interface{}{p.ImageFileUrl, p.ImageName}
}

const createImagePageQuery = `
	INSERT INTO image_imagepage(abstractpage_id, image_file_url, image_name)
	VALUES ($1, $2, $3)`

func (p CreateImagePage) GetInsertQuery() string {
	return createImagePageQuery
}

const updateImagePageQuery = `
	UPDATE
		image_imagepage
	SET
		image_file_url = COALESCE($2, image_file_url),
		image_name = COALESCE($3, image_name)
	WHERE abstractpage_id = $1`

func (p UpdateImagePage) GetUpdateQuery() string {
	return updateImagePageQuery
}

func (p *VideoPage) GetCommonFields() *BasePage {
	return &p.BasePage
}

func (p *CreateVideoPage) GetCommonFields() *CreateBasePage {
	return &p.CreateBasePage
}

func (p *UpdateVideoPage) GetCommonFields() *UpdateBasePage {
	return &p.UpdateBasePage
}

func (p VideoPage) GetContentTypeSpecificFields() []interface{} {
	return []interface{}{p.VideoFileUrl, p.VideoName}
}

func (p CreateVideoPage) GetContentTypeSpecificFields() []interface{} {
	return []interface{}{p.VideoFileUrl, p.VideoName}
}

func (p UpdateVideoPage) GetContentTypeSpecificFields() []interface{} {
	return []interface{}{p.VideoFileUrl, p.VideoName}
}

const createVideoPageQuery = `
	INSERT INTO video_videopage(abstractpage_id, video_file_url, video_name)
	VALUES ($1, $2, $3)`

func (p CreateVideoPage) GetInsertQuery() string {
	return createVideoPageQuery
}

const updateVideoPageQuery = `
	UPDATE video_videopage
	SET
		video_file_url = COALESCE($2, video_file_url),
		video_name = COALESCE($3, video_name)
	WHERE abstractpage_id = $1`

func (p UpdateVideoPage) GetUpdateQuery() string {
	return updateVideoPageQuery
}

func (p *PDFPage) GetCommonFields() *BasePage {
	return &p.BasePage
}

func (p *CreatePDFPage) GetCommonFields() *CreateBasePage {
	return &p.CreateBasePage
}

func (p *UpdatePDFPage) GetCommonFields() *UpdateBasePage {
	return &p.UpdateBasePage
}

func (p PDFPage) GetContentTypeSpecificFields() []interface{} {
	return []interface{}{p.PdfFileUrl, p.PdfName}
}

func (p CreatePDFPage) GetContentTypeSpecificFields() []interface{} {
	return []interface{}{p.PdfFileUrl, p.PdfName}
}

func (p UpdatePDFPage) GetContentTypeSpecificFields() []interface{} {
	return []interface{}{p.PdfFileUrl, p.PdfName}
}

const createPDFPageQuery = `
	INSERT INTO pdf_pdfpage(abstractpage_id, pdf_file_url, pdf_name)
	VALUES ($1, $2, $3)`

func (p CreatePDFPage) GetInsertQuery() string {
	return createPDFPageQuery
}

const updatePDFPageQuery = `
	UPDATE pdf_pdfpage
	SET
		pdf_file_url = COALESCE($2, pdf_file_url),
		pdf_name = COALESCE($3, pdf_name)
	WHERE abstractpage_id = $1`

func (p UpdatePDFPage) GetUpdateQuery() string {
	return updatePDFPageQuery
}
