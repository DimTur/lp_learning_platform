package config

type RabbitMQ struct {
	UserName      string        `yaml:"username"`
	Password      string        `yaml:"password"`
	Host          string        `yaml:"host"`
	Port          int           `yaml:"port"`
	ShareExchange ShareExchange `yaml:"share_exchange"`
	Channel       Channel       `yaml:"channel"`
	Plan          Plan          `yaml:"plan"`
}

type ShareExchange struct {
	Name        string       `yaml:"name"`
	Kind        string       `yaml:"kind"`
	Durable     bool         `yaml:"durable"`
	AutoDeleted bool         `yaml:"auto_deleted"`
	Internal    bool         `yaml:"internal"`
	NoWait      bool         `yaml:"no_wait"`
	Args        ExchangeArgs `yaml:"args"`
}

type ExchangeArgs struct {
	AltExchange string `yaml:"alternate_exchange"`
}

type QueueArgs struct {
	XMessageTtl int32 `yaml:"x_message_ttl"`
}

type ConsumerArgs struct {
	XConsumerTtl       int32 `yaml:"x-consumer-timeout"`
	XConsumerPrefCount int32 `yaml:"x-consumer-prefetch-count"`
}

func (e ExchangeArgs) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"alternate-exchange": e.AltExchange,
	}
}

func (q QueueArgs) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"x-message-ttl": q.XMessageTtl,
	}
}

func (c ConsumerArgs) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"x-consumer-timeout":        c.XConsumerTtl,
		"x-consumer-prefetch-count": c.XConsumerPrefCount,
	}
}
