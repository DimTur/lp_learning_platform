package config

type Channel struct {
	ChannelQueue      QueueConfig     `yaml:"channel_queue"`
	ChannelConsumer   ChannelConsumer `yaml:"channel_consumer"`
	ChannelRoutingKey string          `yaml:"channel_routing_key"`
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
