package config

type Notification struct {
	NotificationQueue      QueueConfig `yaml:"notification_queue"`
	NotificationRoutingKey string      `yaml:"notification_routing_key"`
}
