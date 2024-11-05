package config

type Channel struct {
	ChannelQueue      ChannelQueue    `yaml:"channel_queue"`
	ChannelConsumer   ChannelConsumer `yaml:"channel_consumer"`
	ChannelRoutingKey string          `yaml:"channel_routing_key"`
}
type ChannelQueue struct {
	Name        string    `yaml:"name"`
	Durable     bool      `yaml:"durable"`
	AutoDeleted bool      `yaml:"auto_deleted"`
	Exclusive   bool      `yaml:"exclusive"`
	NoWait      bool      `yaml:"no_wait"`
	Args        QueueArgs `yaml:"args"`
}

type ChannelConsumer struct {
	Queue        string       `yaml:"queue"`
	Consumer     string       `yaml:"consumer"`
	AutoAck      bool         `yaml:"autoAck"`
	Exclusive    bool         `yaml:"exclusive"`
	NoLocal      bool         `yaml:"noLocal"`
	NoWait       bool         `yaml:"noWait"`
	ConsumerArgs ConsumerArgs `yaml:"args"`
}
