package models

type Channel struct {
	ID             int64
	Name           string
	Description    string
	CreatedBy      int64
	LastModifiedBy int64
	Public         bool
	Plans          []Plan
}
