package rabbitmq

type CommonConfig struct {
	// Host and port
	Host string
	// RabbitMQ user
	User string
	// RabbitMQ password
	Password string
	// AMQP Channel count on a connection
	AMQPChannelSize int
}
