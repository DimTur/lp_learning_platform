package channels

import (
	"database/sql"
	"time"
)

type DBChannel struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	CreatedBy      int64     `db:"created_by"`
	LastModifiedBy int64     `db:"last_modified_by"`
	CreatedAt      time.Time `db:"created_at"`
	Modified       time.Time `db:"modified"`
}

type DBChannelWithPlans struct {
	ID             int64              `db:"id"`
	Name           string             `db:"name"`
	Description    string             `db:"description"`
	CreatedBy      int64              `db:"created_by"`
	LastModifiedBy int64              `db:"last_modified_by"`
	CreatedAt      time.Time          `db:"created_at"`
	Modified       time.Time          `db:"modified"`
	Plans          []DBPlanInChannels `db:"plans"`
}

// We can use DBPlan, but for flexibility and control returning fields init new struct
type DBPlanInChannels struct {
	ID             sql.NullInt64  `db:"id"`
	Name           sql.NullString `db:"name"`
	Description    sql.NullString `db:"description"`
	CreatedBy      sql.NullInt64  `db:"created_by"`
	LastModifiedBy sql.NullInt64  `db:"last_modified_by"`
	IsPublished    sql.NullBool   `db:"is_published"`
	Public         sql.NullBool   `db:"public"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	Modified       sql.NullTime   `db:"modified"`
}
