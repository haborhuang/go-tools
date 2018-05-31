package subscriber

import (
	"fmt"
	"time"
	"net/url"

	"github.com/streadway/amqp"

	"github.com/haborhuang/go-tools/rabbitmq"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	l, err := zap.NewDevelopment()
	if nil != err {
		panic(err)
	}

	logger = l.Sugar()
}

const (
	defaultChannelCount = 2
	defaultMaxMultiMsgs = 5
	defaultMultiMsgsTimeout = 1000 // 1 second
)

func defaultBindingKeys() []string {
	return []string{"#"}
}

type PoolConfig struct {
	rabbitmq.CommonConfig
	// Name of exchange
	ExchangeName string
	// Kind of exchange
	ExchangeKind string
	// Name of queue
	QueueName string
	// Route key of binding
	QueueBindingKeys []string
	// Acknowledge message automatically or not
	AutoAck bool
	// Configurations of multiple acknowledgements. Take effect only when AutoAck is false.
	// If AutoAck is false and MultipleAck is nil.
	MultipleAck *MultipleAckConfig
	// Handler to handle a message at a time.
	// Should be specified when AutoAck is true or MultipleAck is nil
	MsgHandler MsgHandler
	// Handler to handle batch messages.
	// Should be specified when AutoAck is false and MultipleAck is not nil
	BatchHandler MultipleHandler
}

type MultipleAckConfig struct {
	// Maximum count of messages to be handled at a time
	MaxCount int
	// Maximum milliseconds to wait to receive a new message
	Timeout int64
}

// Check validation and set default values
func (c *PoolConfig) check() error {
	if "" == c.Host || "" == c.User || "" == c.Password {
		return fmt.Errorf("Host, user and password should be specified")
	}

	if "" == c.ExchangeName || "" == c.QueueName {
		return fmt.Errorf("Name of exchange and queue should be specified")
	}

	if !c.AutoAck && nil != c.MultipleAck {
		if nil == c.BatchHandler {
			return fmt.Errorf("Message handler should be specified")
		}
		if c.MultipleAck.MaxCount < 2 {
			c.MultipleAck.MaxCount = defaultMaxMultiMsgs
		}
		if c.MultipleAck.Timeout < 1 {
			c.MultipleAck.Timeout = defaultMultiMsgsTimeout
		}
		c.MultipleAck.Timeout = c.MultipleAck.Timeout * int64(time.Millisecond)
	} else {
		if nil == c.MsgHandler {
			return fmt.Errorf("Message handler should be specified")
		}
	}

	if c.AMQPChannelSize < 1 {
		c.AMQPChannelSize = defaultChannelCount
	}

	if c.ExchangeKind == "" {
		c.ExchangeKind = amqp.ExchangeFanout
	}

	if len(c.QueueBindingKeys) < 1 {
		c.QueueBindingKeys = defaultBindingKeys()
	}

	return nil
}

type consumer struct {
	MsgHandler
	MultipleHandler
	autoAck bool
	multipleAck *MultipleAckConfig
}

type MsgHandler interface {
	HandleMQMessage(data []byte) (ack bool, err error)
}

type MultipleHandler interface {
	HandleMQMultiple(data [][]byte) (ack bool, err error)
}

type SubscribersPool struct {
	conf *PoolConfig
	// Established connection
	conn *amqp.Connection
	// Subscribers list
	subscribers []*subscriber
	// Channel to get reconnection signal
	reconnCh chan reconnSignal
	// Revision of last connected pool
	lastRev int64
}

// Signal to notify reconnection
type reconnSignal struct {
	// Revision of pool that should be reconnected
	rev int64
}

func NewSubscribersPool(conf PoolConfig) (*SubscribersPool, error) {
	if err := conf.check(); nil != err {
		return nil, err
	}

	reconnectCh := make(chan reconnSignal, 1)

	// Init pool
	pool := &SubscribersPool{
		conf:     &conf,
		reconnCh: reconnectCh,
	}

	// Establish connections
	if err := pool.connect(); nil != err {
		return nil, err
	}

	// Running pool's handle loop of reconnection
	go pool.reconnLoop()
	logger.Infof("%s subscribers pool started", conf.ExchangeName)
	return pool, nil
}

// Running subscribers in pool
func (c *SubscribersPool) Run() {
	for _, s := range c.subscribers {
		go s.run()
	}
}

