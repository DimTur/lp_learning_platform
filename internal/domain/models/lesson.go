package models

type Lesson struct {
	ID             int64
	Name           string
	Description    string
	PassPercentage int64
	CreatedBy      int64
	LastModifiedBy int64
	Plans          []int64
}
