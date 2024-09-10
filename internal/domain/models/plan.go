package models

import "time"

type Plan struct {
	ID             int64
	Name           string
	Description    string
	CreatedBy      int64
	LastModifiedBy int64
	IsPublished    bool
	PublishedAt    time.Time
	Public         bool
	Lessons        []Lesson
	Channels       []int64
}
