package config

type Plan struct {
	PlanQueue      PlanQueue    `yaml:"plan_queue"`
	PlanConsumer   PlanConsumer `yaml:"plan_consumer"`
	PlanRoutingKey string       `yaml:"plan_routing_key"`
}
type PlanQueue struct {
	Name        string    `yaml:"name"`
	Durable     bool      `yaml:"durable"`
	AutoDeleted bool      `yaml:"auto_deleted"`
	Exclusive   bool      `yaml:"exclusive"`
	NoWait      bool      `yaml:"no_wait"`
	Args        QueueArgs `yaml:"args"`
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
