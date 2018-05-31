package publisher


import (
	"github.com/streadway/amqp"

	"encoding/json"
	"fmt"
	"sync"
	"time"
	"net/url"

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

type PublishTimeoutErr struct{}

func NewTimeoutErr() PublishTimeoutErr {
	return PublishTimeoutErr{}
}

func (PublishTimeoutErr) Error() string {
	return "publish timeout"
}

func IsPublishTimeout(err error) bool {
	_, ok := err.(PublishTimeoutErr)
	return ok
}


const (
	defaultChannelCount = 2
	defaultPubTimeout   = time.Second * 5
)

type PoolConf struct {
	rabbitmq.CommonConfig
	// Name of RabbitMQ exchange
	ExchangeName string
	// Kind of exchange
	ExchangeKind string
	// Enable confirmation
	ConfirmationOn bool
	// Function to encode object into bytes. Default json.Marshal
	EncodeFunc DataEncodeFunc
	// Publishing timeout
	PubTimeout time.Duration
}

// A confirmable kind of message should define how to confirm this message
type ConfirmableMsg interface {
	ConfirmCallback(confirmed bool)
}

// A routable kind of message should define how to get the route key of this message
type RoutableMsg interface {
	GetRouteKey() string
}

type DataEncodeFunc func(data interface{}) ([]byte, error)

func (conf *PoolConf) check() error {
	if err := checkNilString(conf.Host, "Host"); nil != err {
		return err
	}
	if err := checkNilString(conf.User, "User"); nil != err {
		return err
	}
	if err := checkNilString(conf.Password, "Password"); nil != err {
		return err
	}
	if err := checkNilString(conf.ExchangeName, "ExchangeName"); nil != err {
		return err
	}

	// Set default values
	if nil == conf.EncodeFunc {
		conf.EncodeFunc = json.Marshal
	}
	if conf.AMQPChannelSize < 1 {
		conf.AMQPChannelSize = defaultChannelCount
	}
	if conf.PubTimeout < 1 {
		conf.PubTimeout = defaultPubTimeout
	}
	if conf.ExchangeKind == "" {
		conf.ExchangeKind = amqp.ExchangeFanout
	}

	return nil
}

func checkNilString(field, name string) error {
	if "" == field {
		return fmt.Errorf("%s should be specified", name)
	}
	return nil
}

// Pool of publisher
type PublishersPool struct {
	// Pool manager
	poolManager *poolManager
	// Connection manager
	connManager *connManager
	// Channel to get pool
	getCh        <-chan *pool
	exchangeName string
	pubTimeout   time.Duration
}

func NewPublishersPool(conf PoolConf) (*PublishersPool, error) {
	if err := conf.check(); nil != err {
		return nil, err
	}

	getPoolCh := make(chan *pool)
	setPoolCh := make(chan *pool)

	// New managers
	pm := newPoolManager(conf.ExchangeName, getPoolCh, setPoolCh)
	cm := newConnManager(setPoolCh, &conf)

	// Connect to Rabbitmq and init publishers
	pool, err := cm.connect()
	if nil != err {
		return nil, err
	}
	logger.Debugf("init %s pool with revision %d", conf.ExchangeName, pool.rev)
	pm.pool = pool

	// Run managers
	go pm.loop()
	go cm.loop()

	// Init pool
	return &PublishersPool{
		exchangeName: conf.ExchangeName,
		pubTimeout:   conf.PubTimeout,
		poolManager:  pm,
		connManager:  cm,
		getCh:        getPoolCh,
	}, nil
}

func (p *PublishersPool) Size() int {
	return len(p.getPublishers())
}

func (p *PublishersPool) Run() {
	// start a goroutine for each publisher
	for _, pub := range p.getPublishers() {
		go func(pub *publisher) {
			pub.loop()
		}(pub)
	}
}

// Run and wait for all publishers startup.
// WaitGroup should be allocated and WaitGroup.Add should be invoked with PublishersPool.Size() before.
func (p *PublishersPool) RunAndWait(wg *sync.WaitGroup) {
	// start a goroutine for each publisher
	for _, pub := range p.getPublishers() {
		go func(pub *publisher) {
			wg.Done()
			pub.loop()
		}(pub)
	}
}

// Publish a message request to RabbitMQ. Callback will be invoked after publishing confirmation is received.
// mq.PublishTimeoutErr{} will be returned if timeout
func (p *PublishersPool) Publish(data interface{}) error {
	pool := <-p.getCh
	if nil == pool || 0 == len(pool.publishers) {
		return fmt.Errorf("connection is closed")
	}
	logger.Debugf("Get %s pool with revision %d", p.exchangeName, pool.rev)

	timeout := time.After(p.pubTimeout)
	// Wait for an available publisher or timeout
	for {
		select {
		case <-timeout:
			// Return timeout error
			return NewTimeoutErr()
		default:
			// Range publishers to select a writable request channel
			for _, pub := range pool.publishers {
				select {
				case pub.dataCh <- data:
					// Write data to channel and read error channel.
					// Read source codes of publisher.loop() for more details
					err := <-pub.errCh
					return err
				default:
					//Noop so that next publisher can be selected
				}
			}
			// Sleep for a while if there is no available publisher
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (p *PublishersPool) getPublishers() []*publisher {
	pool := <-p.getCh
	if nil == pool {
		return nil
	}
	return pool.publishers
}

type connManager struct {
	conf *PoolConf
	// Established connection
	conn *amqp.Connection
	// Channel to exchange the 'reconnect' signal
	reconnCh chan reconnectSignal
	// Channel to set pool
	setCh chan<- *pool
	// Revision of last pool
	lastRev int64
}

func newConnManager(setPoolCh chan<- *pool, conf *PoolConf) *connManager {
	return &connManager{
		setCh:    setPoolCh,
		conf:     conf,
		reconnCh: make(chan reconnectSignal, 1),
	}
}

// connect and init pool
func (m *connManager) connect() (p *pool, gErr error) {
	uri := (&url.URL{
		Scheme: "amqp",
		User: url.UserPassword(m.conf.User, m.conf.Password),
		Host: m.conf.Host,
	}).String()
	// Connect to RabbitMQ server
	conn, err := amqp.Dial(uri)
	if nil != err {
		return nil, fmt.Errorf("connnect RabbitMQ error: %s", err)
	}
	defer func() {
		// Close connection only when error was returned
		if nil != gErr {
			conn.Close()
		}
	}()

	m.conn = conn
	publishers := make([]*publisher, 0, m.conf.AMQPChannelSize)
	rev := newRevision()

	for i := 0; i < m.conf.AMQPChannelSize; i++ {
		ch, err := conn.Channel()
		if nil != err {
			return nil, fmt.Errorf("open amqp channel error: %s", err)
		}
		defer func(c *amqp.Channel) {
			// Close channel only when error was returned
			// Note: channel should be passed as parameter so that all can be closed
			if nil != gErr {
				c.Close()
			}
		}(ch)

		// Declare exchange
		if err := ch.ExchangeDeclare(m.conf.ExchangeName, m.conf.ExchangeKind, true, false, false, false, nil); nil != err {
			return nil, fmt.Errorf("declare exchange error: %s", err)
		}

		var confirmCh chan amqp.Confirmation
		if m.conf.ConfirmationOn {
			// Set listener for publishing confirmation.
			// A confirmation can be read from confirmCh after a message was published
			confirmCh = ch.NotifyPublish(make(chan amqp.Confirmation, 1))
			if err := ch.Confirm(false); nil != err {
				return nil, fmt.Errorf("confirm.select method send error: %s", err)
			}
		}
		// Add a publisher
		publishers = append(publishers,
			m.newPublisher(ch, confirmCh, rev,
				// Set listener for closed connection/channel notification
				ch.NotifyClose(make(chan *amqp.Error, 1)),
			),
		)
	}

	// update revision
	m.lastRev = rev
	return &pool{
		publishers: publishers,
		rev:        rev,
	}, nil
}

func (m *connManager) newPublisher(ch *amqp.Channel, confirm chan amqp.Confirmation, rev int64,
	closeCh <-chan *amqp.Error) *publisher {

	return &publisher{
		dataCh:       make(chan interface{}),
		errCh:        make(chan error),
		amqpCh:       ch,
		confirmCh:    confirm,
		exchangeName: m.conf.ExchangeName,
		encodeFunc:   m.conf.EncodeFunc,
		stopCh:       make(chan stopSignal, 1),
		reconnCh:     m.reconnCh,
		rev:          rev,
		closeCh:      closeCh,
	}
}

// reconnect pool if nessesary
func (m *connManager) loop() {
	for {
		select {
		case signal := <-m.reconnCh:
			// Wait for reconnection signal
			logger.Debugf("Received reconnection signal of %s whose last revision is %d", m.conf.ExchangeName, m.lastRev)
			// Do not reconnect if signal is old
			if m.lastRev > signal.rev {
				logger.Debugf("%s recieved old signal with revision %d", m.conf.ExchangeName, signal.rev)
				continue
			}

			// Reclaim connection
			if nil != m.conn {
				m.conn.Close()
				m.conn = nil
			}
			// Reclaim publishers
			m.setCh <- nil

			for {
				pool, err := m.connect()
				if nil != err {
					logger.Errorf("Failed to reconnect %s: %#v\n", m.conf.ExchangeName, err)
					// Try to reconnect after 2 seconds
					time.Sleep(2 * time.Second)
					continue
				}

				logger.Infof("Reconnected %s", m.conf.ExchangeName)
				wg := sync.WaitGroup{}
				wg.Add(len(pool.publishers))
				for _, p := range pool.publishers {
					go func(pub *publisher) {
						wg.Done()
						pub.loop()
					}(p)
				}
				wg.Wait()
				logger.Infof("%s pool started", m.conf.ExchangeName)
				m.setCh <- pool
				break
			}
		}
	}
}

func newRevision() int64 {
	// Create revision according to utc
	return time.Now().UTC().UnixNano()
}

type reconnectSignal struct {
	// Revision of current disconnected pool
	rev int64
}

type pool struct {
	// Publishers list
	publishers []*publisher
	// Revision of pool
	rev int64
}

func (p *pool) close() {
	// Close each publisher
	for _, pub := range p.publishers {
		pub.close()
	}
}

type poolManager struct {
	// Channel to get the pool
	getCh chan<- *pool
	// Channel to set pool
	setCh <-chan *pool
	// Exchage name
	exchange string
	*pool
}

func newPoolManager(name string, getPoolCh chan<- *pool, setPoolCh <-chan *pool) *poolManager {
	return &poolManager{
		getCh:    getPoolCh,
		setCh:    setPoolCh,
		exchange: name,
	}
}

func (m *poolManager) loop() {
	for {
		select {
		case m.getCh <- m.pool:
			// Send publishers to getter
		case pool := <-m.setCh:
			// Get new pool and update
			if nil == pool && nil != m.pool {
				logger.Debugf("Close disconnected %s pool with revision %d", m.exchange, m.pool.rev)
				m.pool.close()
				m.pool = nil
			}
			logger.Debugf("set new %s pool %#v", m.exchange, pool)
			m.pool = pool
		}
	}
}

type publisher struct {
	// Channel for sending message
	dataCh chan interface{}
	// Channel for return error
	errCh chan error
	// AMQP Channel to communicate with RabbitMQ
	amqpCh *amqp.Channel
	// Channel for publishing confirmation
	confirmCh chan amqp.Confirmation
	// ExchangeName
	exchangeName string
	// Function to encode object into bytes.
	encodeFunc DataEncodeFunc
	// Channel to exchange stop signal
	stopCh chan stopSignal
	// Channel to send reconnect signal
	reconnCh chan<- reconnectSignal
	// Revision of pool
	rev int64
	// Channel to get notifications of closed connection
	closeCh <-chan *amqp.Error
}

type stopSignal struct{}

func (p *publisher) close() {
	p.stopCh <- stopSignal{}
}

func (p *publisher) loop() {
	// Read message request
	for {
		select {
		case err := <-p.closeCh:
			if nil != err {
				logger.Errorf("Try to reconnect %s pool for close notification: %#v", p.exchangeName, err)
				p.reconnect()
			}
			return
		case <-p.stopCh:
			// stop publisher if close signal is received
			p.amqpCh.Close()
			return
		case data := <-p.dataCh:
			// Encode data to json
			bytes, err := p.encodeFunc(data)
			if nil != err {
				p.errCh <- fmt.Errorf("encode data error: %s", err)
				continue
			}

			routeKey := ""
			if rdata, ok := data.(RoutableMsg); ok {
				routeKey = rdata.GetRouteKey()
			}
			// Publish the request
			err = p.amqpCh.Publish(p.exchangeName, routeKey, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Body:         bytes,
			})
			p.errCh <- err

			if amqp.ErrClosed == err {
				p.reconnect()
			}

			if nil != p.confirmCh {
				logger.Debug("Waiting for confimation")
				// Read confirmation
				confirmation, open := <-p.confirmCh
				if !open {
					logger.Warn("Confirm channel was closed")
				}
				if cdata, ok := data.(ConfirmableMsg); ok {
					logger.Debug("Invoke confirmation callback")
					// Run callback asynchronously
					go cdata.ConfirmCallback(confirmation.Ack)
				} else {
					logger.Debug("Discard confirmation")
				}
			}

		}

	}
	// Never happens
	logger.Warn("AMQP channel was closed")
}

func (p *publisher) reconnect() {
	// Try to send signal to reconnect if connection is closed
	select {
	case p.reconnCh <- reconnectSignal{rev: p.rev}:
	default:
		// Do nothing if reconnection is in progress
	}
}
