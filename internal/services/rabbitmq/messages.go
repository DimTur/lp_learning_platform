package rabbitmq

type Spfu struct {
	LearningGroupID string   `json:"learning_group_id"`
	UserIDs         []string `json:"user_ids"`
	CreatedBy       string   `json:"created_by"`
}
