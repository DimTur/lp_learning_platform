package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/attempts"
)

func (r *RedisClient) SavePageAttempt(ctx context.Context, pageAttempt *SavePageAttempt) error {
	const op = "storage.redis.SavePageAttempt"

	key := fmt.Sprintf("lesson_attempt:%d", pageAttempt.LessonAttemptID)

	field := fmt.Sprintf("page_attempt:%d", pageAttempt.PageAttemptID)

	value, err := json.Marshal(map[string]interface{}{
		"user_answer": pageAttempt.UserAnswer,
		"is_correct":  pageAttempt.IsCorrect,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := r.client.HSet(ctx, key, field, value).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RedisClient) GetPageAttempts(ctx context.Context, lessonAttemptID int64) ([]attempts.QuestionPageAttempt, error) {
	const op = "storage.redis.GetPageAttempts"

	key := fmt.Sprintf("lesson_attempt:%d", lessonAttemptID)

	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get all page attempts: %w", op, err)
	}

	var pageAttempts []attempts.QuestionPageAttempt

	for field, value := range result {
		// Parse PageAttempID from the name fields
		var pageAttemptID int64
		_, err := fmt.Sscanf(field, "page_attempt:%d", &pageAttemptID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to parse field name %q: %w", op, field, err)
		}

		// Deserialize the JSON value into a structure
		var data struct {
			PageID     int64  `json:"page_id"`
			IsCorrect  bool   `json:"is_correct"`
			UserAnswer string `json:"user_answer"`
		}
		if err := json.Unmarshal([]byte(value), &data); err != nil {
			return nil, fmt.Errorf("%s: failed to unmarshal value: %w", op, err)
		}

		pageAttempts = append(pageAttempts, attempts.QuestionPageAttempt{
			ID:              pageAttemptID,
			PageID:          data.PageID,
			LessonAttemptID: lessonAttemptID,
			IsCorrect:       data.IsCorrect,
		})
	}

	return pageAttempts, nil
}