func (c *SubscribersPool) connect() (gErr error) {
	uri := (&url.URL{
		Scheme: "amqp",
		User: url.UserPassword(c.conf.User, c.conf.Password),
		Host: c.conf.Host,
	}).String()
	logger.Debugf("connecting %s", uri)
	// Connect to RabbitMQ server
	conn, err := amqp.Dial(uri)
	if nil != err {
		return fmt.Errorf("connnect RabbitMQ error: %s", err)
	}
	defer func() {
		// Close connection only when error was returned
		if nil != gErr {
			conn.Close()
		}
	}()

	// Init pool
	c.conn = conn
	c.subscribers = make([]*subscriber, 0, c.conf.AMQPChannelSize)
	c.lastRev = newRevision()

	for i := 0; i < c.conf.AMQPChannelSize; i++ {
		ch, err := conn.Channel()
		if nil != err {
			return fmt.Errorf("open amqp channel error: %s", err)
		}
		defer func(c *amqp.Channel) {
			// Close channel only when error was returned
			// Note: channel should be passed as parameter so that all can be closed
			if nil != gErr {
				c.Close()
			}
		}(ch)

		// Bind queue
		if err := ch.ExchangeDeclare(c.conf.ExchangeName, c.conf.ExchangeKind, true, false, false, false, nil); nil != err {
			return fmt.Errorf("declare exchange error: %s", err)
		}
		if _, err := ch.QueueDeclare(c.conf.QueueName, true, false, false, false, nil); nil != err {
			return fmt.Errorf("declare queue error: %s", err)
		}

		// Bind queue for each route key
		for _, rk := range c.conf.QueueBindingKeys {
			logger.Debugf("binding queue %s to exchange %s with route key %s", c.conf.QueueName, c.conf.ExchangeName, rk)
			if err := ch.QueueBind(c.conf.QueueName, rk, c.conf.ExchangeName, false, nil); nil != err {
				return fmt.Errorf("bind queue error: %s", err)
			}
		}

		// Get channel for consuming
		deliveryCh, err := ch.Consume(c.conf.QueueName, "", c.conf.AutoAck, false, false, false, nil)
		if nil != err {
			return fmt.Errorf("start consuming error: %s", err)
		}

		// Add subscriber
		sub := newSubscriber(deliveryCh, &consumer{
			MsgHandler: c.conf.MsgHandler,
			MultipleHandler: c.conf.BatchHandler,
			autoAck: c.conf.AutoAck,
			multipleAck: c.conf.MultipleAck,
		}, c.lastRev, c.reconnCh)
		c.subscribers = append(c.subscribers, sub)
	}

	return nil
}

func newRevision() int64 {
	// Create revision according to utc
	return time.Now().UTC().UnixNano()
}

func (c *SubscribersPool) reconnLoop() {
	for {
		// Get reconnection signal
		signal := <-c.reconnCh

		logger.Debugf("Received a reconnection signal of %s pool whose last revision is %d", c.conf.ExchangeName, c.lastRev)
		if c.lastRev > signal.rev {
			logger.Debugf("Received old reconnection signal with revision %d", signal.rev)
			continue
		}
		logger.Infof("Reconnecting %s pool", c.conf.ExchangeName)

		// Close connection
		if nil != c.conn {
			c.conn.Close()
			c.conn = nil
			c.subscribers = nil
		}

		// Establish new connection
		for {
			if err := c.connect(); nil != err {
				logger.Errorf("Failed to reconnect %s pool: %#v", c.conf.ExchangeName, err)
				time.Sleep(2 * time.Second)
				continue
			}
			break
		}

		logger.Infof("Reconnected %s pool with revision %d", c.conf.ExchangeName, c.lastRev)
		// Run subscribers
		c.Run()
	}
}

type subscriber struct {
	// AMQP Channel to communicate with RabbitMQ
	deliveryCh <-chan amqp.Delivery
	// Message handler
	consumer *consumer
	// Pool revision
	rev int64
	// Channel to send reconnection signal
	reconnCh chan<- reconnSignal
}

func newSubscriber(ch <-chan amqp.Delivery, consumer *consumer, rev int64, reconnCh chan<- reconnSignal) *subscriber {
	return &subscriber{
		deliveryCh: ch,
		consumer:   consumer,
		rev:        rev,
		reconnCh:   reconnCh,
	}
}

func (s *subscriber) run() {
	if !s.consumer.autoAck && nil != s.consumer.multipleAck {
		s.batchLoop()
	} else {
		s.loop()
	}
	logger.Info("Subscriber stopped")
	select {
	case s.reconnCh <- reconnSignal{rev: s.rev}:
		// Try to send reconnection signal
		logger.Infof("Trigger reconnection")
	default:
		// Reconnection is in progress
	}
}

func (s *subscriber) batchLoop() {
	msgs := make([][]byte, 0, s.consumer.multipleAck.MaxCount)
	for {
		// Wait for the first message
		d, ok := <- s.deliveryCh
		if !ok {
			// close loop if channel was closed
			return
		}
		msgs = append(msgs, d.Body)

		// Wait for other messages or timeout
		for ; len(msgs) < s.consumer.multipleAck.MaxCount; {
			select {
			case <-time.After(time.Duration(s.consumer.multipleAck.Timeout)):
				// timeout
				break
			case d, ok = <- s.deliveryCh:
				if !ok {
					// close loop if channel was closed
					return
				}
				msgs = append(msgs, d.Body)
			}
		}

		// Call handler to handle messages
		ack, err := s.consumer.HandleMQMultiple(msgs)
		if nil != err {
			logger.Errorf("Handle MQ messages error: %v", err)
		}

		// Acknowledge or reject messages
		if ack {
			d.Ack(true)
		} else {
			d.Nack(true, true)
		}

		// Clear messages slice
		msgs = msgs[:0]
	}
}

func (s *subscriber) loop() {
	for d := range s.deliveryCh {
		// Call handler to handle the message
		ack, err := s.consumer.HandleMQMessage(d.Body)
		if nil != err {
			logger.Errorf("Handle MQ message error: %v", err)
		}

		if !s.consumer.autoAck {
			// Acknowledge or reject message
			if ack {
				d.Ack(false)
			} else {
				d.Reject(true)
			}
		}
	}
}