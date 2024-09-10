package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/DimTur/lp_learning_platform/internal/domain/models"
	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/mattn/go-sqlite3"
)

type SQLLiteStorage struct {
	db *sql.DB
}

// New creates a new instance of the SQLite storage
func New(storagePath string) (SQLLiteStorage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return SQLLiteStorage{}, fmt.Errorf("%s: %w", op, err)
	}

	return SQLLiteStorage{db: db}, nil
}

func (s *SQLLiteStorage) Close() error {
	return s.db.Close()
}

// SaveChannel saves channel to db.
func (s *SQLLiteStorage) SaveChannel(ctx context.Context, name string, description string, userID int64, public bool) (int64, error) {
	const op = "storage.sqlite.channel.SaveChannel"

	stmt, err := s.db.Prepare("INSERT INTO lp_channels(name, description, created_by, last_modified_by, public) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	publicInt := 0
	if public {
		publicInt = 1
	}

	res, err := stmt.ExecContext(
		ctx,
		name,
		description,
		userID,
		userID,
		publicInt,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrChannelExitsts)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	chanID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return chanID, nil
}

// GetChannelByID returns channel by ID.
func (s *SQLLiteStorage) GetChannelByID(ctx context.Context, chanID int64) (models.Channel, error) {
	const op = "storage.sqlite.channel.GetChannelByID"

	stmt, err := s.db.Prepare("SELECT id, name, description, created_by, last_modified_by, public FROM lp_channels WHERE id = ?")
	if err != nil {
		return models.Channel{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, chanID)

	var channel models.Channel
	err = row.Scan(
		&channel.ID,
		&channel.Name,
		&channel.Description,
		&channel.CreatedBy,
		&channel.LastModifiedBy,
		&channel.Public,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Channel{}, fmt.Errorf("%s: %w", op, storage.ErrChannelNotFound)
		}

		return models.Channel{}, fmt.Errorf("%s: %w", op, err)
	}

	return channel, nil
}
