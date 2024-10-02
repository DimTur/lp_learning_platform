package lessons

import "time"

type DBLesson struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	CreatedBy      int64     `db:"created_by"`
	LastModifiedBy int64     `db:"last_modified_by"`
	CreatedAt      time.Time `db:"created_at"`
	Modified       time.Time `db:"modified"`
}
