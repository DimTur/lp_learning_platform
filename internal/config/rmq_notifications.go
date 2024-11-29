package config

type Notification struct {
	NotificationQueue      QueueConfig `yaml:"notification_queue"`
	NotificationRoutingKey string      `yaml:"plan_routing_key"`
}
