package utils

type PaginationQueryParams struct {
	Limit  int64 `json:"limit" validate:"required,min=1"`
	Offset int64 `json:"offset" validate:"min=0"`
}
