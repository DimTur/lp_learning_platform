package utils

type PaginationQueryParams struct {
	Limit  int64 `json:"limit" validate:"required,min=1"`
	Offset int64 `json:"offset" validate:"min=0"`
}

func (p *PaginationQueryParams) SetDefaults() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}
