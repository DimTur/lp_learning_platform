package plans

import "time"

type DBPlan struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	CreatedBy      int64     `db:"created_by"`
	LastModifiedBy int64     `db:"last_modified_by"`
	IsPublished    bool      `db:"is_published"`
	Public         bool      `db:"public"`
	CreatedAt      time.Time `db:"created_at"`
	Modified       time.Time `db:"modified"`
}
