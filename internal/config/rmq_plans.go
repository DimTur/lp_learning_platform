package config

type Plan struct {
	PlanQueue      QueueConfig  `yaml:"plan_queue"`
	PlanConsumer   PlanConsumer `yaml:"plan_consumer"`
	PlanRoutingKey string       `yaml:"plan_routing_key"`
}

type PlanConsumer struct {
	Queue        string       `yaml:"queue"`
	Consumer     string       `yaml:"consumer"`
	AutoAck      bool         `yaml:"autoAck"`
	Exclusive    bool         `yaml:"exclusive"`
	NoLocal      bool         `yaml:"noLocal"`
	NoWait       bool         `yaml:"noWait"`
	ConsumerArgs ConsumerArgs `yaml:"args"`
}
